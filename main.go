package main

import (
	"gitee.com/tzxhy/dlna-home/initial"
	"gitee.com/tzxhy/dlna-home/routers"
)

func main() {
	initial.InitAll()
	api := routers.InitRouter()

	api.Run(":8082")

}
