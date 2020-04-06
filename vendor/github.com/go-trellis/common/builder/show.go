// GNU GPL v3 License
// Copyright (c) 2018 github.com:go-trellis

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
