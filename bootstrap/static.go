package bootstrap

import (
	"github.com/HaliComing/fpp/pkg/util"
	_ "github.com/HaliComing/fpp/statik"
	"github.com/gin-contrib/static"
	"github.com/rakyll/statik/fs"
	"net/http"
)

type GinFS struct {
	FS http.FileSystem
}

// StaticFS 内置静态文件资源
var StaticFS static.ServeFileSystem

// Open 打开文件
func (b *GinFS) Open(name string) (http.File, error) {
	return b.FS.Open(name)
}

// Exists 文件是否存在
func (b *GinFS) Exists(prefix string, filepath string) bool {

	if _, err := b.FS.Open(filepath); err != nil {
		return false
	}
	return true

}

// InitStatic 初始化静态资源文件
func InitStatic() {
	var err error

	StaticFS = &GinFS{}
	StaticFS.(*GinFS).FS, err = fs.New()
	if err != nil {
		util.Log().Panic("[Static] Unable to initialize static resource, Error = %s", err)
	}
}
