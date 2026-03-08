package dao

import (
	"fmt"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/zlog"
)

func Order(out interface{}, order interface{}, query interface{}, args ...interface{}) error {
	if err := GormDB.Order(order).Where(query).Find(out).Error; err != nil {
		zlog.Error(err.Error())
		return err
	}
	return nil
}

func GetUserContact(identifier interface{}) (*model.UserInfo, error) {
	var user model.UserInfo
	db := GormDB

	switch v := identifier.(type) {
	case UserUUID:
		db = db.Where("uuid = ?", string(v))
	case UserTelephone:
		db = db.Where("telephone = ?", string(v))
	default:
		return nil, fmt.Errorf("unsupported identifier type: %T", identifier)
	}

	if err := db.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
