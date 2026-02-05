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

	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/lib/xreq"
	"github.com/yf-networks/ai-gateway-api/model/iauth"
)

var DeleteRoute = &xreq.Endpoint{
	Path:       "/products/{product_name}/api-keys/{api_key_name}",
	Method:     http.MethodDelete,
	Handler:    xreq.Convert(DeleteAction),
	Authorizer: iauth.FAP(iauth.FeatureAPIKey, iauth.ActionDelete),
}

var _ xreq.Handler = DeleteAction

func DeleteAction(req *http.Request) (interface{}, error) {
	oneReq, err := newReq4One(req)
	if err != nil {
		return nil, err
	}

	products, err := container.ProductManager.FetchProducts(req.Context(), &ibasic.ProductFilter{
		Name: oneReq.ProductName,
	})
	if err != nil {
		return nil, err
	}
	if len(products) != 1 {
		return nil, xerror.WrapParamErrorWithMsg("Invalid Product")
	}

	err = container.APIKeyManager.DeleteAPIKey(req.Context(), &icluster_conf.APIKeyFilter{
		Name:        oneReq.APIKeyName,
		ProductName: oneReq.ProductName,
	})
	return nil, err
}
