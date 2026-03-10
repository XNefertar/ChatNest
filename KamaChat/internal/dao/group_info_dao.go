package dao

import (
	"fmt"
	"kama_chat_server/internal/model"
)

type GroupUUID string

func GetGroupInfo(identifier interface{}) (*model.GroupInfo, error) {
	var group model.GroupInfo
	db := GormDB

	switch v := identifier.(type) {
	case GroupUUID:
		db = db.Where("uuid = ?", string(v))
	default:
		return nil, fmt.Errorf("unsupported identifier type: %T", identifier)
	}

	if err := db.First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}
