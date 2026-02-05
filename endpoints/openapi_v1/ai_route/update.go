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
	"net/http"

	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/lib/xreq"
	"github.com/yf-networks/ai-gateway-api/model/iai_route"
	"github.com/yf-networks/ai-gateway-api/model/iauth"
	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/stateful/container"
)

// UpdateRulesRequest defines the request parameters for updating AI route rules
type UpdateRulesRequest struct {
	Rules []*iai_route.Rule `json:"rules"`
}

// UpdateRoute is the endpoint definition for updating AI route rules
var UpdateRoute = &xreq.Endpoint{
	Path:       "/products/{product_name}/ai-route-rules",
	Method:     http.MethodPatch,
	Handler:    xreq.Convert(UpdateAction),
	Authorizer: iauth.FAP(iauth.FeatureAIRoute, iauth.ActionCreate),
}

// newReqParam parses and validates the update request parameters
func newReqParam(req *http.Request) (*UpdateRulesRequest, error) {
	param := &UpdateRulesRequest{}
	if err := xreq.BindJSON(req, param); err != nil {
		return nil, err
	}

	// Validate all rules
	for i, rule := range param.Rules {
		if err := iai_route.ValidateRule(rule, i); err != nil {
			return nil, xerror.WrapParamError(err)
		}
	}

	return param, nil
}

// updateProcess handles the business logic for updating AI route rules
func updateProcess(req *http.Request, param *UpdateRulesRequest) (interface{}, error) {
	// Fetch default route rules for validation
	defaultRules, err := container.RouteRuleManager.FetchDefaultRouteRules(req.Context(), nil)
	if err != nil {
		return nil, err
	}

	if len(defaultRules) == 0 {
		return nil, xerror.WrapParamErrorWithMsg("Must set default route rule")
	}

	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	// Fetch the internal AI route product
	products, err := container.ProductManager.FetchProducts(req.Context(), &ibasic.ProductFilter{
		Name: &product.Name,
	})
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, xerror.WrapRecordNotExist(product.Name)
	}

	// Build route rules parameter
	routeRules := ProductRouteRuleParam{
		AdvanceRouteRules: buildAdvanceRouteRules(req.Context(), param.Rules),
	}

	// Get cluster names from rules and fetch cluster information
	clusterNames := getClusterNames(param.Rules)
	clusterMap := make(map[string]int64)
	if len(clusterNames) > 0 {
		clusterList, err := container.ClusterManager.FetchClusterList(req.Context(), &icluster_conf.ClusterFilter{Names: clusterNames})
		if err != nil {
			return nil, err
		}
		clusterMap = buildClusterMap(clusterList)
	}

	// Convert advance rules to internal configuration format
	advanceRules, err := convertAdvanceRules2IrouteConf(routeRules.AdvanceRouteRules, clusterMap)
	if err != nil {
		return nil, err
	}

	// Create or update AI route rules
	err = container.AIRouteRuleManager.CreateAIRouteRule(req.Context(),
		param.Rules,
		products[0],
		advanceRules)
	if err != nil {
		return nil, err
	}

	// Build response
	response := make(map[string]interface{})
	response["rules"] = param.Rules

	return response, err
}

// buildClusterMap creates a mapping from cluster name to cluster ID
func buildClusterMap(clusterList []*icluster_conf.Cluster) map[string]int64 {
	clusterMap := make(map[string]int64)
	for _, cluster := range clusterList {
		clusterMap[cluster.Name] = cluster.ID
	}
	return clusterMap
}

// getClusterNames extracts unique cluster names from rules
func getClusterNames(rules []*iai_route.Rule) []string {
	names := make([]string, 0)
	nameMap := make(map[string]bool)
	for _, rule := range rules {
		name := rule.Basic.ExpectAction.Forward.ClusterName
		if _, ok := nameMap[name]; ok {
			continue
		}
		nameMap[name] = true

		names = append(names, name)
	}
	return names
}

// UpdateAction implements the xreq.Handler interface for updating AI route rules
var _ xreq.Handler = UpdateAction

// UpdateAction is the main handler for updating AI route rules
func UpdateAction(req *http.Request) (interface{}, error) {
	param, err := newReqParam(req)
	if err != nil {
		return nil, err
	}

	return updateProcess(req, param)
}
