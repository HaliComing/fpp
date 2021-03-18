package conf

// DatabaseConfig 数据库配置
var DatabaseConfig = &database{
	Type:   "sqlite3",
	DBFile: "fpp.db",
}

// SystemConfig 系统公用配置
var SystemConfig = &system{
	Debug:              false,
	Listen:             ":9826",
	NumberOfThreads:    50,
	ExtractionInterval: 60,
	CheckInterval:      15,
}

var SSLConfig = &ssl{
	Listen:   ":443",
	CertPath: "",
	KeyPath:  "",
}
