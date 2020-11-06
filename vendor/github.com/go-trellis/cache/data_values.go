/*
Copyright © 2016 Henry Huang <hhh@rutcode.com>

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

package cache

import "time"

// ValueMode define value mode
type ValueMode int

// ValueMode
const (
	// only one value
	ValueModeUnique ValueMode = iota
	// The table is a bag table, which can have many objects
	// but only one instance of each object, per key.
	ValueModeBag
	// The table is a duplicate_bag table, which can have many objects,
	// including multiple copies of the same object, per key.
	ValueModeDuplicateBag
)

// DataValues define k-vs struct
type DataValues struct {
	Key    interface{}
	Values []interface{}
	Exists map[interface{}]bool
	Expire *time.Time
}
