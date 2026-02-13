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

package route_conf

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yf-networks/ai-gateway-api/lib"
	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/model/iroute_conf"
	"github.com/yf-networks/ai-gateway-api/model/iversion_control"
	"github.com/yf-networks/ai-gateway-api/storage/rdb/internal/dao"
)

var _ iroute_conf.RouteRuleStorager = &RouteRuleStorager{}

func NewRouteRuleStorager(dbCtxFactory lib.DBContextFactory,
	versionControlStorager iversion_control.VersionControlStorager) *RouteRuleStorager {
	return &RouteRuleStorager{
		dbCtxFactory:           dbCtxFactory,
		versionControlStorager: versionControlStorager,
	}
}

type RouteRuleStorager struct {
	dbCtxFactory lib.DBContextFactory

	versionControlStorager iversion_control.VersionControlStorager
}

func (rs *RouteRuleStorager) UpsertAdvanceProductRule(ctx context.Context, product *ibasic.Product,
	rules []*iroute_conf.AdvanceRouteRule) error {

	now := lib.PTime(time.Now())

	daoAdvanceRules := []*dao.TRouteAdvanceRuleParam{}
	for _, one := range rules {
		daoAdvanceRule := &dao.TRouteAdvanceRuleParam{
			ProductID:   &product.ID,
			Name:        lib.PString(one.Name),
			ClusterID:   lib.PInt64(one.ClusterID),
			Expression:  lib.PString(one.Expression),
			Description: lib.PString(one.Description),
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		daoAdvanceRules = append(daoAdvanceRules, daoAdvanceRule)
	}

	dbCtx, err := rs.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	if _, err := dao.TRouteAdvanceRuleList(dbCtx, &dao.TRouteAdvanceRuleParam{
		ProductID: &product.ID,
		LockMode:  &dao.ModeForUpdate,
	}); err != nil {
		return err
	}

	if _, err := dao.TRouteAdvanceRuleDelete(dbCtx, &dao.TRouteAdvanceRuleParam{
		ProductID: &product.ID,
	}); err != nil {
		return err
	}

	if len(daoAdvanceRules) > 0 {
		if _, err := dao.TRouteAdvanceRuleCreate(dbCtx, daoAdvanceRules...); err != nil {
			return err
		}
	}

	return nil
}

func (rs *RouteRuleStorager) UpsertDefaultProductRule(ctx context.Context,
	product *ibasic.Product,
	rule *iroute_conf.DefaultRouteRule) error {
	dbCtx, err := rs.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	if _, err := dao.TRouteDefaultRuleDelete(dbCtx, &dao.TRouteDefaultRuleParam{
		ProductID: &product.ID,
	}); err != nil {
		return err
	}

	if rule != nil {
		paramBytes, err := json.Marshal(rule.Params)
		if err != nil {
			return err
		}

		defaultRule := &dao.TRouteDefaultRuleParam{
			ProductID:   &product.ID,
			Cmd:         &rule.Cmd,
			Params:      lib.PString(string(paramBytes)),
			Description: &rule.Description,
		}

		if rule.RouteAction != nil {
			b, _ := json.Marshal(rule.RouteAction)
			defaultRule.RouteAction = lib.PString(string(b))
		} else {
			defaultRule.RouteAction = lib.PString("")
		}

		if _, err := dao.TRouteDefaultRuleCreate(dbCtx, defaultRule); err != nil {
			return err
		}
	}

	return nil
}

func (rs *RouteRuleStorager) FetchProductRule(ctx context.Context, product *ibasic.Product,
	clusterList []*icluster_conf.Cluster) (*iroute_conf.ProductRouteRule, error) {
	m, err := rs.FetchRouteRules(ctx, []*ibasic.Product{product}, clusterList)
	if err != nil {
		return nil, err
	}

	return m[product.ID], nil
}

func (rs *RouteRuleStorager) FetchRouteRules(ctx context.Context, products []*ibasic.Product,
	clusterList []*icluster_conf.Cluster) (map[int64]*iroute_conf.ProductRouteRule, error) {

	// 1. Prepare product mappings
	productIDs := prepareProductMappings(products)

	// 2. Initialize database context
	dbCtx, err := rs.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Initialize result map
	product2ProductRouteRule := initializeResultMap(products)

	// 4. Fetch and process advance rules
	clusterMap := icluster_conf.ClusterList2MapByID(clusterList)
	err = fetchAndProcessAdvanceRules(dbCtx, productIDs, clusterMap, product2ProductRouteRule)
	if err != nil {
		return nil, err
	}

	// 5. Fetch and process default rules
	err = fetchAndProcessDefaultRules(dbCtx, product2ProductRouteRule)
	if err != nil {
		return nil, err
	}

	return product2ProductRouteRule, nil
}

// prepareProductMappings prepares product ID list and ID to name mapping
func prepareProductMappings(products []*ibasic.Product) []int64 {
	productIDs := make([]int64, 0, len(products))
	for _, product := range products {
		productIDs = append(productIDs, product.ID)
	}
	return productIDs
}

// initializeResultMap initializes the result map with empty ProductRouteRule for each product
func initializeResultMap(products []*ibasic.Product) map[int64]*iroute_conf.ProductRouteRule {
	result := make(map[int64]*iroute_conf.ProductRouteRule, len(products))

	for _, product := range products {
		result[product.ID] = &iroute_conf.ProductRouteRule{
			AdvanceRouteRules: make([]*iroute_conf.AdvanceRouteRule, 0),
			DefaultRouteRule:  nil,
		}
	}

	return result
}

// fetchAndProcessAdvanceRules fetches advance rules from database and processes them
func fetchAndProcessAdvanceRules(dbCtx *lib.DBContext, productIDs []int64,
	clusterMap map[int64]*icluster_conf.Cluster,
	product2ProductRouteRule map[int64]*iroute_conf.ProductRouteRule) error {
	filter := &dao.TRouteAdvanceRuleParam{}
	if len(productIDs) > 0 {
		filter.ProductIDs = productIDs
	}

	advanceRules, err := dao.TRouteAdvanceRuleList(dbCtx, filter)
	if err != nil {
		return err
	}

	for _, rule := range advanceRules {
		productRule, exists := product2ProductRouteRule[rule.ProductID]
		if !exists {
			// Skip rules for products not in the input list
			continue
		}

		cluster, clusterExists := clusterMap[rule.ClusterID]
		if !clusterExists {
			// Skip rules with non-existent clusters
			continue
		}

		advanceRouteRule := &iroute_conf.AdvanceRouteRule{
			Name:        rule.Name,
			Description: rule.Description,
			Expression:  rule.Expression,
			ClusterName: cluster.Name,
			ClusterID:   rule.ClusterID,
		}

		productRule.AdvanceRouteRules = append(productRule.AdvanceRouteRules, advanceRouteRule)
	}

	return nil
}

// fetchAndProcessDefaultRules fetches default rules from database and processes them
func fetchAndProcessDefaultRules(dbCtx *lib.DBContext,
	product2ProductRouteRule map[int64]*iroute_conf.ProductRouteRule) error {

	defaultRules, err := dao.TRouteDefaultRuleList(dbCtx, &dao.TRouteDefaultRuleParam{})
	if err != nil {
		return err
	}

	for _, defaultRule := range defaultRules {
		productRule, exists := product2ProductRouteRule[defaultRule.ProductID]
		if !exists {
			// Skip default rules for products not in the input list
			continue
		}

		routeAction := &iroute_conf.RouteAction{}

		// Safely unmarshal JSON
		if defaultRule.RouteAction != "" {
			if err := json.Unmarshal([]byte(defaultRule.RouteAction), routeAction); err != nil {
				return fmt.Errorf("unmarshal:%s is error:%s", defaultRule.RouteAction, err.Error())
			}
		}

		// Create default route rule, append to advance rules.
		advanceRouteRule := &iroute_conf.AdvanceRouteRule{
			Expression:  defaultRule.Cmd,
			Description: defaultRule.Description,
			ClusterName: routeAction.Forward.ClusterName,
		}

		productRule.AdvanceRouteRules = append(productRule.AdvanceRouteRules, advanceRouteRule)
	}

	return nil
}

func (rs *RouteRuleStorager) FetchDefaultRouteRules(ctx context.Context, products []*ibasic.Product) ([]*iroute_conf.DefaultRouteRule, error) {
	var productIDs []int64
	for _, one := range products {
		productIDs = append(productIDs, one.ID)
	}

	dbCtx, err := rs.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	defaultRules, err := dao.TRouteDefaultRuleList(dbCtx, &dao.TRouteDefaultRuleParam{
		ProductIDs: productIDs,
	})
	if err != nil {
		return nil, err
	}

	rules := make([]*iroute_conf.DefaultRouteRule, 0)
	for _, one := range defaultRules {
		defaultRule, err := newRouteDefaultRule(one)
		if err != nil {
			return nil, err
		}
		rules = append(rules, defaultRule)
	}

	return rules, nil
}

func (rs *RouteRuleStorager) FetchAdvanceRouteRules(ctx context.Context, products []*ibasic.Product, clusters []*icluster_conf.Cluster) ([]*iroute_conf.AdvanceRouteRule, error) {
	var productIDs []int64
	for _, one := range products {
		productIDs = append(productIDs, one.ID)
	}

	clusterMap := icluster_conf.ClusterList2MapByID(clusters)

	dbCtx, err := rs.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	advanceRules, err := dao.TRouteAdvanceRuleList(dbCtx, &dao.TRouteAdvanceRuleParam{
		ProductIDs: productIDs,
	})
	if err != nil {
		return nil, err
	}

	rules := make([]*iroute_conf.AdvanceRouteRule, 0)
	for _, one := range advanceRules {
		advanceRule := &iroute_conf.AdvanceRouteRule{
			Name:        one.Name,
			Description: one.Description,
			Expression:  one.Expression,
			ClusterID:   one.ClusterID,
		}

		if one.ClusterID > 0 {
			cluster := clusterMap[one.ClusterID]
			if cluster == nil {
				continue
			}

			advanceRule.ClusterName = cluster.Name
		}

		advanceRule.ClusterName = one.Name
		rules = append(rules, advanceRule)
	}

	return rules, nil
}

func newRouteDefaultRule(param *dao.TRouteDefaultRule) (*iroute_conf.DefaultRouteRule, error) {
	rule := &iroute_conf.DefaultRouteRule{
		Cmd:         param.Cmd,
		Description: param.Description,
	}

	if err := json.Unmarshal([]byte(param.Params), &rule.Params); err != nil {
		return nil, xerror.WrapDirtyDataErrorWithMsg("Params: %s, err: %v", string(param.Params), err)
	}

	if param.RouteAction != "" {
		routeAction := &iroute_conf.RouteAction{}
		json.Unmarshal([]byte(param.RouteAction), &routeAction)
		rule.RouteAction = routeAction
	} else {
		routeAction := &iroute_conf.RouteAction{}

		switch rule.Cmd {
		case iroute_conf.DefaultRouteRuleCmdForward:
			routeAction.Forward = &iroute_conf.ActionForward{
				ClusterName: rule.Params[0],
			}
		case iroute_conf.DefaultRouteRuleCmdSendResponse:
			routeAction.Response = &iroute_conf.ActionResponse{
				StatusCode: rule.Params[0],
			}
		}
		rule.RouteAction = routeAction
	}

	return rule, nil
}
