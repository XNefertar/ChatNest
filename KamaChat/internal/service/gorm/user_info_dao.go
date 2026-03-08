package gorm

import (
	"errors"
	"fmt"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/zlog"

	"gorm.io/gorm"
)

type userInfoDao struct {
}

var UserInfoDao = new(userInfoDao)

// IsRecordNotFound 检查错误是否是记录未找到
func (u *userInfoDao) IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (u *userInfoDao) GetUserInfo(identifier interface{}) (*model.UserInfo, error) {
	var user model.UserInfo
	db := dao.GormDB

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

func (u *userInfoDao) Create(value interface{}) error {
	if err := dao.GormDB.Create(value).Error; err != nil {
		zlog.Error(err.Error())
		return err
	}
	return nil
}

func (u *userInfoDao) Save(identifier interface{}) error {
	if err := dao.GormDB.Save(identifier).Error; err != nil {
		zlog.Error(err.Error())
		return err
	}
	return nil
}

func (u *userInfoDao) GetAllUsersExcept(identifier interface{}) ([]model.UserInfo, error) {
	var users []model.UserInfo
	db := dao.GormDB
	switch v := identifier.(type) {
	case UserOwnerID:
		db = db.Where("uuid != ?", string(v))
	}

	if err := db.Unscoped().Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userInfoDao) GetUsersByUUIDs(uuidList []string) ([]model.UserInfo, error) {
	var users []model.UserInfo
	if err := dao.GormDB.Model(model.UserInfo{}).Where("uuid in (?)", uuidList).Find(&users).Error; err != nil {
		zlog.Error(err.Error())
		return nil, err
	}
	return users, nil
}

func (u *userInfoDao) Find(out interface{}, query interface{}, args ...interface{}) error {
	if err := dao.GormDB.Where(query, args...).Find(out).Error; err != nil {
		zlog.Error(err.Error())
		return err
	}
	return nil
}
