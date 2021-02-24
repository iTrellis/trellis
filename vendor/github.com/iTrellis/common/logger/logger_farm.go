/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package logger

import (
	"fmt"

	"github.com/google/uuid"
)

// WithPrefix with prefix
func WithPrefix(logger Logger, prefixes ...interface{}) Logger {
	return &context{
		id: uuid.NewString(),

		logger:   logger,
		prefixes: prefixes,
	}
}

type context struct {
	id       string
	logger   Logger
	prefixes []interface{}
}

// Debug 调试
func (p *context) Debug(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.logger.Debug(logs...)
}

// Debugf 调试
func (p *context) Debugf(msg string, kvs ...interface{}) {
	logs := append(p.prefixes, "msg", fmt.Sprintf(msg, kvs...))
	p.logger.Debug(logs...)
}

// Info 信息
func (p *context) Info(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.logger.Info(logs...)
}

// Infof 信息
func (p *context) Infof(msg string, kvs ...interface{}) {
	logs := append(p.prefixes, "msg", fmt.Sprintf(msg, kvs...))
	p.logger.Info(logs...)
}

// Warn 警告
func (p *context) Warn(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.logger.Warn(logs...)
}

// Warnf 警告
func (p *context) Warnf(msg string, kvs ...interface{}) {
	logs := append(p.prefixes, "msg", fmt.Sprintf(msg, kvs...))
	p.logger.Warn(logs...)
}

// Error 错误
func (p *context) Error(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.logger.Error(logs...)
}

// Errorf 错误
func (p *context) Errorf(msg string, kvs ...interface{}) {
	logs := append(p.prefixes, "msg", fmt.Sprintf(msg, kvs...))
	p.logger.Error(logs...)
}

// Critical 严重的
func (p *context) Critical(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.logger.Critical(logs...)
}

// Criticalf 严重的
func (p *context) Criticalf(msg string, kvs ...interface{}) {
	logs := append(p.prefixes, "msg", fmt.Sprintf(msg, kvs...))
	p.logger.Critical(logs...)
}

// Panic panic
func (p *context) Panic(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.logger.Panic(logs...)
}

// Panicf panic
func (p *context) Panicf(msg string, kvs ...interface{}) {
	logs := append(p.prefixes, "msg", fmt.Sprintf(msg, kvs...))
	p.logger.Panic(logs...)
}

func (p *context) GetID() string {
	return p.logger.GetID()
}

func (p *context) Log(kvs ...interface{}) error {
	logs := append(p.prefixes, kvs...)
	return p.logger.Log(logs...)
}

func (p *context) Publish(vals ...interface{}) error {
	for _, v := range vals {
		switch t := v.(type) {
		case Event:
			t.Fields = append(p.prefixes, t.Fields...)
			p.logger.Publish(&t)
		case *Event:
			evt := *t
			evt.Fields = append(p.prefixes, evt.Fields...)
			p.logger.Publish(&evt)
		case Level:
			p.logger.Publish(t)
		}
	}
	return nil
}

func (p *context) SetLevel(lvl Level) {
	p.logger.SetLevel(lvl)
}

func (p *context) Stop() {
	p.logger.Stop()
}
