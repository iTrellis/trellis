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

package internal

import "strings"

func StringInSlice(needle string, haystack []string) bool {
	for _, hay := range haystack {
		if hay == needle {
			return true
		}
	}
	return false
}

func SuffixInSlice(needle string, haystack []string) bool {
	for _, h := range haystack {
		if strings.HasSuffix(needle, h) {
			return true
		}
	}
	return false
}

func StringContainedInSlice(needle string, haystack []string) bool {
	for _, h := range haystack {
		if strings.Contains(needle, h) {
			return true
		}
	}

	return false
}
