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

// Level log level
type Level int32

// define levels
const (
	TraceLevel = Level(iota)
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel
	PanicLevel

	LevelNameUnknown  = "NULL"
	LevelNameTrace    = "TRAC"
	LevelNameDebug    = "DEBU"
	LevelNameInfo     = "INFO"
	LevelNameWarn     = "WARN"
	LevelNameError    = "ERRO"
	LevelNameCritical = "CRIT"
	LevelNamePanic    = "PANC"

	levelColorDebug    = "\033[32m%s\033[0m" // grenn
	levelColorInfo     = "\033[37m%s\033[0m" // white
	levelColorWarn     = "\033[34m%s\033[0m" // blue
	levelColorError    = "\033[33m%s\033[0m" // yellow
	levelColorCritical = "\033[31m%s\033[0m" // red
	levelColorPanic    = "\033[35m%s\033[0m" // perple

)

// LevelColors printer's color
var LevelColors = map[Level]string{
	TraceLevel:    levelColorInfo,
	DebugLevel:    levelColorDebug,
	InfoLevel:     levelColorInfo,
	WarnLevel:     levelColorWarn,
	ErrorLevel:    levelColorError,
	CriticalLevel: levelColorCritical,
}

// ToLevelName 等级转换为名称
func ToLevelName(lvl Level) string {
	switch lvl {
	case TraceLevel:
		return LevelNameTrace
	case DebugLevel:
		return LevelNameDebug
	case InfoLevel:
		return LevelNameInfo
	case WarnLevel:
		return LevelNameWarn
	case ErrorLevel:
		return LevelNameError
	case CriticalLevel:
		return LevelNameCritical
	case PanicLevel:
		return LevelNamePanic
	default:
		return LevelNameUnknown
	}
}
