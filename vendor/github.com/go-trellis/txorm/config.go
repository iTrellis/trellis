// GNU GPL v3 License
// Copyright (c) 2019 github.com:go-trellis

package txorm

import (
	"fmt"
	"sync"

	"github.com/go-trellis/config"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

var locker = &sync.Mutex{}

// NewEnginesFromFile initial xorm engine from file
func NewEnginesFromFile(file string) (map[string]*xorm.Engine, error) {
	conf, err := config.NewConfigOptions(config.OptionFile(file))
	if err != nil {
		return nil, err
	}
	return NewEnginesFromConfig(conf, "mysql")
}

// NewEnginesFromConfig initial xorm engine from config
func NewEnginesFromConfig(conf config.Config, name string) (map[string]*xorm.Engine, error) {

	engines := make(map[string]*xorm.Engine)

	locker.Lock()
	defer locker.Unlock()

	cfg := conf.GetValuesConfig(name)
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	for _, databaseName := range cfg.GetKeys() {
		_engine, err := xorm.NewEngine("mysql", GetMysqlDSNFromConfig(databaseName, cfg.GetValuesConfig(databaseName)))
		if err != nil {
			return nil, err
		}

		_engine.SetMaxIdleConns(cfg.GetInt(databaseName+".max_idle_conns", 10))

		_engine.SetMaxOpenConns(cfg.GetInt(databaseName+".max_open_conns", 100))

		_engine.ShowSQL(cfg.GetBoolean(databaseName + ".show_sql"))

		_engine.Logger().SetLevel(core.LogLevel(cfg.GetInt(databaseName+".log_level", 0)))

		if _isD := cfg.GetBoolean(databaseName + ".is_default"); _isD {
			engines[DefaultDatabase] = _engine
		}

		engines[databaseName] = _engine
	}

	return engines, nil
}
