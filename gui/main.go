// Package main provides various examples of Fyne API capabilities.
package main

import (
	"log"
	"net/url"
	"time"

	"gitee.com/tzxhy/dlna-home/constants"
	"gitee.com/tzxhy/dlna-home/gui/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gitee.com/tzxhy/dlna-home/initial"
	"gitee.com/tzxhy/dlna-home/routers"
)

func startServer() {
	log.SetPrefix("from gui:")
	initial.InitAll()
	api := routers.InitRouter()

	api.Run(":8111")
}

func main() {
	a := app.NewWithID("com.gitee.dlnahome")
	constants.StorageRoot = a.Storage().RootURI().String()

	a.Settings().SetTheme(&theme.MyTheme{})

	logLifecycle(a)
	w := a.NewWindow("dlna home controller")

	w.SetMaster()

	w.SetContent(makeNav(a))

	w.Resize(fyne.NewSize(640, 460))
	go startServer()
	// go delayOpen(a)
	w.ShowAndRun()
}
func delayOpen(a fyne.App, delaySecond uint8) {
	<-time.After(time.Duration(delaySecond) * time.Second)
	u, _ := url.Parse("http://localhost:8111")
	a.OpenURL(u)
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeNav(a fyne.App) fyne.CanvasObject {
	show := widget.NewButton("点我去配置", func() {
		delayOpen(a, 0)
	})
	center := layout.NewCenterLayout()
	content := container.New(center, show)
	return container.NewBorder(nil, nil, nil, nil, content)
}
