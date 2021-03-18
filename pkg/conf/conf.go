package conf

import (
	"github.com/HaliComing/fpp/pkg/util"
	"github.com/go-ini/ini"
	"gopkg.in/go-playground/validator.v9"
)

// database 数据库
type database struct {
	Type   string
	DBFile string
}

// system 系统通用配置
type system struct {
	Listen             string `validate:"required"`
	Token              string `validate:"required"`
	Debug              bool
	NumberOfThreads    int
	ExtractionInterval int
	CheckInterval      int
}

type ssl struct {
	CertPath string `validate:"omitempty,required"`
	KeyPath  string `validate:"omitempty,required"`
	Listen   string `validate:"required"`
}

var cfg *ini.File

const defaultConf = `[System]
Listen = :9826
NumberOfThreads = 50
Token = {Token}
`

// Init 初始化配置文件
func Init(path string) {
	var err error

	if path == "" || !util.Exists(path) {
		// 创建初始配置文件
		confContent := util.Replace(map[string]string{
			"{Token}": util.RandStringRunes(32),
		}, defaultConf)
		f, err := util.CreatNestedFile(path)
		if err != nil {
			util.Log().Panic("[Conf] Unable to create configuration file, Error = %s", err)
		}

		// 写入配置文件
		_, err = f.WriteString(confContent)
		if err != nil {
			util.Log().Panic("[Conf] Unable to write to configuration file, Error = %s", err)
		}

		f.Close()
	}

	cfg, err = ini.Load(path)
	if err != nil {
		util.Log().Panic("[Conf] The configuration file could not be resolved, Path = '%s', Error = %s", path, err)
	}

	sections := map[string]interface{}{
		"System":   SystemConfig,
		"SSL":      SSLConfig,
		"Database": DatabaseConfig,
	}
	for sectionName, sectionStruct := range sections {
		err = mapSection(sectionName, sectionStruct)
		if err != nil {
			util.Log().Panic("[Conf] %s Section parsing failed, Error = %s", sectionName, err)
		}
	}

	// 重设log等级
	if !SystemConfig.Debug {
		util.Level = util.LevelInformational
		util.GloablLogger = nil
		util.Log()
	}

}

// mapSection 将配置文件的 Section 映射到结构体上
func mapSection(section string, confStruct interface{}) error {
	err := cfg.Section(section).MapTo(confStruct)
	if err != nil {
		return err
	}

	// 验证合法性
	validate := validator.New()
	err = validate.Struct(confStruct)
	if err != nil {
		return err
	}

	return nil
}
