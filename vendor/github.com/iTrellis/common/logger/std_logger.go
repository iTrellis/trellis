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
	"io"
	"os"
	"reflect"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

// STDOptions std options
type STDOptions struct {
	level  Level
	writer io.Writer

	publisher Publisher
}

// STDPublisher set publisher
func STDPublisher(pub Publisher) STDOption {
	return func(c *STDOptions) {
		c.publisher = pub
	}
}

type stdLogger struct {
	id      string
	options STDOptions
	logger  kitlog.Logger

	prefixes []interface{}
}

// NewStdLogger new std logger
func NewStdLogger(opts ...STDOption) Logger {

	l := &stdLogger{
		id: uuid.NewString(),
	}

	for _, o := range opts {
		o(&l.options)
	}

	if l.options.writer == nil {
		l.options.writer = os.Stdout
	}

	l.logger = kitlog.NewLogfmtLogger(l.options.writer)

	if l.options.publisher != nil {
		l.options.publisher.Subscriber(l)
	}
	return l
}

func (p *stdLogger) GetID() string {
	return p.id
}

func (p *stdLogger) Publish(evts ...interface{}) error {
	for _, evt := range evts {
		switch t := evt.(type) {
		case Event:
			p.logEvent(&t)
		case *Event:
			p.logEvent(t)
		case Level:
			p.options.level = t
		default:
			return fmt.Errorf("unsupported event type: %+v", reflect.TypeOf(evt))
		}
	}
	return nil
}

func (p *stdLogger) SetLevel(lvl Level) {
	p.options.level = lvl
}

func (p *stdLogger) Stop() {
	p.options.publisher = nil
}

func (p *stdLogger) Log(kvs ...interface{}) error {
	return p.logger.Log(kvs...)
}

func (p *stdLogger) pubLog(level Level, kvs ...interface{}) {
	p.Publish(&Event{
		Time:   time.Now(),
		Level:  level,
		Fields: kvs,
	})
}

func (p *stdLogger) logEvent(log *Event) error {
	if log.Level < p.options.level {
		return nil
	}
	return p.Log(genLogs(log)...)
}

// Debug 调试
func (p *stdLogger) Debug(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.pubLog(DebugLevel, logs...)
}

// Debugf 调试
func (p *stdLogger) Debugf(msg string, kvs ...interface{}) {
	p.Debug("msg", fmt.Sprintf(msg, kvs...))
}

// Info 信息
func (p *stdLogger) Info(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.pubLog(InfoLevel, logs...)
}

// Infof 信息
func (p *stdLogger) Infof(msg string, kvs ...interface{}) {
	p.Info("msg", fmt.Sprintf(msg, kvs...))
}

// Warn 警告
func (p *stdLogger) Warn(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.pubLog(WarnLevel, logs...)
}

// Warnf 警告
func (p *stdLogger) Warnf(msg string, kvs ...interface{}) {
	p.Warn("msg", fmt.Sprintf(msg, kvs...))
}

// Error 错误
func (p *stdLogger) Error(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.pubLog(ErrorLevel, logs...)
}

// Errorf 错误
func (p *stdLogger) Errorf(msg string, kvs ...interface{}) {
	p.Error("msg", fmt.Sprintf(msg, kvs...))
}

// Critical 严重的
func (p *stdLogger) Critical(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.pubLog(CriticalLevel, logs...)
}

// Criticalf 严重的
func (p *stdLogger) Criticalf(msg string, kvs ...interface{}) {
	p.Critical("msg", fmt.Sprintf(msg, kvs...))
}

// Panic panic
func (p *stdLogger) Panic(kvs ...interface{}) {
	logs := append(p.prefixes, kvs...)
	p.pubLog(PanicLevel, logs...)
}

// Panicf panic
func (p *stdLogger) Panicf(msg string, kvs ...interface{}) {
	p.Panic("msg", fmt.Sprintf(msg, kvs...))
}

func (p *stdLogger) WithPrefix(kvs ...interface{}) Logger {
	return &stdLogger{
		id:       uuid.NewString(),
		options:  p.options,
		prefixes: append(kvs, p.prefixes...),
		logger:   p.logger,
	}
}
