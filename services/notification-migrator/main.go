package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	kafka "github.com/segmentio/kafka-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
	mysqlModel "github.com/alphafast/asmt-fw/libs/noti/repository/mysql"
	"github.com/alphafast/asmt-fw/libs/utils/env"
)

func main() {
	// initial kafka topics
	kafkaHost := env.RequiredEnv("KAFKA_HOST")
	notiReqTopic := env.RequiredEnv("KAFKA_NOTIFICATION_REQUEST_TOPIC_NAME")
	notiReqPartitionSize := env.ToInt(env.RequiredEnv("KAFKA_NOTIFICATION_REQUEST_TOPIC_PARTITION_SIZE"))
	notiProcessedTopic := env.RequiredEnv("KAFKA_NOTIFICATION_PROCESSED_TOPIC_NAME")
	notiProcessedPartitionSize := env.ToInt(env.RequiredEnv("KAFKA_NOTIFICATION_PROCESSED_TOPIC_PARTITION_SIZE"))

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
		log.Panic().Err(err).Msg("[main]: failed to connect to database")
	}

	// migrate kafka topics
	conn, err := kafka.Dial("tcp", kafkaHost)
	if err != nil {
		log.Panic().Err(err).Msg("[main]: failed to connect to kafka")
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Panic().Err(err).Msg("[main]: failed to get controller")
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Panic().Err(err).Msg("[main]: failed to connect to controller")
	}
	defer controllerConn.Close()

	topicsCreator := []kafka.TopicConfig{
		{
			Topic:             notiReqTopic,
			NumPartitions:     notiReqPartitionSize,
			ReplicationFactor: 1,
		},
		{
			Topic:             notiProcessedTopic,
			NumPartitions:     notiProcessedPartitionSize,
			ReplicationFactor: 1,
		},
	}

	for _, topic := range topicsCreator {
		err = controllerConn.CreateTopics(topic)
		if err != nil {
			log.Panic().Err(err).Msgf("[main]: failed to create topic %s", topic.Topic)
		}
	}

	log.Info().Msg("[main]: kafka topics created successfully")

	// migrate seed data
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn: dbConn,
		}),
		&gorm.Config{},
	)
	if err != nil {
		log.Panic().Err(err).Msg("[main]: failed to bind connection to gorm")
	}

	// drop table
	err = db.Migrator().DropTable(&mysqlModel.MySqlNotiResult{}, &mysqlModel.MySqlNotiUser{})
	if err != nil {
		log.Panic().Err(err).Msg("[main]: failed to drop tables")
	}

	// create table
	err = db.AutoMigrate(&mysqlModel.MySqlNotiResult{}, &mysqlModel.MySqlNotiUser{})
	if err != nil {
		log.Panic().Err(err).Msg("[main]: failed to auto migrate tables")
	}

	// insert seed data
	content, err := os.ReadFile("./seed/noti_user.json")
	if err != nil {
		log.Panic().Err(err).Msg("[main]: failed to read seed file")
	}

	var notiUsers []model.NotiUser
	err = json.Unmarshal(content, &notiUsers)
	if err != nil {
		log.Panic().Err(err).Msg("[main]: failed to unmarshal seed data")
	}

	for _, notiUser := range notiUsers {
		err = db.Create(mysqlModel.ToMySqlNotiUser(&notiUser)).Error
		if err != nil {
			log.Panic().Err(err).Msg("[main]: failed to insert seed data")
		}
	}

	log.Info().Msg("[main]: seed data inserted successfully")
}
