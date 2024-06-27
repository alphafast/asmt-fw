module github.com/alphafast/asmt-fw/services/notification-migrater

go 1.21.0

require (
	github.com/alphafast/asmt-fw/libs v0.0.0
	github.com/rs/zerolog v1.33.0
	github.com/segmentio/kafka-go v0.4.47
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.10
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.19.0 // indirect
)

replace github.com/alphafast/asmt-fw/libs => ../../libs
