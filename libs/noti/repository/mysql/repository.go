package mysql

import (
	"context"
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
	"github.com/pkg/errors"
)

type NotiMySQLRepository struct {
	db *gorm.DB
}

func New(conn *sql.DB) (*NotiMySQLRepository, error) {
	// Bind connection to gorm
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn: conn,
		}),
		&gorm.Config{},
	)
	if err != nil {
		return nil, errors.Wrap(err, "[NotiMySQLRepository.New] error while binding connection to gorm")
	}

	if err := db.AutoMigrate(&MySqlNotiResult{}, &MySqlNotiUser{}); err != nil {
		return nil, errors.Wrap(err, "[NotiMySQLRepository.New] error while auto migrating tables")
	}

	return &NotiMySQLRepository{
		db: db,
	}, nil
}

func (r *NotiMySQLRepository) GetNotifyResultsByReqID(ctx context.Context, reqID string) ([]model.NotiResult, error) {
	targetResults := []MySqlNotiResult{}
	txn := r.db.Where(&MySqlNotiResult{ReqID: reqID}).Find(&targetResults)
	if txn.Error != nil {
		return nil, errors.Wrap(txn.Error, "[NotiMySQLRepository.GetNotifyResultsByRequestID] error while finding notify results")
	}

	results := []model.NotiResult{}
	for _, result := range targetResults {
		results = append(results, *ToNotiResult(&result))
	}

	return results, nil
}

func (r *NotiMySQLRepository) FindUserNotification(ctx context.Context, userID string) (*model.NotiUser, error) {
	targetUser := MySqlNotiUser{UserID: userID}
	txn := r.db.Take(&targetUser)
	if txn.Error != nil {
		if errors.Is(txn.Error, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(model.NotiUserNotFoundError, "[NotiMySQLRepository.FindUserNotification] error while finding user notification")
		}
	}

	return ToNotiUser(&targetUser), nil
}

func (r *NotiMySQLRepository) UpsertNotifyResult(ctx context.Context, result model.NotiResult) error {
	mySqlNotiResult := ToMySqlNotiResult(&result)
	txn := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(mySqlNotiResult)
	if txn.Error != nil {
		return errors.Wrap(txn.Error, "[NotiMySQLRepository.UpsertNotifyResult] error while upserting notify result")
	}

	return nil
}
