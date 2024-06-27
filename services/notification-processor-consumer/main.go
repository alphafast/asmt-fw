package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"

	notiModel "github.com/alphafast/asmt-fw/libs/domain/noti/model"
	fcmAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/fcm"
	multiAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/multi"
	senGridAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/sengrid"
	sesAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/ses"
	notiConsumerHandler "github.com/alphafast/asmt-fw/libs/noti/handler/queue/kafka"
	notiMysqlRepo "github.com/alphafast/asmt-fw/libs/noti/repository/mysql"
	notiUsercase "github.com/alphafast/asmt-fw/libs/noti/usecase"
	env "github.com/alphafast/asmt-fw/libs/utils/env"
	kafkaConsumer "github.com/alphafast/asmt-fw/libs/utils/kafka/consumer"
	kafkaProducer "github.com/alphafast/asmt-fw/libs/utils/kafka/producer"
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
	ucDeps := notiUsercase.NotiUsecaseDeps{
		NotiRepo:    notiRepo,
		NotiAdapter: multiNotiAdapter,
	}
	ucConf := notiUsercase.NotiUsecaseConf{
		NotiChannelByEventSource: notiModel.DefaultNotiChannelBySourceEvent,
	}
	notiService := notiUsercase.New(ucDeps, ucConf)

	// initialize kafka consumer and producer
	host := env.RequiredEnv("KAFKA_HOST")
	topic := env.RequiredEnv("KAFKA_NOTIFICATION_REQUEST_TOPIC")
	groupId := env.RequiredEnv("KAFKA_CONSUMER_GROUP_ID")
	consumer := kafkaConsumer.NewConsumer(kafkaConsumer.ConsumerConfig{
		Context: ctx,
		Host:    host,
		GroupId: groupId,
		Topic:   topic,
	})
	processedTopic := env.RequiredEnv("KAFKA_NOTIFICATION_PROCESSED_TOPIC")
	producer := kafkaProducer.NewProducer(host)

	// initialize notification request consumer handler
	deps := notiConsumerHandler.ProcessNotiRequestHandlerDeps{
		NotiUseCase: notiService,
		Producer:    producer,
	}
	confs := notiConsumerHandler.ProcessNotiRequestHandlerOps{
		ProcessedTopic: processedTopic,
	}
	consumerHandler := notiConsumerHandler.NewProcessNotiRequestHandler(ctx, deps, confs)

	logger.Info().Msg("[main]: consumer started, being consume messages from kafka")

	// start consuming messages from kafka
	consumer.Consume(consumerHandler.HandleNotiRequest)
}
