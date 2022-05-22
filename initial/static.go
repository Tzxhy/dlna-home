package initial

import (
	"net/http"

	_ "gitee.com/tzxhy/dlna-home/statik"
	"github.com/rakyll/statik/fs"
)

var StatikFS http.FileSystem

func InitStatic() {
	StatikFS, _ = fs.New()
}
