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

func GetPlayList() (*[]PlayListItem, error) {
	var playLists []PlayListItem
	err := DB.Find(&playLists).Error
	return &playLists, err
}

type ListItemParam struct {
	Name string
	Url  string
}

func RenamePlayList(pid, name string) {
	DB.Where(&PlayListItem{
		Pid: pid,
	}).Updates(&PlayListItem{
		Name: name,
	})
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

func DeletePlayList(pid string) bool {
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 删除原有绑定
		err := DB.Where(&AudioItem{
			PlayListItemID: pid,
		}).Delete(&AudioItem{}).Error
		if err != nil {
			return err
		}

		err = DB.Where(&PlayListItem{
			Pid: pid,
		}).Delete(&PlayListItem{}).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err == nil
}

const ALL_ITEMS_PID = "ALL"

func GetPlayListItems(pid string) *[]AudioItem {
	if pid == ALL_ITEMS_PID {
		return GetAllPlayListItems()
	}
	var items []AudioItem
	DB.Where(&AudioItem{
		PlayListItemID: pid,
	}).Find(&items).Order("create_date desc")
	return &items
}
func GetAllPlayListItems() *[]AudioItem {
	var items []AudioItem
	DB.Find(&items).Order("create_date desc")
	return &items
}

func DeleteSingleResource(aid string) {
	DB.Where(&AudioItem{
		Aid: aid,
	}).Delete(&AudioItem{})
}
