package main

import (
	"database/sql"
	"fmt"
	"net/url"

	echo "github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	env "github.com/alphafast/asmt-fw/libs/utils/env"
	kafkaProducer "github.com/alphafast/asmt-fw/libs/utils/kafka/producer"
	logUtil "github.com/alphafast/asmt-fw/libs/utils/log"

	notiModel "github.com/alphafast/asmt-fw/libs/domain/noti/model"
	fcmAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/fcm"
	multiAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/multi"
	senGridAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/sengrid"
	sesAdapter "github.com/alphafast/asmt-fw/libs/noti/adapter/ses"
	notiHandler "github.com/alphafast/asmt-fw/libs/noti/handler/http/echo"
	notiMysqlRepo "github.com/alphafast/asmt-fw/libs/noti/repository/mysql"
	notiUsercase "github.com/alphafast/asmt-fw/libs/noti/usecase"
)

func main() {
	logger := logUtil.SetupZeroLog()

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
	requestNotiTopic := env.RequiredEnv("KAFKA_NOTIFICATION_REQUEST_TOPIC")
	producer := kafkaProducer.NewProducer(host)

	// initialize handler
	port := env.RequiredEnv("APP_PORT")
	e := echo.New()
	deps := notiHandler.NotiHandlerDeps{
		NotiUseCase: notiService,
		Producer:    producer,
	}
	opts := notiHandler.NotiHandlerOps{
		NotiRequestTopicName: requestNotiTopic,
	}
	notiHandler.NewRequestNotificationHandler(e, deps, opts)

	err = e.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Panic().Err(err).Msg("[main] error start echo server")
	}
}
