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
	"fmt"
	"net/http"
	"time"

	"github.com/yf-networks/ai-gateway-api/lib"
	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/stateful/container"

	"github.com/gofrs/uuid"
	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/lib/xreq"
	"github.com/yf-networks/ai-gateway-api/model/iauth"
)

var GenerateTokenRoute = &xreq.Endpoint{
	Path:       "/products/{product_name}/api-keys/actions/generate",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(GenerateTokenAction),
	Authorizer: iauth.FAP(iauth.FeatureAPIKey, iauth.ActionRead),
}

var _ xreq.Handler = GenerateTokenAction

func GenerateTokenAction(req *http.Request) (interface{}, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewV7()
	if err != nil {
		return nil, xerror.WrapParamErrorWithMsg(fmt.Sprintf("generate key is error:%s", err.Error()))
	}

	key := fmt.Sprintf("%s-%s-%d", product.Name, uid.String(), time.Now().Nanosecond())
	id, err := container.APIKeyManager.CreateAPIKeyToken(req.Context(), &icluster_conf.APIKeyTokenParam{
		Key: lib.PString(key),
	})
	if err != nil {
		return nil, err
	}

	response := make(map[string]string)
	response["key"] = fmt.Sprintf("%s-%d", key, id)
	return response, nil
}
