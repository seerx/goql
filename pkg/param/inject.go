package param

import (
	"context"
	"net/http"

	"github.com/seerx/goql/internal/inject"
)

// InjectStoreContext 存储Conext 和 Request，以备注入时使用
func InjectStoreContext(ctx context.Context, r *http.Request, root map[string]interface{}) {
	root[inject.KeyOfConext] = ctx
	root[inject.KeyOfRequest] = r
}
