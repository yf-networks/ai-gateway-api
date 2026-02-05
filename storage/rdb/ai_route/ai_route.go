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
	"encoding/json"

	"github.com/yf-networks/ai-gateway-api/lib"
	"github.com/yf-networks/ai-gateway-api/model/iai_route"
	"github.com/yf-networks/ai-gateway-api/storage/rdb/internal/dao"
)

type RDBAIRouteRuleStorager struct {
	dbCtxFactory lib.DBContextFactory
}

var _ iai_route.AIRouteRuleStorager = &RDBAIRouteRuleStorager{}

func NewRDBAIRouteRuleStorager(dbCtxFactory lib.DBContextFactory) *RDBAIRouteRuleStorager {
	return &RDBAIRouteRuleStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

func (rpps *RDBAIRouteRuleStorager) FetchAIRouteRules(ctx context.Context, filter *iai_route.AIRouteFilter) ([]*iai_route.Rule, error) {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	param := &dao.TAIRouteRuleParam{}
	if filter != nil {
		param.ProductName = filter.ProductName
	}

	list, err := dao.TAIRouteRuleList(dbCtx, param)
	if err != nil {
		return nil, err
	}

	rules := make([]*iai_route.Rule, 0)

	for _, one := range list {
		tmp := ConvertToRule(one)
		rules = append(rules, tmp)
	}

	return rules, nil
}

// ConvertToRule 将TAIRouteRule转换为Rule（不进行错误校验）
func ConvertToRule(dbRule *dao.TAIRouteRule) *iai_route.Rule {
	if dbRule == nil {
		return nil
	}

	rule := &iai_route.Rule{
		Name: dbRule.Name,
	}

	// 转换Basic
	if dbRule.Basic != "" {
		var basic iai_route.BasicInfo
		json.Unmarshal([]byte(dbRule.Basic), &basic)
		rule.Basic = &basic
	}

	return rule
}

// ConvertToTAIRouteRule 将Rule转换为TAIRouteRule（不进行错误校验）
func ConvertToTAIRouteRule(rule *iai_route.Rule) *dao.TAIRouteRuleParam {
	if rule == nil {
		return nil
	}

	dbRule := &dao.TAIRouteRuleParam{
		Name:  &rule.Name,
		Basic: lib.PString(""),
	}

	// 转换Basic
	if rule.Basic != nil {
		basicJSON, _ := json.Marshal(rule.Basic)
		dbRule.Basic = lib.PString(string(basicJSON))
	}

	dbRule.UpdatedAt = lib.PTimeNow()

	return dbRule
}

func (rpps *RDBAIRouteRuleStorager) CreateAIRouteRules(ctx context.Context, rules []*iai_route.Rule) error {
	dbRules := make([]*dao.TAIRouteRuleParam, 0)

	for i, _ := range rules {
		dbRule := ConvertToTAIRouteRule(rules[i])
		dbRule.IDX = lib.PInt64(int64(i + 1))
		dbRule.CreatedAt = dbRule.UpdatedAt
		dbRules = append(dbRules, dbRule)
	}

	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TAIRouteRuleDelete(dbCtx, &dao.TAIRouteRuleParam{})
	if err != nil {
		return err
	}

	if len(rules) > 0 {
		_, err = dao.TAIRouteRuleCreate(dbCtx, dbRules...)
		if err != nil {
			return err
		}
	}

	return nil
}
