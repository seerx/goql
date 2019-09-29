package types

import (
	"reflect"
)

// IsInt 是否整数类型
func IsInt(kind reflect.Kind) bool {
	return kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64
}

// IsFloat 是否浮点类型
func IsFloat(kind reflect.Kind) bool {
	return kind == reflect.Float32 ||
		kind == reflect.Float64
}

// IsString 是否字符串类型
func IsString(kind reflect.Kind) bool {
	return kind == reflect.String
}

// IsBoolType 是否布尔类型
func IsBool(kind reflect.Kind) bool {
	return kind == reflect.Bool
}

//// IsTimeType 是否时间类型
//func IsTimeType(typ reflect.Type) bool {
//	return typ == reflect.TypeOf(time.Time{})
//}

//// IsStructType 是否结构类型
//func IsStructType(kind reflect.Kind) bool {
//	return kind == reflect.Struct
//}
