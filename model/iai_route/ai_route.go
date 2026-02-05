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

package iai_route

import (
	"context"
	"fmt"

	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/iroute_conf"
	"github.com/yf-networks/ai-gateway-api/model/itxn"
	"github.com/yf-networks/ai-gateway-api/model/iversion_control"
	"github.com/yf-networks/ai-gateway-api/stateful"
)

type AIRouteRuleManager struct {
	storager AIRouteRuleStorager
	txn      itxn.TxnStorager

	versionControlManager *iversion_control.VersionControlManager
	routeStorager         iroute_conf.RouteRuleStorager
}

type AIRouteFilter struct {
	ProductName *string
}

type AIRouteRuleStorager interface {
	FetchAIRouteRules(ctx context.Context, filter *AIRouteFilter) ([]*Rule, error)
	CreateAIRouteRules(ctx context.Context, param []*Rule) error
}

func NewAIRouteRuleManager(txn itxn.TxnStorager, storager AIRouteRuleStorager,
	versionControlManager *iversion_control.VersionControlManager,
	routeStorager iroute_conf.RouteRuleStorager) *AIRouteRuleManager {
	return &AIRouteRuleManager{
		txn:                   txn,
		storager:              storager,
		versionControlManager: versionControlManager,
		routeStorager:         routeStorager,
	}
}

func (rlm *AIRouteRuleManager) CreateAIRouteRule(ctx context.Context,
	rules []*Rule, product *ibasic.Product,
	advanceRules []*iroute_conf.AdvanceRouteRule) (err error) {
	err = rlm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		err = rlm.storager.CreateAIRouteRules(ctx, rules)
		if err != nil {
			return err
		}

		return rlm.routeStorager.UpsertAdvanceProductRule(ctx, product, advanceRules)
	})

	return
}

func (rlm *AIRouteRuleManager) FetchAIRouteRules(ctx context.Context, filter *AIRouteFilter) (rules []*Rule, err error) {
	err = rlm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		rules, err = rlm.storager.FetchAIRouteRules(ctx, filter)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func BuildAIRouteCond(ctx context.Context, basicInfo *BasicInfo) string {
	cond := buildDomainCondition(basicInfo.Domain)
	cond = buildPathCondition(cond, basicInfo.PathFilter)
	cond = buildMethodCondition(cond, basicInfo.Method)
	cond = buildHeaderConditions(cond, basicInfo.HeaderFilters)
	cond = buildModelCondition(cond, basicInfo.ModelFilter)

	stateful.AccessLogger.Debug(fmt.Sprintf("BuildAIRouteCond:%s", cond))
	return cond
}

// buildDomainCondition constructs domain condition for AI routing
func buildDomainCondition(domain *string) string {
	if domain != nil && *domain != "" {
		return fmt.Sprintf("req_host_in(\"%s\")", *domain)
	}
	return ""
}

// buildPathCondition constructs path condition for AI routing
func buildPathCondition(currentCond string, pathFilter *PathFilter) string {
	if pathFilter == nil || pathFilter.MatchMode == nil {
		return currentCond
	}

	var pathCond string
	switch *pathFilter.MatchMode {
	case MatchModePrefix:
		pathCond = fmt.Sprintf("req_path_prefix_in(\"%s\", %t)",
			*pathFilter.Path,
			*pathFilter.IgnoreCase)
	case MatchModeExact:
		pathCond = fmt.Sprintf("req_path_in(\"%s\", %t)",
			*pathFilter.Path,
			*pathFilter.IgnoreCase)
	case MatchModeSuffix:
		pathCond = fmt.Sprintf("req_path_suffix_in(\"%s\", %t)",
			*pathFilter.Path,
			*pathFilter.IgnoreCase)
	}

	return combineConditions(currentCond, pathCond)
}

// buildMethodCondition constructs HTTP method condition for AI routing
func buildMethodCondition(currentCond string, method *string) string {
	if method == nil || *method == "" {
		return currentCond
	}

	methodCond := fmt.Sprintf("req_method_in(\"%s\")", *method)
	return combineConditions(currentCond, methodCond)
}

// buildHeaderConditions constructs header conditions for AI routing
func buildHeaderConditions(currentCond string, headerFilters []*BasicHeaderFilter) string {
	for _, header := range headerFilters {
		if header.MatchMode == nil {
			continue
		}

		var headerCond string
		switch *header.MatchMode {
		case MatchModePrefix:
			headerCond = fmt.Sprintf("req_header_value_prefix_in(\"%s\", \"%s\", %t)",
				*header.Key,
				*header.Value,
				*header.IgnoreCase)
		case MatchModeExact:
			headerCond = fmt.Sprintf("req_header_value_in(\"%s\", \"%s\", %t)",
				*header.Key,
				*header.Value,
				*header.IgnoreCase)
		case MatchModeSuffix:
			headerCond = fmt.Sprintf("req_header_value_suffix_in(\"%s\", \"%s\", %t)",
				*header.Key,
				*header.Value,
				*header.IgnoreCase)
		}

		if headerCond != "" {
			currentCond = combineConditions(currentCond, headerCond)
		}
	}

	return currentCond
}

// buildModelCondition constructs model condition for AI routing
func buildModelCondition(currentCond string, modelFilter *ModelFilter) string {
	if modelFilter == nil || modelFilter.Name == nil || *modelFilter.Name == "" {
		return currentCond
	}

	modelCond := fmt.Sprintf("req_body_json_in(\"%s\", \"%s\", %t)",
		*modelFilter.Name,
		*modelFilter.Pattern,
		*modelFilter.IgnoreCase)
	return combineConditions(currentCond, modelCond)
}

// combineConditions combines multiple conditions with logical AND operator
func combineConditions(currentCond, newCond string) string {
	if newCond == "" {
		return currentCond
	}
	if currentCond == "" {
		return newCond
	}
	return fmt.Sprintf("%s&&%s", currentCond, newCond)
}
