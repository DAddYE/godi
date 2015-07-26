package godi

import (
	"fmt"
	"reflect"
)

// --------
//
// typeInfo struct bundles up a string type name with the reflect type
// into a handy tuple.
//
// --------
type typeInfo struct {
	typeName    string
	reflectType *reflect.Type
}

func newtypeInfo(typeName string, reflectType *reflect.Type) *typeInfo {
	ti := &typeInfo{typeName: typeName, reflectType: reflectType}

	if reflectType != nil {
		str := typeToString(*reflectType)
		ti.typeName = str
	}
	return ti
}

func (p *typeInfo) Type() reflect.Type {
	if p.reflectType == nil {
		t := typeMap[p.typeName]
		if t == nil {
			panic(fmt.Sprintf("Can't find type '%s', did you forget to register it?", p.typeName))
		}
		p.reflectType = t
	}
	return *p.reflectType
}

//
// ----------- typeInfo
//
