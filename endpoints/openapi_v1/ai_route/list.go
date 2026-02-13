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

package ai_route

import (
	"context"
	"fmt"
	"net/http"

	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/lib/xreq"
	"github.com/yf-networks/ai-gateway-api/model/iai_route"
	"github.com/yf-networks/ai-gateway-api/model/iauth"
	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/iroute_conf"
	"github.com/yf-networks/ai-gateway-api/stateful/container"
)

// AdvanceRouteRule defines the structure for advanced AI routing rules
type AdvanceRouteRule struct {
	Name          string                      `json:"name"`
	Description   string                      `json:"description"`
	Expression    string                      `json:"expression" validate:"required,min=1"`
	ClusterName   string                      `json:"cluster_name"`
	RouteAction   *iroute_conf.RouteAction    `json:"action,omitempty"`
	ExtendActions []*iroute_conf.ExtendAction `json:"extend_actions"`
}

// DefaultRouteRule defines the structure for default AI routing rules
type DefaultRouteRule struct {
	Cmd           string                      `json:"cmd"`
	Params        []string                    `json:"params"`
	Description   string                      `json:"description"`
	RouteAction   *iroute_conf.RouteAction    `json:"action,omitempty"`
	ExtendActions []*iroute_conf.ExtendAction `json:"extend_actions"`
}

// ProductRouteRuleParam defines the request parameters for product route rules
type ProductRouteRuleParam struct {
	AdvanceRouteRules []*AdvanceRouteRule `json:"forward_rules" validate:"dive"`
}

// ListRoute is the endpoint definition for listing AI route rules
var ListRoute = &xreq.Endpoint{
	Path:       "/products/{product_name}/ai-route-rules",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(ListAction),
	Authorizer: iauth.FAP(iauth.FeatureAIRoute, iauth.ActionRead),
}

// listActionProcess processes the request to fetch AI route rules
func listActionProcess(req *http.Request) (interface{}, error) {
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	// Fetch AI route rules from the manager
	return container.AIRouteRuleManager.FetchAIRouteRules(req.Context(), &iai_route.AIRouteFilter{
		ProductName: &product.Name,
	})
}

// ListAction implements the xreq.Handler interface for listing AI route rules
var _ xreq.Handler = ListAction

// ListAction is the main handler for listing AI route rules
func ListAction(req *http.Request) (interface{}, error) {
	return listActionProcess(req)
}

// buildAdvanceRouteRules converts internal rule structures to API response format
func buildAdvanceRouteRules(ctx context.Context, rules []*iai_route.Rule) []*AdvanceRouteRule {
	var advanceRouteRules []*AdvanceRouteRule
	for _, rule := range rules {
		advanceRouteRules = append(advanceRouteRules, &AdvanceRouteRule{
			Name:        rule.Name,
			Expression:  iai_route.BuildAIRouteCond(ctx, rule.Basic),
			ClusterName: rule.Basic.ExpectAction.Forward.ClusterName,
		})
	}
	return advanceRouteRules
}

// convertAdvanceRules2IrouteConf converts advanced route rules to internal configuration format
func convertAdvanceRules2IrouteConf(rules []*AdvanceRouteRule, clusterMap map[string]int64) ([]*iroute_conf.AdvanceRouteRule, error) {
	var advanceRouteRules []*iroute_conf.AdvanceRouteRule

	for _, rule := range rules {
		// Validate if cluster exists in the map
		if _, ok := clusterMap[rule.ClusterName]; !ok {
			return nil, xerror.WrapParamErrorWithMsg(fmt.Sprintf("not found cluster:%s", rule.ClusterName))
		}

		advanceRouteRules = append(advanceRouteRules, &iroute_conf.AdvanceRouteRule{
			Name:        rule.Name,
			Expression:  rule.Expression,
			ClusterName: rule.ClusterName,
			ClusterID:   clusterMap[rule.ClusterName],
		})
	}

	return advanceRouteRules, nil
}
