package parser

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/seerx/goql/pkg/require"

	"github.com/graphql-go/graphql"

	"github.com/seerx/goql/internal/parser/types"

	"github.com/seerx/goql/internal/reflects"
)

// InjectQuery 注入类型查询函数
type InjectQuery func(injectType reflect.Type) (info *InjectInfo, e error)

// FuncProp 函数属性
type FuncProp struct {
	Desc string // 描述信息

}

type FuncParam struct {
	Arg FuncArg           // 参数
	Typ *reflects.TypeDef // 参数类型
}

// FuncDef 函数信息定义
type FuncDef struct {
	Struct     *StructDef         // 所属结构类型
	Func       reflect.Value      // 函数
	Result     *reflects.TypeDef  // 返回值定义
	RequestArg *reflects.TypeDef  // graphql 提交的参数类型
	Args       []*FuncParam       // 函数输入参数列表
	Info       *reflects.FuncInfo // 函数相关信息
	Prop       *FuncProp
}

func (fd *FuncDef) CreateResolver(ivp *InputVarsPool,
	ovp *OutputVarsPool) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (i interface{}, e error) {
		// 解析请求参数
		inputType := fd.RequestArg
		require := require.New()
		var input reflect.Value
		if inputType != nil {
			// 解析
			rp := &RequestParam{InputVars: ivp}
			var err error
			input, err = rp.Parse(&p, inputType, fd.Struct != nil && fd.Struct.ExplodeParams)
			if err != nil {
				panic(err)
			}
			require = rp.Require
		}

		//rFunc := *fd
		var closers []io.Closer
		defer func() {
			for _, c := range closers {
				c.Close()
			}
		}()
		// 构造参数表
		ctx := &ArgContext{
			Input:          input,
			Param:          &p,
			Require:        require,
			InjectValueMap: map[reflect.Type]reflect.Value{},
		}
		args := make([]reflect.Value, len(fd.Args))
		for n, a := range fd.Args {
			val := a.Arg.CreateValue(ctx)
			if a.Arg.IsInjectInterface() {
				// 注入的接口，判断是否需要 Close
				if closer := getCloser(val); closer != nil {
					closers = append(closers, closer)
				}
			}
			args[n] = val
		}
		// 执行函数
		res := fd.Func.Call(args)
		if res == nil || len(res) != 2 {
			// 没有返回值，或这返回值的数量不是两个
			panic(fmt.Errorf("Resolver <%s> error, need return values", fd.Info.String()))
		}
		out := res[0].Interface()
		errOut := res[1].Interface()
		var err error = nil
		if errOut != nil {
			ok := false
			err, ok = errOut.(error)
			if !ok {
				panic(fmt.Errorf("Resolver <%s> error, second return must be error", fd.Info.String()))
			}
		}

		return out, err
	}
}

// IsDescFunc 是否描述函数
// 描述函数特点，以 Desc 为后缀，且返回一个 string
func IsDescFunc(method reflect.Method) bool {
	mn := method.Name
	if strings.HasSuffix(mn, "Desc") {
		typ := method.Type
		if typ.NumIn() == 1 && typ.NumOut() == 1 {
			out := typ.Out(0)
			return out.Kind() == reflect.String
		}
	}
	return false
}

// GetResolveName 获取 graphql 接口名称
func (fd *FuncDef) GetResolveName() string {
	if fd.Struct == nil {
		return fd.Info.Name
	}
	return fd.Struct.Prefix + fd.Info.Name
}

// ParseDescription 解析函数的描述信息
func (fd *FuncDef) ParseDescription(structInstance interface{}) {
	descFuncName := fd.Info.Name + "Desc"
	typ := reflect.TypeOf(structInstance)
	method, ok := typ.MethodByName(descFuncName)
	if ok {
		if IsDescFunc(method) {
			out := method.Func.Call([]reflect.Value{reflect.ValueOf(structInstance)})
			desc, ok := out[0].Interface().(string)
			if ok {
				fd.Prop.Desc = desc
			}
		}
	}
}

