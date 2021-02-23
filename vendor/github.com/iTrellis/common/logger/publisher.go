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

// NewPublisher new a publisher
func NewPublisher() Publisher {
	return &publisher{
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

type publisher struct {
	prefixes   []interface{}
	hasCaller  bool
	eventGroup event.SubscriberGroup
}

// Debug 调试
func (p *publisher) Debug(fields ...interface{}) {
	p.publishLog(DebugLevel, fields...)
}

// Debugf 调试
func (p *publisher) Debugf(msg string, fields ...interface{}) {
	p.Debug(fmt.Sprintf(msg, fields...))
}

// Info 信息
func (p *publisher) Info(fields ...interface{}) {
	p.publishLog(InfoLevel, fields...)
}

// Infof 信息
func (p *publisher) Infof(msg string, fields ...interface{}) {
	p.Info(fmt.Sprintf(msg, fields...))
}

// Warn 警告
func (p *publisher) Warn(fields ...interface{}) {
	p.publishLog(WarnLevel, fields...)
}

// Warnf 警告
func (p *publisher) Warnf(msg string, fields ...interface{}) {
	p.Warn(fmt.Sprintf(msg, fields...))
}

// Error 错误
func (p *publisher) Error(fields ...interface{}) {
	p.publishLog(ErrorLevel, fields...)
}

// Errorf 错误
func (p *publisher) Errorf(msg string, fields ...interface{}) {
	p.Error(fmt.Sprintf(msg, fields...))
}

// Critical 严重的
func (p *publisher) Critical(fields ...interface{}) {
	p.publishLog(CriticalLevel, fields...)
}

// Criticalf 严重的
func (p *publisher) Criticalf(msg string, fields ...interface{}) {
	p.Critical(fmt.Sprintf(msg, fields...))
}

// Panic panic
func (p *publisher) Panic(fields ...interface{}) {
	p.publishLog(PanicLevel, fields...)
}

// Panicf panic
func (p *publisher) Panicf(msg string, fields ...interface{}) {
	p.Panic(fmt.Sprintf(msg, fields...))
}

// SetLevel set looger's level
func (p *publisher) SetLevel(lvl Level) {
	p.Publish(lvl)
}

// Log 打印
func (p *publisher) Log(kvs ...interface{}) error {
	p.publishLog(InfoLevel, kvs...)
	return nil
}

// With 异常
func (p *publisher) With(params ...interface{}) Publisher {
	if len(params) == 0 {
		return p
	}
	newPrefixes := append(p.prefixes, params...)

	return &publisher{
		prefixes:   newPrefixes[:len(newPrefixes):len(newPrefixes)],
		hasCaller:  p.hasCaller || containsCaller(newPrefixes),
		eventGroup: p.eventGroup,
	}
}

// WithPrefix 加载前缀
func (p *publisher) WithPrefix(prefixes ...interface{}) Publisher {
	if len(prefixes) == 0 {
		return p
	}
	newPrefixes := append(prefixes, p.prefixes...)
	return &publisher{
		prefixes:   newPrefixes,
		hasCaller:  p.hasCaller || containsCaller(newPrefixes),
		eventGroup: p.eventGroup,
	}
}

func (p *publisher) publishLog(lvl Level, kvs ...interface{}) {
	if len(kvs) == 0 {
		return
	}
	prifixes := p.prefixes
	if p.hasCaller {
		if len(prifixes) != 0 {
			bindCallers(prifixes)
		}
	}
	p.Publish(&Event{
		Time:   time.Now(),
		Level:  lvl,
		Fields: append(prifixes, kvs...),
	})
}

func (p *publisher) Subscriber(s interface{}) (event.Subscriber, error) {
	return p.eventGroup.Subscriber(s)
}

func (p *publisher) RemoveSubscriber(ids ...string) error {
	return p.eventGroup.RemoveSubscriber(ids...)
}

func (p *publisher) Publish(values ...interface{}) {
	p.eventGroup.Publish(values...)
}

func (p *publisher) ClearSubscribers() {
	p.eventGroup.ClearSubscribers()
}
