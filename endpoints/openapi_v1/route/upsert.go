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
//limitations under the License. All rights reserved.

// Copyright (c) 2021 The BFE Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package route

import (
	"net/http"

	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/lib/xreq"

	"github.com/yf-networks/ai-gateway-api/model/iauth"
	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/iroute_conf"
	"github.com/yf-networks/ai-gateway-api/stateful/container"
)

type ProductRouteRuleParam struct {
	DefaultRouteRule *DefaultRouteRule `json:"default_forward_rule"`
}

type DefaultRouteRule struct {
	Cmd         string                   `json:"cmd"`
	Params      []string                 `json:"params"`
	Description string                   `json:"description"`
	RouteAction *iroute_conf.RouteAction `json:"action,omitempty"`
}

type ProductRouteRuleData struct {
	DefaultRouteRule *DefaultRouteRule `json:"default_route_rule"`
}

func newProductRouteRuleData(pfr *ProductRouteRuleParam) *ProductRouteRuleData {
	return &ProductRouteRuleData{
		DefaultRouteRule: pfr.DefaultRouteRule,
	}
}

func routeRuleParam2routeRule(p *ProductRouteRuleParam) *iroute_conf.ProductRouteRule {
	return &iroute_conf.ProductRouteRule{
		DefaultRouteRule: &iroute_conf.DefaultRouteRule{
			Cmd:         p.DefaultRouteRule.Cmd,
			Params:      p.DefaultRouteRule.Params,
			Description: p.DefaultRouteRule.Description,
			RouteAction: p.DefaultRouteRule.RouteAction,
		},
	}
}

// UpsertRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var UpsertEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/routes",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(UpsertAction),
	Authorizer: iauth.FAP(iauth.FeatureRoute, iauth.ActionUpdate),
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newRuleInfoFromReq(req *http.Request) (*ProductRouteRuleParam, error) {
	rule := &ProductRouteRuleParam{}
	err := xreq.BindJSON(req, rule)
	if err != nil {
		return nil, err
	}

	if rule.DefaultRouteRule == nil {
		return nil, xerror.WrapParamErrorWithMsg("default_forward_rule cant be nil")
	}

	rule.DefaultRouteRule.Cmd = iroute_conf.DefaultExpression

	return rule, err
}

func UpsertActionProcess(req *http.Request, rule *ProductRouteRuleParam) (*ProductRouteRuleData, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	ipfr := routeRuleParam2routeRule(rule)
	err = container.RouteRuleManager.UpsertDefaultProductRule(req.Context(), product, ipfr.DefaultRouteRule)
	if err != nil {
		return nil, err
	}

	return newProductRouteRuleData(rule), nil
}

var _ xreq.Handler = UpsertAction

// UpsertAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func UpsertAction(req *http.Request) (interface{}, error) {
	rule, err := newRuleInfoFromReq(req)
	if err != nil {
		return nil, err
	}

	return UpsertActionProcess(req, rule)
}