// ParseFunc 解析函数信息
// @param fn 函数
// @param info 函数信息
// @param structInstance 函数所在结构
func ParseFunc(fnObj reflect.Value,
	fnType reflect.Type,
	info *reflects.FuncInfo,
	structDef *StructDef,
	injectQuery InjectQuery) (*FuncDef, error) {
	// 函数开始
	//fnObj := reflect.ValueOf(fn)
	//fnType := reflect.TypeOf(fn)

	outCount := fnType.NumOut()
	if outCount != 2 {
		return nil, errors.New("功能函数必须有两个返回值，且第二个必须是 error 类型")
	}

	def := &FuncDef{
		Struct: structDef,
		Info:   info,
		Func:   fnObj,
		Prop:   &FuncProp{},
	}

	// 分析返回值，返回值必须有两个，第一个是返回的内容，第二个是 error
	// 解析第一个返回参数
	def.Result = reflects.ParseType(fnType.Out(0))
	// 检查返回值的类型
	if err := def.Result.CheckType(); err != nil {
		return nil, err
	}
	// 检查第二个返回值是否是 error
	if !types.IsTypeError(fnType.Out(1)) {
		// 不是 error
		return nil, errors.New("第二个返回值必须是 error")
	}

	// 分析函数参数
	inCount := fnType.NumIn()
	def.Args = make([]*FuncParam, inCount)
	requestArgName := ""
	for n := 0; n < inCount; n++ {
		argType := fnType.In(n)
		typ := reflects.ParseType(argType)
		if n == 0 && structDef != nil {
			// 第一个参数，且该函数在结构体中
			def.Args[n] = &FuncParam{
				Arg: &StructArg{
					Def:   structDef,
					IsPtr: typ.IsPtr,
				},
				Typ: typ,
			}
		} else if types.IsTypeResolveParams(typ.RealType) {
			// ResolveParams 参数
			def.Args[n] = &FuncParam{
				Arg: &ResolveParamsArg{
					IsPtr: typ.IsPtr,
				},
				Typ: typ,
			}
		} else if require.IsRequirement(typ.Type) {
			// Requirement
			def.Args[n] = &FuncParam{
				Arg: &RequireArg{
					IsPtr: typ.IsPtr,
				},
				Typ: typ,
			}
		} else if typ.IsInterface {
			// 接口类型参数，必须是注入类型
			inject, err := injectQuery(typ.RealType)
			if err != nil {
				return nil, fmt.Errorf("注入类型 %s: %s", typ.Name, err.Error())
			}
			if inject == nil {
				return nil, fmt.Errorf("找不到注入类型 %s: ", typ.Name)
			}
			def.Args[n] = &FuncParam{
				Arg: &InjectArg{
					IsInterface: true,
					Inject:      inject,
				},
				Typ: typ,
			}
		} else if typ.IsStruct {
			// 结构类型，必须是指针
			if !typ.IsPtr {
				// 不是指针
				return nil, fmt.Errorf("函数接收的 struct 参数必须是指针类型: %s", typ.Name)
			}
			// 优先认定注入类型，从注入类型中查找
			inject, err := injectQuery(typ.RealType)
			if err != nil {
				return nil, fmt.Errorf("Unsupprt %s: %s", typ.Name, err.Error())
			}
			if inject != nil {
				// 注入类型
				def.Args[n] = &FuncParam{
					Arg: &InjectArg{
						IsInterface: false,
						Inject:      inject,
					},
					Typ: typ,
				}
			} else {
				// 接收参数的类型
				if requestArgName != "" {
					// 已经有一个接收请求参数的结构
					return nil, fmt.Errorf("不可以有多个接收参数的 struct %s，已经存在<%s>", typ.Name, requestArgName)
				} else {
					// 还没有接收请求的结构类型参数
					requestArgName = typ.Name
					// 检查类型是否合法
					if e := typ.CheckType(); e != nil {
						return nil, fmt.Errorf("Unsupport data type: %s", e.Error())
					}

					def.RequestArg = typ
					def.Args[n] = &FuncParam{
						Arg: &RequestArg{
							ArgType: typ,
						},
						Typ: typ,
					}
				}
			}
		} else {
			// 不支持的数据类型
			return nil, fmt.Errorf("Unsupport data type: %s", typ.Name)
		}
	}

	return def, nil
}

// getCloser 判断该值是否时实现了 io.Closer 接口
// 如果实现，则返回 io.Closer 接口
// 否则返回 nil
func getCloser(value reflect.Value) io.Closer {
	val := value.Interface()
	if val != nil {
		closer, ok := val.(io.Closer)
		if ok {
			return closer
		}
	}
	return nil
}
