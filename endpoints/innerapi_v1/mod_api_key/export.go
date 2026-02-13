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

package mod_api_key

import (
	"net/http"

	"github.com/yf-networks/ai-gateway-api/endpoints/innerapi_v1/export_util"
	"github.com/yf-networks/ai-gateway-api/lib/xreq"
	"github.com/yf-networks/ai-gateway-api/model/iauth"
	"github.com/yf-networks/ai-gateway-api/stateful/container"
)

// ExportRoute route
var ExportRoute = &xreq.Endpoint{
	Path:       "/configs/mod-api-key",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ExportAction),
	Authorizer: iauth.FA(iauth.FeatureAPIKey, iauth.ActionExport),
}

var _ xreq.Handler = ExportAction

// ExportAction action
func ExportAction(req *http.Request) (interface{}, error) {
	return exportActionProcess(req)
}

func exportActionProcess(req *http.Request) (interface{}, error) {
	param, err := export_util.NewExportFromReq(req)
	if err != nil {
		return nil, err
	}

	return container.APIKeyRuleManager.ConfigExport(req.Context(), param.Version)
}
