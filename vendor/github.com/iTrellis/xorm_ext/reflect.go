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

import (
	"fmt"
	"reflect"

	"github.com/iTrellis/common/errors"
)

// MapErrorTypes 可以支持的返回的错误类型
var mapErrorTypes = map[reflect.Type]bool{
	// 普通错误类型
	reflect.TypeOf((*error)(nil)).Elem(): true,
	// common error错误类型
	reflect.TypeOf((*errors.ErrorCode)(nil)).Elem():   true,
	reflect.TypeOf((*errors.SimpleError)(nil)).Elem(): true,
}

// AddErrorTypes 增加支持的错误类型
func AddErrorTypes(errType reflect.Type) {
	mapErrorTypes[errType] = true
}

// Function Flags
const (
	Logic = iota
	BeforeLogic
	AfterLogic
	OnError
	AfterCommit
)

// LogicFuncs logic functions
type LogicFuncs struct {
	BeforeLogic interface{}
	AfterLogic  interface{}
	OnError     interface{}
	Logic       interface{}
	AfterCommit interface{}
}

// DeepFields relect interface deep fields
func DeepFields(iface interface{}, vType reflect.Type, fields []reflect.Value) interface{} {

	ift := reflect.TypeOf(iface)
	if ift == vType {
		return iface
	}
	ifv := reflect.ValueOf(iface)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
		ift = ifv.Type()
	}

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		switch v.Kind() {
		case reflect.Struct:
			var deepIns interface{}
			if v.CanAddr() {
				deepIns = DeepFields(v.Addr().Interface(), vType, fields)
			} else {
				deepIns = DeepFields(v.Interface(), vType, fields)
			}

			if deepIns != nil {
				return deepIns
			}
		}
	}
	return nil
}

// GetLogicFuncs reflect logic function
func GetLogicFuncs(fn interface{}) (funcs LogicFuncs) {
	switch fn := fn.(type) {
	case TXFunc, func(repos []interface{}) (err error):
		{
			funcs.Logic = fn
		}
	case map[int]interface{}:
		{
			if hookBeforefn, exist := fn[BeforeLogic]; exist {
				funcs.BeforeLogic = hookBeforefn
			}

			if logicfn, exist := fn[Logic]; exist {
				funcs.Logic = logicfn
			}

			if hookAfterfn, exist := fn[AfterLogic]; exist {
				funcs.AfterLogic = hookAfterfn
			}

			if errfn, exist := fn[OnError]; exist {
				funcs.OnError = errfn
			}

			if afterCommitfn, exist := fn[AfterCommit]; exist {
				funcs.AfterCommit = afterCommitfn
			}
		}
	default:
		funcs.Logic = fn
	}

	return
}

// CallFunc execute transaction function with logic functions and args
func CallFunc(fn interface{}, funcs LogicFuncs, args []interface{}) ([]interface{}, error) {
	if fn == nil {
		return nil, nil
	}

	switch _logicFunc := fn.(type) {
	case TXFunc:
		{
			return nil, _logicFunc(args)
		}
	case func(repos []interface{}) (err error):
		{
			if err := _logicFunc(args); err != nil {
				return nil, err
			}
			return nil, nil
		}
	default:
		values, err := call(fn, args...)
		if err != nil {
			if funcs.OnError != nil {
				_, _ = call(funcs.OnError, err)
			}
			return nil, err
		}
		return values, nil
	}
}

func call(fn interface{}, args ...interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(fn)
	if !v.IsValid() {
		return nil, fmt.Errorf("call of nil")
	}
	typ := v.Type()
	if typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("non-function of type %s", typ)
	}
	if !goodFunc(typ) {
		return nil, fmt.Errorf("the last return value should be an error type")
	}
	numIn := typ.NumIn()
	var dddType reflect.Type
	if typ.IsVariadic() {
		if len(args) < numIn-1 {
			return nil, fmt.Errorf("wrong number of args: got %d want at least %d, type: %v", len(args), numIn-1, typ)
		}
		dddType = typ.In(numIn - 1).Elem()
	} else {
		if len(args) != numIn {
			return nil, fmt.Errorf("wrong number of args: got %d want %d, type: %v", len(args), numIn, typ)
		}
	}
	argv := make([]reflect.Value, len(args))
	for i, arg := range args {
		value := reflect.ValueOf(arg)
		// Compute the expected type. Clumsy because of variadics.
		var argType reflect.Type
		if !typ.IsVariadic() || i < numIn-1 {
			argType = typ.In(i)
		} else {
			argType = dddType
		}

		var err error
		if argv[i], err = prepareArg(value, argType); err != nil {
			return nil, fmt.Errorf("arg %d: %s", i, err)
		}
	}

	result := v.Call(argv)
	resultLen := len(result)

	var resultValues []interface{}

	for _, v := range result {
		resultValues = append(resultValues, v.Interface())
	}

	if resultLen == 1 {
		if resultValues[0] != nil {
			return nil, resultValues[0].(error)
		}
	} else if resultLen > 1 {
		if resultValues[resultLen-1] != nil {
			return resultValues[0 : resultLen-1], resultValues[resultLen-1].(error)
		}
		return resultValues[0 : resultLen-1], nil
	}

	return nil, nil
}

func goodFunc(typ reflect.Type) bool {
	if typ.NumOut() == 0 ||
		(typ.NumOut() > 0 && mapErrorTypes[typ.Out(typ.NumOut()-1)]) {
		return true
	}

	return false
}

func prepareArg(value reflect.Value, argType reflect.Type) (reflect.Value, error) {
	if !value.IsValid() {
		if !canBeNil(argType) {
			return reflect.Value{}, fmt.Errorf("value is nil; should be of type %s", argType)
		}
		value = reflect.Zero(argType)
	}
	if !value.Type().AssignableTo(argType) {
		return reflect.Value{}, fmt.Errorf("value has type %s; should be %s", value.Type(), argType)
	}
	return value, nil
}

func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	}
	return false
}
