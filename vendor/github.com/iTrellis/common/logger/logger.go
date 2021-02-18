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
	"time"

	"github.com/iTrellis/common/event"
)

// Logger 日志对象
type Logger interface {
	Debug(msg string, fields ...interface{})
	Debugf(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Infof(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Warnf(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Errorf(msg string, fields ...interface{})
	Critical(msg string, fields ...interface{})
	Criticalf(msg string, fields ...interface{})
	With(params ...interface{}) Logger
	WithPrefix(prefixes ...interface{}) Logger

	event.SubscriberGroup
}

// NewLogger 获取日志实例
func NewLogger() Logger {
	return &context{
		eventGroup: event.NewSubscriberGroup(),
	}
}

// Caller fileds function
type Caller func() interface{}

func containsCaller(fileds []interface{}) bool {
	for i := 0; i < len(fileds); i++ {
		switch fileds[i].(type) {
		case Caller, func() interface{}:
			return true
		}
	}
	return false
}

func bindCallers(fileds []interface{}) {
	for i := 0; i < len(fileds); i++ {
		switch fn := fileds[i].(type) {
		case Caller:
			fileds[i] = fn()
		case func() interface{}:
			fileds[i] = fn()
		}
	}
}

type context struct {
	prefixes   []interface{}
	hasCaller  bool
	eventGroup event.SubscriberGroup
}

// Debug 调试
func (p *context) Debug(msg string, fields ...interface{}) {
	p.publishLog(DebugLevel, msg, fields...)
}

// Debugf 调试
func (p *context) Debugf(msg string, fields ...interface{}) {
	p.Debug(fmt.Sprintf(msg, fields...))
}

// Info 信息
func (p *context) Info(msg string, fields ...interface{}) {
	p.publishLog(InfoLevel, msg, fields...)
}

// Infof 信息
func (p *context) Infof(msg string, fields ...interface{}) {
	p.Info(fmt.Sprintf(msg, fields...))
}

// Warn 警告
func (p *context) Warn(msg string, fields ...interface{}) {
	p.publishLog(WarnLevel, msg, fields...)
}

// Warnf 警告
func (p *context) Warnf(msg string, fields ...interface{}) {
	p.Warn(fmt.Sprintf(msg, fields...))
}

// Error 错误
func (p *context) Error(msg string, fields ...interface{}) {
	p.publishLog(ErrorLevel, msg, fields...)
}

// Errorf 错误
func (p *context) Errorf(msg string, fields ...interface{}) {
	p.Error(fmt.Sprintf(msg, fields...))
}

// Critical 严重的
func (p *context) Critical(msg string, fields ...interface{}) {
	p.publishLog(CriticalLevel, msg, fields...)
}

// Criticalf 严重的
func (p *context) Criticalf(msg string, fields ...interface{}) {
	p.Critical(fmt.Sprintf(msg, fields...))
}

// With 异常
func (p *context) With(params ...interface{}) Logger {
	if len(params) == 0 {
		return p
	}
	newPrefixes := append(p.prefixes, params...)

	return &context{
		prefixes:   newPrefixes[:len(newPrefixes):len(newPrefixes)],
		hasCaller:  p.hasCaller || containsCaller(newPrefixes),
		eventGroup: p.eventGroup,
	}
}

// WithPrefix 加载前缀
func (p *context) WithPrefix(prefixes ...interface{}) Logger {
	if len(prefixes) == 0 {
		return p
	}
	newPrefixes := append(prefixes, p.prefixes...)
	return &context{
		prefixes:   newPrefixes,
		hasCaller:  p.hasCaller || containsCaller(newPrefixes),
		eventGroup: p.eventGroup,
	}
}

func (p *context) publishLog(lvl Level, msg string, fields ...interface{}) {
	if len(msg) == 0 {
		panic("message should not be empty")
	}
	prifixes := p.prefixes
	if p.hasCaller {
		if len(prifixes) != 0 {
			bindCallers(prifixes)
		}
	}
	p.Publish(&Event{
		Time:     time.Now(),
		Level:    lvl,
		Prefixes: prifixes,
		Fields:   append([]interface{}{msg}, fields...)})
}

func (p *context) Subscriber(s interface{}) (event.Subscriber, error) {
	return p.eventGroup.Subscriber(s)
}

func (p *context) RemoveSubscriber(ids ...string) error {
	return p.eventGroup.RemoveSubscriber(ids...)
}

func (p *context) Publish(values ...interface{}) {
	p.eventGroup.Publish(values...)
}

func (p *context) ClearSubscribers() {
	p.eventGroup.ClearSubscribers()
}
