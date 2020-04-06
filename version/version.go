package version

import (
	"fmt"

	"github.com/go-trellis/common/builder"
)

// ShowAllInfo 展示详细信息
func ShowAllInfo() {
	builder.Show()
}

// Version 版本信息
func Version() string {
	return fmt.Sprintf("(version: %s, %s)", builder.CompilerVersion, builder.ProgramVersion)
}

// BuildInfo returns goVersion, Author and buildTime information.
func BuildInfo() string {
	return fmt.Sprintf("(go=%s, user=%s, date=%s)", builder.CompilerVersion, builder.Author, builder.BuildTime)
}
