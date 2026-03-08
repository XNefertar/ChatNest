package dao

import (
	"errors"
	"fmt"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/zlog"

	"gorm.io/gorm"
)

type UserUUID string
type UserTelephone string
type UserOwnerID string

// IsRecordNotFound 检查错误是否是记录未找到
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func GetUserInfo(identifier interface{}) (*model.UserInfo, error) {
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

func Create(value interface{}) error {
	if err := GormDB.Create(value).Error; err != nil {
		zlog.Error(err.Error())
		return err
	}
	return nil
}

func Save(identifier interface{}) error {
	if err := GormDB.Save(identifier).Error; err != nil {
		zlog.Error(err.Error())
		return err
	}
	return nil
}

func GetAllUsersExcept(identifier interface{}) ([]model.UserInfo, error) {
	var users []model.UserInfo
	db := GormDB
	switch v := identifier.(type) {
	case UserOwnerID:
		db = db.Where("uuid != ?", string(v))
	}

	if err := db.Unscoped().Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsersByUUIDs(uuidList []string) ([]model.UserInfo, error) {
	var users []model.UserInfo
	if err := GormDB.Model(model.UserInfo{}).Where("uuid in (?)", uuidList).Find(&users).Error; err != nil {
		zlog.Error(err.Error())
		return nil, err
	}
	return users, nil
}

func Find(out interface{}, query interface{}, args ...interface{}) error {
	if err := GormDB.Where(query, args...).Find(out).Error; err != nil {
		zlog.Error(err.Error())
		return err
	}
	return nil
}
