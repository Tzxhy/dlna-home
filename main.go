package main

import (
	"gitee.com/tzxhy/dlna-home/initial"
	"gitee.com/tzxhy/dlna-home/routers"
)

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func main() {
	initial.InitAll()
	// user := models.GetUserById(1)
	// fmt.Print(user)
	api := routers.InitRouter()

	api.Run(":8081")
	// ex, _ := os.Executable()
	// fmt.Print(ex)
	// fmt.Print(path.Join(path.Dir(ex)))
}
