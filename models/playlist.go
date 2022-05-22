package models

import (
	"gitee.com/tzxhy/dlna-home/utils"
	"gorm.io/gorm"
)

type PlayListItem struct {
	Pid        string `json:"pid" gorm:"primaryKey;type:string;"`
	Name       string `json:"name" gorm:"type:string not null;"`
	CreateDate uint64 `json:"create_date" gorm:"autoUpdateTime:milli"`
}

type AudioItem struct {
	PlayListItem   PlayListItem `json:"-" gorm:"references:Pid"`
	PlayListItemID string       `json:"pid" gorm:"type:string not null;"`
	Aid            string       `json:"aid" gorm:"type: string not null;"`
	Url            string       `json:"url" gorm:"type: string not null;"`
	Name           string       `json:"name" gorm:"type: string not null;"`
	CreateDate     uint64       `json:"create_date" gorm:"autoUpdateTime:milli"`
}

func CreatePlayList(name string) (string, error) {
	pid := utils.GenerateRid()
	err := DB.Create(&PlayListItem{
		Pid:  pid,
		Name: name,
	}).Error
	if err != nil {
		return "", err
	}
	return pid, nil
}

type ListItemParam struct {
	Name string
	Url  string
}

func SetPlayList(pid string, list []ListItemParam) bool {
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有绑定
		err := DB.Where(&AudioItem{
			PlayListItemID: pid,
		}).Delete(&AudioItem{}).Error
		if err != nil {
			return err
		}

		// 增加现在
		newList := utils.Map(&list, func(item ListItemParam) AudioItem {
			aid := utils.GenerateRid()
			return AudioItem{
				Aid:            aid,
				PlayListItemID: pid,
				Name:           item.Name,
				Url:            item.Url,
			}
		})

		err = DB.Create(&newList).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err == nil

}

func GetPlayListItems(pid string) *[]AudioItem {
	var items []AudioItem
	DB.Where(&AudioItem{
		PlayListItemID: pid,
	}).Find(&items)
	return &items
}
