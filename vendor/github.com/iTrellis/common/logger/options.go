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

import "io"

// STDOption option
type STDOption func(*STDOptions)

// STDLevel set std logger level
func STDLevel(lvl Level) STDOption {
	return func(c *STDOptions) {
		c.level = lvl
	}
}

// STDWriter io writer
func STDWriter(w io.Writer) STDOption {
	return func(c *STDOptions) {
		c.writer = w
	}
}

// FileOption 操作配置函数
type FileOption func(*FileOptions)

// FileLevel 设置等级
func FileLevel(lvl Level) FileOption {
	return func(f *FileOptions) {
		f.level = lvl
	}
}

// FileBuffer 设置Chan的大小
func FileBuffer(buffer int) FileOption {
	return func(f *FileOptions) {
		f.chanBuffer = buffer
	}
}

// FileSeparator 设置打印分隔符
func FileSeparator(separator string) FileOption {
	return func(f *FileOptions) {
		f.separator = separator
	}
}

// FileFileName 设置文件名
func FileFileName(name string) FileOption {
	return func(f *FileOptions) {
		f.fileName = name
	}
}

// FileMaxLength 设置最大文件大小
func FileMaxLength(length int64) FileOption {
	return func(f *FileOptions) {
		f.maxLength = length
	}
}

// FileMaxBackupFile 文件最大数量
func FileMaxBackupFile(num int) FileOption {
	return func(f *FileOptions) {
		f.maxBackupFile = num
	}
}

// FileMoveFileType 设置移动文件的类型
func FileMoveFileType(typ MoveFileType) FileOption {
	return func(f *FileOptions) {
		f.moveFileType = typ
	}
}

// FilePublisher set publisher
func FilePublisher(pub Publisher) FileOption {
	return func(c *FileOptions) {
		c.publisher = pub
	}
}

// LogrusOption 操作配置函数
type LogrusOption func(*LogrusOptions)

// LogrusLevel 设置等级
func LogrusLevel(lvl Level) LogrusOption {
	return func(f *LogrusOptions) {
		f.level = lvl
	}
}
