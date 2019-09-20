package goql

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"

	"github.com/graphql-go/handler"

	"github.com/seerx/goql/pkg/log"
)

var instance *GQL

// InjectTemplate 注入函数模板
type InjectTemplate func(ctx context.Context, r *http.Request, gp *graphql.ResolveParams) interface{}

// Configure 配置 GQL 基础信息
func Configure(logger log.Logger) {
	if instance == nil {
		instance = NewGQL(logger)
	} else {
		instance.log = logger
	}
}

func Get() *GQL {
	if instance == nil {
		Configure(nil)
	}
	return instance
}

func RegisterQuery(funcLoader interface{}) {
	Get().RegisterQuery(funcLoader)
}

func RegisterMutate(funcLoader interface{}) {
	Get().RegisterMutate(funcLoader)
}

func RegisterInject(injectFn interface{}) {
	Get().RegisterInject(injectFn)
}

func CreateHandler(cfg *handler.Config) http.Handler {
	return Get().CreateHandler(cfg)
}

func GenerateSchma() (*graphql.Schema, error) {
	return Get().GenerateSchema()
}

// SilentLogger 静默日志
type SilentLogger struct {
}

func (SilentLogger) Info(msg string) {
}

func (SilentLogger) Debug(msg string) {
}

func (SilentLogger) Error(msg string) {
}

func (SilentLogger) Warn(msg string) {
}
