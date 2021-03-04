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

package version

import (
	"fmt"

	"github.com/iTrellis/common/builder"
)

// ShowAllInfo 展示详细信息
func ShowAllInfo() {
	builder.Show()
}

// Version 版本信息
func Version() string {
	return fmt.Sprintf("%s, version: %s (branch: %s, revision: %s)",
		builder.ProgramName, builder.ProgramVersion,
		builder.ProgramBranch, builder.ProgramRevision,
	)
}

// BuildInfo returns goVersion, Author and buildTime information.
func BuildInfo() string {
	return fmt.Sprintf("(go=%s, user=%s, date=%s)",
		builder.CompilerVersion, builder.Author, builder.BuildTime)
}
