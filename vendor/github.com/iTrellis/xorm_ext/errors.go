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

import "errors"

// define connector errors
var (
	ErrNotFoundDefaultDatabase    = errors.New("not found default database")
	ErrAtLeastOneRepo             = errors.New("input one repo at least")
	ErrNotFoundTransationFunction = errors.New("not found transation function")
	ErrStructCombineWithRepo      = errors.New("your repository struct should combine repo")
	ErrFailToCreateRepo           = errors.New("fail to create an new repo")
	ErrFailToConvetTXToNonTX      = errors.New("could not convert TX to NON-TX")
	ErrTransactionIsAlreadyBegin  = errors.New("transaction is already begin")
	ErrNonTransactionCantCommit   = errors.New("non-transaction can't commit")
	ErrTransactionSessionIsNil    = errors.New("transaction session is nil")
	ErrNotFoundXormEngine         = errors.New("not found xorm engine")
)
