package goql

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/seerx/goql/internal/parser"

	"github.com/seerx/goql/pkg/log"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/seerx/goql/internal/reflects"
	"github.com/seerx/goql/pkg/gqlh"
	"github.com/seerx/goql/pkg/param"
)

// GQL gql 信息
type GQL struct {
	log           log.Logger
	queryManager  *gqlh.ResolverManager
	mutateManager *gqlh.ResolverManager
	injectManager *gqlh.InjectManager

	schema       *graphql.Schema
	customConfig *handler.Config // 用户传入的 cfg

	tobeQueries   []interface{} // 查询函数
	tobeMutations []interface{} // 操作函数
	inputVars     *parser.InputVarsPool
}

// NewGQL
func NewGQL(logger log.Logger) *GQL {
	if logger == nil {
		logger = &log.DefaultLogger{}
	}
	return &GQL{
		log:           logger,
		queryManager:  gqlh.NewResolverManager(logger),
		mutateManager: gqlh.NewResolverManager(logger),
		injectManager: gqlh.NewInjectManager(logger),
	}
}

// CreateHandler 创建处理器
func (g *GQL) CreateHandler(cfg *handler.Config) http.Handler {
	schema, err := g.GenerateSchema()
	if err != nil {
		panic(err)
	}

	handlerCfg := &handler.Config{
		Schema:           schema,
		Pretty:           cfg.Pretty,
		Playground:       cfg.Playground,
		GraphiQL:         cfg.GraphiQL,
		ResultCallbackFn: cfg.ResultCallbackFn,
		FormatErrorFn:    cfg.FormatErrorFn,
		RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
			var root map[string]interface{}
			if cfg.RootObjectFn != nil {
				root = cfg.RootObjectFn(ctx, r)
			}
			if root == nil {
				root = map[string]interface{}{}
			}
			param.InjectStoreContext(ctx, r, root)
			return root
		},
	}

	return handler.New(handlerCfg)
}

func (g *GQL) GenerateSchema() (*graphql.Schema, error) {
	g.inputVars = parser.NewInputVarsPool(g.log)
	ovp := parser.NewOutputVarsPool(g.log)

	cfg := graphql.SchemaConfig{}

	// 解析查询函数
	for _, obj := range g.tobeQueries {
		if err := g.parseFuncs(g.queryManager, obj); err != nil {
			g.log.Error("Parser Query: " + err.Error())
		}
	}
	queryFields := graphql.Fields{}
	g.queryManager.GenerateResolvers(g.inputVars, ovp, func(name string, resolver *graphql.Field) {
		queryFields[name] = resolver
	})
	if len(queryFields) > 0 {
		cfg.Query = graphql.NewObject(
			graphql.ObjectConfig{
				Name:   "Query",
				Fields: queryFields,
			})
	}

	// 解析操作函数
	for _, obj := range g.tobeMutations {
		if err := g.parseFuncs(g.mutateManager, obj); err != nil {
			g.log.Error("Parser Mutate: " + err.Error())
		}
	}
	mutateFields := graphql.Fields{}
	g.mutateManager.GenerateResolvers(g.inputVars, ovp, func(name string, resolver *graphql.Field) {
		mutateFields[name] = resolver
	})
	if len(mutateFields) > 0 {
		cfg.Mutation = graphql.NewObject(
			graphql.ObjectConfig{
				Name:   "Mutation",
				Fields: mutateFields,
			})
	}

	schema, err := graphql.NewSchema(cfg)
	return &schema, err
}

// RegisterInject 注册注入函数
func (g *GQL) RegisterInject(injectFn interface{}) {
	g.injectManager.AddInject(injectFn)
}

// RegisterQuery 注册查询
func (g *GQL) RegisterQuery(funcLoader interface{}) {
	g.tobeQueries = append(g.tobeQueries, funcLoader)
	//if err := g.parseFuncs(g.queryManager, funcLoader); err != nil {
	//	g.log.Error("RegisterQuery: " + err.Error())
	//}
}

// RegisterMutate 注册操作
func (g *GQL) RegisterMutate(funcLoader interface{}) {
	g.tobeMutations = append(g.tobeMutations, funcLoader)
	//if err := g.parseFuncs(g.mutateManager, funcLoader); err != nil {
	//	g.log.Error("RegisterMutate: " + err.Error())
	//}
}

func (g *GQL) parseFuncs(manager *gqlh.ResolverManager, resolveObject interface{}) error {
	injectQuery := func(injectType reflect.Type) (info *parser.InjectInfo, e error) {
		// 查询是否是注入类型，如果是，直接返回
		inject := g.injectManager.FindInject(injectType)
		return inject, nil
	}

	typ := reflects.ParseType(reflect.TypeOf(resolveObject))
	if typ.IsStruct {
		// 结构体
		manager.ParseStruct(resolveObject, injectQuery)
		return nil
	} else if typ.IsFunc {
		// 独立函数
		manager.ParserFunction(resolveObject, injectQuery)
		return nil
	}
	return fmt.Errorf("Unsupport type: %s", typ.Name)
}
