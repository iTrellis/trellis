// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package txorm

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/go-trellis/config"
)

// defines
const (
	DefaultDatabase = "go-trellis::txorm::default"
)

// GetMysqlDSNFromConfig get mysql dsn from gogap config
func GetMysqlDSNFromConfig(name string, conf config.Config) string {
	if name == "" {
		panic("database's name not exist")
	}

	dsn := mysql.Config{
		DBName:  name,
		Net:     "tcp",
		Timeout: conf.GetTimeDuration("timeout", time.Second*5),

		User:   conf.GetString("user", "root"),
		Passwd: conf.GetString("password", ""),
		Addr:   fmt.Sprintf("%s:%d", conf.GetString("host", "localhost"), conf.GetInt("port", 3306)),
		Params: map[string]string{
			"charset":              conf.GetString("charset", "utf8"),
			"parseTime":            conf.GetString("parseTime", "True"),
			"loc":                  conf.GetString("location", "Local"),
			"allowNativePasswords": conf.GetString("allowNativePasswords", "true"),
		},
	}
	return dsn.FormatDSN()
}
