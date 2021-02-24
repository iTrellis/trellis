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
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	id string

	options LogrusOptions
	logger  logrus.FieldLogger
}

// LogrusOptions options
type LogrusOptions struct {
	level Level
}

var mapToLogrusLevel = map[Level]logrus.Level{
	TraceLevel:    logrus.TraceLevel,
	DebugLevel:    logrus.DebugLevel,
	InfoLevel:     logrus.InfoLevel,
	WarnLevel:     logrus.WarnLevel,
	ErrorLevel:    logrus.ErrorLevel,
	CriticalLevel: logrus.FatalLevel,
	PanicLevel:    logrus.PanicLevel,
}

// NewLogrusLogger logrus logger
func NewLogrusLogger(l logrus.FieldLogger, opts ...LogrusOption) Logger {

	log := &logrusLogger{
		id: uuid.NewString(),
	}

	for _, o := range opts {
		o(&log.options)
	}

	return log
}

func (p *logrusLogger) logEvent(evt *Event) {
	if evt.Level < p.options.level {
		return
	}

	vals := genLogs(evt)

	fields := logrus.Fields{}
	for i := 0; i < len(vals); i += 2 {
		key := toString(vals[i])
		fields[key] = vals[i+1]
	}

	switch mapToLogrusLevel[evt.Level] {
	case logrus.TraceLevel:
		p.logger.WithFields(fields).Trace()
	case logrus.DebugLevel:
		p.logger.WithFields(fields).Debug()
	case logrus.InfoLevel:
		p.logger.WithFields(fields).Info()
	case logrus.WarnLevel:
		p.logger.WithFields(fields).Warn()
	case logrus.ErrorLevel:
		p.logger.WithFields(fields).Error()
	case logrus.FatalLevel:
		p.logger.WithFields(fields).Fatal()
	case logrus.PanicLevel:
		p.logger.WithFields(fields).Panic()
	}

	return
}

func (p *logrusLogger) Publish(evts ...interface{}) error {
	for _, evt := range evts {
		switch eType := evt.(type) {
		case Event:
			p.logEvent(&eType)
		case *Event:
			p.logEvent(eType)
		case Level:
			p.options.level = eType
		default:
			return fmt.Errorf("unsupported event type: %s", reflect.TypeOf(evt).Name())
		}
	}
	return nil
}

func (p *logrusLogger) Log(kvs ...interface{}) error {
	p.Publish(Event{
		Time:   time.Now(),
		Level:  InfoLevel,
		Fields: kvs,
	})
	return nil
}

func (p *logrusLogger) pubLog(level Level, kvs ...interface{}) {
	p.Publish(&Event{
		Time:   time.Now(),
		Level:  level,
		Fields: kvs,
	})
}

func (p *logrusLogger) GetID() string {
	return p.id
}

func (p *logrusLogger) Stop() {}

func (p *logrusLogger) SetLevel(lvl Level) {
	p.options.level = lvl
}

// Debug 调试
func (p *logrusLogger) Debug(kvs ...interface{}) {
	p.pubLog(DebugLevel, kvs...)
}

// Debugf 调试
func (p *logrusLogger) Debugf(msg string, kvs ...interface{}) {
	p.Debug("msg", fmt.Sprintf(msg, kvs...))
}

// Info 信息
func (p *logrusLogger) Info(kvs ...interface{}) {
	p.pubLog(InfoLevel, kvs...)
}

// Infof 信息
func (p *logrusLogger) Infof(msg string, kvs ...interface{}) {
	p.Info("msg", fmt.Sprintf(msg, kvs...))
}

// Warn 警告
func (p *logrusLogger) Warn(kvs ...interface{}) {
	p.pubLog(WarnLevel, kvs...)
}

// Warnf 警告
func (p *logrusLogger) Warnf(msg string, kvs ...interface{}) {
	p.Warn("msg", fmt.Sprintf(msg, kvs...))
}

// Error 错误
func (p *logrusLogger) Error(kvs ...interface{}) {
	p.pubLog(ErrorLevel, kvs...)
}

// Errorf 错误
func (p *logrusLogger) Errorf(msg string, kvs ...interface{}) {
	p.Error("msg", fmt.Sprintf(msg, kvs...))
}

// Critical 严重的
func (p *logrusLogger) Critical(kvs ...interface{}) {
	p.pubLog(CriticalLevel, kvs...)
}

// Criticalf 严重的
func (p *logrusLogger) Criticalf(msg string, kvs ...interface{}) {
	p.Critical("msg", fmt.Sprintf(msg, kvs...))
}

// Panic panic
func (p *logrusLogger) Panic(kvs ...interface{}) {
	p.pubLog(PanicLevel, kvs...)
}

// Panicf panic
func (p *logrusLogger) Panicf(msg string, kvs ...interface{}) {
	p.Panic("msg", fmt.Sprintf(msg, kvs...))
}
