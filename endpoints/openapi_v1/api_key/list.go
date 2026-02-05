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
	"net/http"

	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/stateful/container"

	"github.com/yf-networks/ai-gateway-api/lib/xreq"
	"github.com/yf-networks/ai-gateway-api/model/iauth"
)

var ListRoute = &xreq.Endpoint{
	Path:       "/products/{product_name}/api-keys",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ListAction),
	Authorizer: iauth.FAP(iauth.FeatureAPIKey, iauth.ActionReadAll),
}

var _ xreq.Handler = OneAction

func ListAction(req *http.Request) (interface{}, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	list, err := container.APIKeyManager.FetchAPIKeyList(req.Context(), &icluster_conf.APIKeyFilter{
		ProductName: &product.Name,
	})
	if err != nil {
		return nil, err
	}

	return newResponse(list)
}

func newResponse(list []*icluster_conf.APIKeyParam) ([]*icluster_conf.APIKeyParam, error) {
	for i, one := range list {
		if one.Key != nil {
			remainingQuota, err := icluster_conf.GetRemainingQuota(one)
			if err != nil {
				return nil, err
			}
			list[i].RemainingQuota = remainingQuota
		}
	}

	return list, nil
}
