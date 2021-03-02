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
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	kitlog "github.com/go-kit/kit/log"
	"github.com/iTrellis/common/event"
)

// Logger 日志对象
type Logger interface {
	SetLevel(lvl Level)

	WithPrefix(kvs ...interface{}) Logger

	Log(keyvals ...interface{}) error
	Debug(kvs ...interface{})
	Debugf(msg string, kvs ...interface{})
	Info(kvs ...interface{})
	Infof(msg string, kvs ...interface{})
	Warn(kvs ...interface{})
	Warnf(msg string, kvs ...interface{})
	Error(kvs ...interface{})
	Errorf(msg string, kvs ...interface{})
	Critical(kvs ...interface{})
	Criticalf(msg string, kvs ...interface{})
	Panic(kvs ...interface{})
	Panicf(msg string, kvs ...interface{})

	event.Subscriber
}

func genLogs(evt *Event) []interface{} {

	lenFields := len(evt.Fields)
	n := 4 + (lenFields+1)/2*2

	logs := make([]interface{}, 0, n)

	logs = append(logs, "ts", evt.Time.Format("2006/01/02T15:04:05.000"), "level", ToLevelName(evt.Level))

	for i := 0; i < lenFields; i += 2 {
		k := evt.Fields[i]
		var v interface{} = kitlog.ErrMissingValue
		if i+1 < lenFields {
			v = evt.Fields[i+1]
		}
		logs = append(logs, toString(k), toString(v))
	}

	return logs
}

func toString(v interface{}) string {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Ptr, reflect.Struct, reflect.Map:
		bs, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return string(bs)
	case reflect.String:
		return v.(string)
	default:
		return fmt.Sprint(v)
	}
}

// Debug 调试
func Debug(l Logger, fields ...interface{}) {
	l.Debug(fields...)
}

// Debugf 调试
func Debugf(l Logger, msg string, fields ...interface{}) {
	l.Debugf(msg, fields...)
}

// Info 信息
func Info(l Logger, fields ...interface{}) {
	l.Info(fields...)
}

// Infof 信息
func Infof(l Logger, msg string, fields ...interface{}) {
	l.Infof(msg, fields...)
}

// Error 错误
func Error(l Logger, fields ...interface{}) {
	l.Error(fields...)
}

// Errorf 错误
func Errorf(l Logger, msg string, fields ...interface{}) {
	l.Errorf(msg, fields...)
}

// Warn 警告
func Warn(l Logger, fields ...interface{}) {
	l.Warn(fields...)
}

// Warnf 警告
func Warnf(l Logger, msg string, fields ...interface{}) {
	l.Warnf(msg, fields...)
}

// Critical 异常
func Critical(l Logger, fields ...interface{}) {
	l.Critical(fields...)
}

// Criticalf 异常
func Criticalf(l Logger, msg string, fields ...interface{}) {
	l.Criticalf(msg, fields...)
}

// Panic 异常
func Panic(l Logger, fields ...interface{}) {
	l.Panic(fields...)
}

// Panicf 异常
func Panicf(l Logger, msg string, fields ...interface{}) {
	l.Panicf(msg, fields...)
}

// RuntimeCaller stores a stacktrace under the key "stacktrace".
func RuntimeCaller(skip int) func() interface{} {
	return func() interface{} {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			file = "<???>"
			line = 1
		} else {
			slash := strings.LastIndex(file, "/")
			file = file[slash+1:]
		}
		return fmt.Sprintf("%s:%d", file, line)
	}
}

func RuntimeCallers(skip int) func() interface{} {

	return func() interface{} {
		var name, file string
		var line int
		var pc [16]uintptr

		n := runtime.Callers(skip, pc[:])
		for _, pc := range pc[:n] {
			fn := runtime.FuncForPC(pc)
			if fn == nil {
				continue
			}
			file, line = fn.FileLine(pc)

			slash := strings.LastIndex(file, "/")
			file = file[slash+1:]

			name = fn.Name()

			if !strings.HasPrefix(name, "runtime.") {
				slash := strings.LastIndex(name, "/")
				name = name[slash+1:]
				break
			}
		}

		var str string
		switch {
		case name != "":
			str = fmt.Sprintf("%v:%v", name, line)
		case file != "":
			str = fmt.Sprintf("%v:%v", file, line)
		default:
			str = fmt.Sprintf("pc:%x", pc)
		}
		return str
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

func doCaller(hasCaller bool, prefixes []interface{}, keyvals ...interface{}) []interface{} {
	kvs := append(prefixes, keyvals...)
	if !hasCaller {
		return kvs
	}
	bindCallers(kvs)
	return kvs
}
