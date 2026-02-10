package product_cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/lib/xreq"
	"github.com/yf-networks/ai-gateway-api/model/iauth"
)

var _ xreq.Handler = ListModelProvidersAction

var ListModelProvidersRoute = &xreq.Endpoint{
	Path:       "/products/{product_name}/model-providers",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ListModelProvidersAction),
	Authorizer: iauth.FA(iauth.FeatureProductCluster, iauth.ActionCreate),
}

func ListModelProvidersAction(req *http.Request) (interface{}, error) {
	return listModelProvidersProcess(req.Context())
}

type ModelProvider struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func listModelProvidersProcess(ctx context.Context) (interface{}, error) {
	data, err := os.ReadFile("conf/ai/models.json")
	if err != nil {
		return nil, xerror.WrapParamError(fmt.Errorf("failed to read config conf/ai/models.json: %w", err))
	}

	var models []ModelProvider
	err = json.Unmarshal(data, &models)
	return models, err
}
