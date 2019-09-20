package gqlh

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/seerx/goql/internal/parser"

	"github.com/seerx/goql/pkg/log"

	"github.com/seerx/goql/internal/parser/types"

	"github.com/seerx/goql/internal/reflects"
)

// InjectManager 注入管理
type InjectManager struct {
	log       log.Logger
	injectMap map[reflect.Type]*parser.InjectInfo
}

var errOfInject = errors.New("Param injectFn must be a func  like -- func(ctx context.Context, r *http.Request, gp *graphql.ResolveParams) *returnType, and returnType must be a struct or interface")

// NewInjectManager 创建
func NewInjectManager(log log.Logger) *InjectManager {
	return &InjectManager{
		log:       log,
		injectMap: map[reflect.Type]*parser.InjectInfo{},
	}
}

// FindInject 根据类型查找注入信息
func (im *InjectManager) FindInject(typ reflect.Type) *parser.InjectInfo {
	info, ok := im.injectMap[typ]
	if ok {
		return info
	}
	return nil
}

// AddInject 添加注入函数
func (im *InjectManager) AddInject(injectFn interface{}) {
	typ := reflect.TypeOf(injectFn)
	if typ.Kind() != reflect.Func {
		panic(errOfInject)
	}

	oc := typ.NumOut()
	if oc != 1 {
		// 返回参数必须只能有一个
		panic(errOfInject)
	}
	out := typ.Out(0)
	outp := reflects.ParseType(out)
	if !outp.IsInterface {
		// 返回值不是 interface
		if outp.IsPrimitive { // 返回值是原始类型
			panic(errOfInject)
		}
		if !outp.IsStruct { // 返回值不是结构类型
			panic(errOfInject)
		}
		if !outp.IsPtr { // 返回值不是指针类型
			panic(errOfInject)
		}
	}

	old, ok := im.injectMap[outp.RealType]
	if ok { // 已经存在
		panic(fmt.Errorf("Type [%s] is Registered by func [%s]", outp.Name, old.Info.String()))
	}

	// 判断输入参数
	ic := typ.NumIn()
	if ic != 3 {
		panic(errOfInject)
	}
	// 第一个参数
	p1 := reflects.ParseType(typ.In(0))
	if !types.IsContext(p1.RealType) {
		panic(errOfInject)
	}
	// 第二个参数
	p2 := reflects.ParseType(typ.In(1))
	if !p2.IsPtr || !types.IsHttpRequest(p2.RealType) {
		panic(errOfInject)
	}
	// 第三个参数
	p3 := reflects.ParseType(typ.In(2))
	if !p3.IsPtr || !types.IsResolveParams(p3.RealType) {
		panic(errOfInject)
	}

	// 注册
	fnInfo := reflects.ParseFuncInfo(injectFn)
	inject := &parser.InjectInfo{
		Type: outp.RealType,
		Info: fnInfo,
		Func: reflect.ValueOf(injectFn),
	}
	im.injectMap[inject.Type] = inject
}
