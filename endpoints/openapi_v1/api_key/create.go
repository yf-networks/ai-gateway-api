// Copyright(c) 2026 Beijing Yingfei Networks Technology Co.Ltd.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http: //www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package api_key

import (
	"context"
	"net/http"

	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/lib/xreq"
	"github.com/yf-networks/ai-gateway-api/model/iauth"
	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/stateful/container"
)

var _ xreq.Handler = APIKeyCreateAction

var APIKeyCreateRoute = &xreq.Endpoint{
	Path:       "/products/{product_name}/api-keys",
	Method:     http.MethodPost,
	Handler:    xreq.Convert(APIKeyCreateAction),
	Authorizer: iauth.FA(iauth.FeatureAPIKey, iauth.ActionCreate),
}

func APIKeyCreateAction(req *http.Request) (interface{}, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	param := &icluster_conf.APIKeyParam{}
	if err := xreq.BindJSON(req, param); err != nil {
		return nil, err
	}

	return APIKeyCreateProcess(req.Context(), param, product)
}

func APIKeyCreateProcess(ctx context.Context, param *icluster_conf.APIKeyParam, product *ibasic.Product) (*ibasic.Product, error) {
	if err := checkCreateAPIKey(param, product.Name); err != nil {
		return nil, xerror.WrapParamError(err)
	}

	err := container.APIKeyManager.CreateAPIKey(ctx, &icluster_conf.APIKeyParam{
		Name:          param.Name,
		Enable:        param.Enable,
		Key:           param.Key,
		IsLimit:       param.IsLimit,
		Limit:         param.Limit,
		ExpiredTime:   param.ExpiredTime,
		AllowedModels: param.AllowedModels,
		AllowedCIDR:   param.AllowedCIDR,
		ProductName:   &product.Name,
	})

	return nil, err
}
