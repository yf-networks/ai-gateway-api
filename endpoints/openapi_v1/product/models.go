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

package product

import (
	"net/http"

	"github.com/yf-networks/ai-gateway-api/lib/xreq"
	"github.com/yf-networks/ai-gateway-api/model/iauth"
	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/stateful/container"
)

var _ xreq.Handler = ProductModelsAction

var ProductModelsRoute = &xreq.Endpoint{
	Path:       "/products/{product_name}/models",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ProductModelsAction),
	Authorizer: iauth.FA(iauth.FeatureProduct, iauth.ActionRead),
}

func ProductModelsAction(req *http.Request) (interface{}, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	list, err := container.ClusterManager.FetchClusterList(req.Context(), &icluster_conf.ClusterFilter{
		Product: product,
	})
	if err != nil {
		return nil, err
	}

	models := make([]string, 0)
	modelMap := make(map[string]string)

	for _, pp := range list {
		if pp.LLMConfig != nil && pp.LLMConfig.Enable != nil && *pp.LLMConfig.Enable {
			for _, mapping := range pp.LLMConfig.ModelMappings {
				modelMap[*mapping.Value] = *mapping.Key
			}

			for _, model := range pp.LLMConfig.Models {
				if _, ok := modelMap[model]; ok {
					continue
				}

				modelMap[model] = ""
			}
		}
	}

	for serverModel, clientModel := range modelMap {
		if clientModel != "" {
			models = append(models, clientModel)
			continue
		}

		models = append(models, serverModel)
	}

	return models, nil
}
