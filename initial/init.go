package initial

import (
	"gitee.com/tzxhy/dlna-home/models"
	"gitee.com/tzxhy/dlna-home/utils"
)

func InitAll() {
	utils.ShowAppInfo()
	models.InitSqlite3()
	InitStatic()
}
