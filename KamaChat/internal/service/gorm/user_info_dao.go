package gorm

import (
	"errors"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/model"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/zlog"

	"gorm.io/gorm"
)

type userInfoDao struct {
}

var UserInfoDao = new(userInfoDao)

func (u *userInfoDao) GetUserInfoByTelephone(user *model.UserInfo, telephone string) (string, int) {
	res := dao.GormDB.First(&user, "telephone = ?", telephone)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "用户不存在，请注册"
			zlog.Error(message)
			return message, -2
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	return "", 0
}

func (u *userInfoDao) SaveUser(user *model.UserInfo) int {
	res := dao.GormDB.Save(user)
	if res.Error != nil {
		zlog.Error(res.Error.Error())
		return -1
	}
	return 0
}
