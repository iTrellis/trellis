/*
Copyright © 2019 Henry Huang <hhh@rutcode.com>

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

package xorm_ext

// Committer Defination
type Committer interface {
	TX(fn interface{}, repos ...interface{}) error
	TXWithName(fn interface{}, name string, repos ...interface{}) error
	NonTX(fn interface{}, repos ...interface{}) error
	NonTXWithName(fn interface{}, name string, repos ...interface{}) error
}

// TXFunc Transation function
type TXFunc func(repos []interface{}) error

// Inheritor inherit function
type Inheritor interface {
	Inherit(repo interface{}) error
}

// Inherit inherit a new repository from origin repository
func Inherit(new, origin interface{}) error {
	if inheritor, ok := new.(Inheritor); ok {
		return inheritor.Inherit(origin)
	}
	return nil
}

// Deriver derive function
type Deriver interface {
	Derive() (repo interface{}, err error)
}

// Derive derive from developer function
func Derive(origin interface{}) (interface{}, error) {
	if deriver, ok := origin.(Deriver); ok {
		return deriver.Derive()
	}
	return nil, nil
}
