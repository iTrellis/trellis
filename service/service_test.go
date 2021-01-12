package service

import (
	"testing"

	"github.com/go-trellis/common/testutils"
)

func TestParseService(t *testing.T) {
	path := "/service1/v1"

	s, err := ParseService(path)
	testutils.Ok(t, err)
	testutils.Equals(t, &Service{Domain: defDomain, Name: "service1", Version: "v1"}, s)

	path2 := "/test/service2/v2"
	s, err = ParseService(path2)
	testutils.Ok(t, err)
	testutils.Equals(t, &Service{Domain: "test", Name: "service2", Version: "v2"}, s)

	path3 := "/test"
	_, err = ParseService(path3)
	testutils.NotOk(t, err)
}
