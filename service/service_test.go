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

package service

import (
	"testing"

	"github.com/iTrellis/common/testutils"
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
