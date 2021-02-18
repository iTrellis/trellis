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

const (
	// 默认的通道大小
	defaultChanBuffer int = 10000
)

// Debug 调试
func Debug(l Logger, msg string, fields ...interface{}) {
	l.Debug(msg, fields...)
}

// Debugf 调试
func Debugf(l Logger, msg string, fields ...interface{}) {
	l.Debugf(msg, fields...)
}

// Info 信息
func Info(l Logger, msg string, fields ...interface{}) {
	l.Info(msg, fields...)
}

// Infof 信息
func Infof(l Logger, msg string, fields ...interface{}) {
	l.Infof(msg, fields...)
}

// Error 错误
func Error(l Logger, msg string, fields ...interface{}) {
	l.Error(msg, fields...)
}

// Errorf 错误
func Errorf(l Logger, msg string, fields ...interface{}) {
	l.Errorf(msg, fields...)
}

// Warn 警告
func Warn(l Logger, msg string, fields ...interface{}) {
	l.Warn(msg, fields...)
}

// Warnf 警告
func Warnf(l Logger, msg string, fields ...interface{}) {
	l.Warnf(msg, fields...)
}

// Critical 异常
func Critical(l Logger, msg string, fields ...interface{}) {
	l.Critical(msg, fields...)
}

// Criticalf 异常
func Criticalf(l Logger, msg string, fields ...interface{}) {
	l.Criticalf(msg, fields...)
}

// With 增加默认的消息
func With(l Logger, params ...interface{}) Logger {
	return l.With(params...)
}

// WithPrefix 在最前面增加消息
func WithPrefix(l Logger, prefixes ...interface{}) Logger {
	return l.WithPrefix(prefixes...)
}
