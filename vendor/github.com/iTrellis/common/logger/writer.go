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
	"strings"

	"github.com/iTrellis/common/event"
)

// Writer 写对象
type Writer interface {
	event.Subscriber
}

func generateLogs(evt *Event, separator string) string {

	logs := make([]string, 0, 2+len(evt.Prefixes)+len(evt.Fields))
	logs = append(logs, evt.Time.Format("2006/01/02T15:04:05.000"), ToLevelName(evt.Level))

	for _, pref := range evt.Prefixes {
		logs = append(logs, replacerString(pref.(string), separator))
	}

	for _, v := range evt.Fields {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Ptr, reflect.Struct, reflect.Map:
			bs, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			logs = append(logs, string(bs))
		default:
			logs = append(logs, fmt.Sprintf("%+v", v))
		}
	}
	return strings.Join(logs, separator) + "\n"
}

func replacerString(origin, replacer string) string {
	str := ReplaceString(origin, " ")
	return ReplaceString(str, replacer)
}

// ReplaceString 替换字符串
func ReplaceString(origin, replacer string) string {
	return strings.Replace(origin, replacer, "", -1)
}
