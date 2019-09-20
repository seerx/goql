package types

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/graphql-go/graphql"
)

var (
	errorType         = reflect.TypeOf((*error)(nil)).Elem()
	resolveParamsType = reflect.TypeOf(graphql.ResolveParams{})

	typeOfHTTPRequest = reflect.TypeOf(http.Request{})
	typeOfContext     = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeOfGqlParams   = reflect.TypeOf(graphql.ResolveParams{})
	typeOfTime        = reflect.TypeOf(time.Time{})
)

// IsTime 判断是否 time.Time 类型
func IsTime(typ reflect.Type) bool {
	return typ == typeOfTime
}

// IsHttpRequest 判断是否 http.Request 类型
func IsHttpRequest(typ reflect.Type) bool {
	return typ == typeOfHTTPRequest
}

// IsContext 判断是否 context.Context 类型
func IsContext(typ reflect.Type) bool {
	return typ == typeOfContext
}

// IsResolveParams 判断是否 graphql.ResolveParams 类型
func IsResolveParams(typ reflect.Type) bool {
	return typ == typeOfGqlParams
}

// IsTypeError 判断是否是 error 接口
func IsTypeError(typ reflect.Type) bool {
	//return typ.Implements(errorType)
	return errorType == typ
}

// IsTypeResolveParams 判断是否 graphql.ResolveParams 类型
func IsTypeResolveParams(typ reflect.Type) bool {
	return resolveParamsType == typ
}
