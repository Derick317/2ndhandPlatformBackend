package service

import (
	"errors"
	"secondHand/backend"
	"secondHand/model"
	"secondHand/util"

	"gorm.io/gorm"
)

func Search(keywords []string, tag model.TagType) ([]model.Item, error) {
	items := []model.Item{}
	qKeys := []string{"status"}
	qTargets := []interface{}{model.Available}
	qEqual := []bool{true}
	if len(keywords) > 0 {
		qKeys = append(qKeys, "title")
		qTargets = append(qTargets, keywords)
		qEqual = append(qEqual, false)
	}
	if tag != model.TagAll {
		qKeys = append(qKeys, "tag")
		qTargets = append(qTargets, tag)
		qEqual = append(qEqual, true)
	}
	if err := backend.ReadFromDBEqualOrInclude(&items, qKeys, qTargets, qEqual,
		false); errors.Is(err, gorm.ErrRecordNotFound) {
		return []model.Item{}, util.ErrItemNotFound
	} else if err != nil {
		return []model.Item{}, err
	}
	return items, nil
}
