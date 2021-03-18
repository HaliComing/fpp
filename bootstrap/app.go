package bootstrap

import (
	"encoding/json"
	"fmt"

	"github.com/HaliComing/fpp/pkg/conf"
	"github.com/HaliComing/fpp/pkg/request"
	"github.com/HaliComing/fpp/pkg/util"
	"github.com/hashicorp/go-version"
)

// InitApplication 初始化应用常量
func InitApplication() {
	fmt.Print(`
fpp  V` + conf.Version + `
=============================

`)
	//go CheckUpdate()
}

type GitHubRelease struct {
	URL  string `json:"html_url"`
	Name string `json:"name"`
	Tag  string `json:"tag_name"`
}

// CheckUpdate 检查更新
func CheckUpdate() {
	client := request.HTTPClient{}
	res, err := client.Request("GET", "https://api.github.com/repos/halicoming/fpp/releases", nil).GetResponse()
	if err != nil {
		util.Log().Warning("[CheckUpdate] Check update failed, Error = %s", err)
		return
	}

	var list []GitHubRelease
	if err := json.Unmarshal([]byte(res), &list); err != nil {
		util.Log().Warning("[CheckUpdate] Check update failed, Error = %s", err)
		return
	}

	if len(list) > 0 {
		present, err1 := version.NewVersion(conf.Version)
		latest, err2 := version.NewVersion(list[0].Tag)
		if err1 == nil && err2 == nil && latest.GreaterThan(present) {
			util.Log().Info("[CheckUpdate] A new version is available, Version = [%s], Download url = %s", list[0].Name, list[0].URL)
		}
	}

}
