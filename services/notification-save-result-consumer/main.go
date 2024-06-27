package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	// _ "github.com/go-sql-driver/mysql"

	fcmAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/fcm"
	multiAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/multi"
	senGridAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/sengrid"
	sesAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/ses"
	notiConsumerHandler "github.com/alphafast/asmt-fw/libs/noti/handler/queue/kafka"
	notiMysqlRepo "github.com/alphafast/asmt-fw/libs/noti/repository/mysql"
	notiUsercase "github.com/alphafast/asmt-fw/libs/noti/usecase"
	"github.com/alphafast/asmt-fw/libs/utils/env"
	consumer "github.com/alphafast/asmt-fw/libs/utils/kafka/consumer"
	logUtil "github.com/alphafast/asmt-fw/libs/utils/log"
)

func main() {
	ctx, logger := logUtil.SetupZeroLogWithCtx(context.Background())

	// initialize notification repository
	dbHost := env.RequiredEnv("MYSQL_HOST")
	dbPort := env.RequiredEnv("MYSQL_PORT")
	dbUser := env.RequiredEnv("MYSQL_USER")
	dbPass := env.RequiredEnv("MYSQL_PASSWORD")
	dbName := env.RequiredEnv("MYSQL_DATABASE")
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Bangkok")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil {
		logger.Panic().Err(err).Msg("[main]: failed to connect to database")
	}
	notiRepo, err := notiMysqlRepo.New(dbConn)
	if err != nil {
		logger.Panic().Err(err).Msg("[main]: failed to initialize notification repository")
	}

	// initialize notification provider adapter
	sesMailProviderAdapter := sesAdapter.New()
	sendGridMailProviderAdapter := senGridAdapter.New()
	fcmPushProviderAdapter := fcmAdapter.New()

	// initialize multi adapter  with backup
	multiNotiAdapter := multiAdapter.NewNotiAdapter(sesMailProviderAdapter, fcmPushProviderAdapter, multiAdapter.WithBackupEmailAdapter(sendGridMailProviderAdapter))

	// initialize notification service
	notiService := notiUsercase.New(notiRepo, multiNotiAdapter)

	// initialize kafka consumer
	host := env.RequiredEnv("KAFKA_HOST")
	topic := env.RequiredEnv("KAFKA_NOTIFICATION_PROCESSED_TOPIC")
	groupId := env.RequiredEnv("KAFKA_CONSUMER_GROUP_ID")
	consumer := consumer.NewConsumer(consumer.ConsumerConfig{
		Context: ctx,
		Host:    host,
		GroupId: groupId,
		Topic:   topic,
	})

	// initialize consumer handler
	consumerHandler := notiConsumerHandler.NewSaveNotiResultConsumerHandler(ctx, notiService, notiRepo)

	logger.Info().Msg("[main]: consumer started, being consume messages from kafka")

	// start consuming messages from kafka
	consumer.Consume(consumerHandler.HandleNotiSaveResult)
}
