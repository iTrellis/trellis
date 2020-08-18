/*
Copyright © 2018 Henry Huang <hhh@rutcode.com>

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

package builder

import (
	"fmt"
	"strings"

	"github.com/dimiro1/banner"
	colorable "github.com/mattn/go-colorable"
)

// 编译信息
var (
	ProgramName     string
	ProgramVersion  string
	CompilerVersion string
	BuildTime       string
	Author          string
)

// Show 显示项目信息
func Show() {
	bannerLogo :=
		`
*************************************************************** ***
*************************************************************** ***
***  _______   _______   _______   _       _       _    ______  ***
*** (__   __) |  ___  \ /  _____) / \     / \     / \  /  ____) ***
***    | |    | |___| | | (____   | |     | |     | |  | (____  ***
***    | |    |  ___ /  |  ____)  | |     | |     | |  \____  \ ***
***    | |    | |  \ \  | (_____  | |__/\ | |__/\ | |   ____) | ***
***    \_/    \_/   \_/ \_______) \_____/ \_____/ \_/  (______/ ***
***                                                             ***
*******************************************************************
******************** Compile Environment **************************
*** Program Name     : %s
*** Program Version  : %s
*** Compiler Version : %s
*** Build Time       : %s
*** Author           : %s
*******************************************************************
******************** Running Environment **************************
*** Go running version : {{ .GoVersion }}
*** Go running OS      : {{ .GOOS }} {{ .GOARCH }}
*** Startup time       : {{ .Now "2006-01-02 15:04:05" }}
*******************************************************************
*******************************************************************
`
	newBanner := fmt.Sprintf(bannerLogo, ProgramName, ProgramVersion, CompilerVersion, BuildTime, Author)

	banner.Init(colorable.NewColorableStdout(), true, true, strings.NewReader(newBanner))
}
