package registry

import (
	"testing"

	"github.com/go-trellis/common/testutils"
)

func TestParseServicePath(t *testing.T) {

	act, err := ParseServicePath("trellis://registry/server/")

	testutils.NotOk(t, err)
	testutils.Equals(t, err.Error(), "prefix is not trellis://registry/service/")
	testutils.Equals(t, act, nil)

	act, err = ParseServicePath("trellis://registry/service/")
	testutils.NotOk(t, err)
	testutils.Equals(t, err.Error(), "path is incorrect")
	testutils.Equals(t, act, nil)

	act, err = ParseServicePath("trellis://registry/service/A")
	testutils.NotOk(t, err)
	testutils.Equals(t, act.Name, "A")
	testutils.Equals(t, act.Version, "")

	act, err = ParseServicePath("trellis://registry/service/A/1")
	testutils.Ok(t, err)
	testutils.Equals(t, act.Name, "A")
	testutils.Equals(t, act.Version, "1")
}
