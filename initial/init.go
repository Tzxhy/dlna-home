package initial

import (
	"gitee.com/tzxhy/dlna-home/constants"
	"gitee.com/tzxhy/dlna-home/models"
	"gitee.com/tzxhy/dlna-home/utils"
)

func InitAll() {
	utils.ShowAppInfo()
	models.InitSqlite3()
	utils.MakeSurePathExists(constants.UPLOAD_PATH)
	InitStatic()
}
