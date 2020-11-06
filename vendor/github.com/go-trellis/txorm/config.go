// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package txorm

import (
	"fmt"
	"sync"

	"github.com/go-trellis/config"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var locker = &sync.Mutex{}

// NewEnginesFromFile initial xorm engine from file
func NewEnginesFromFile(file string) (map[string]*xorm.Engine, error) {
	conf, err := config.NewConfigOptions(config.OptionFile(file))
	if err != nil {
		return nil, err
	}
	return NewEnginesFromConfig(conf)
}

// NewEnginesFromConfig initial xorm engine from config
func NewEnginesFromConfig(conf config.Config) (engines map[string]*xorm.Engine, err error) {

	es := make(map[string]*xorm.Engine)
	if conf == nil {
		return nil, fmt.Errorf("config is nil")
	}

	locker.Lock()
	defer locker.Unlock()

	for _, databaseName := range conf.GetKeys() {
		var _engine *xorm.Engine
		databaseConf := conf.GetValuesConfig(databaseName)
		switch driver := databaseConf.GetString("driver", "mysql"); driver {
		case driver:

			_engine, err = xorm.NewEngine(driver, GetMysqlDSNFromConfig(databaseName, databaseConf))
			if err != nil {
				return nil, err
			}

			_engine.SetMaxIdleConns(conf.GetInt(databaseName+".max_idle_conns", 10))

			_engine.SetMaxOpenConns(conf.GetInt(databaseName+".max_open_conns", 100))

			_engine.ShowSQL(conf.GetBoolean(databaseName + ".show_sql"))

		default:
			return nil, fmt.Errorf("unsupported driver: %s", driver)
		}

		_engine.Logger().SetLevel(log.LogLevel(conf.GetInt(databaseName+".log_level", 0)))

		if _isD := conf.GetBoolean(databaseName + ".is_default"); _isD {
			es[DefaultDatabase] = _engine
		}

		es[databaseName] = _engine
	}

	return es, nil
}
