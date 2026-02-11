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

package iroute_conf

import (
	"context"
	"fmt"
	"strings"

	"github.com/bfenetworks/bfe/bfe_basic/condition"
	"github.com/bfenetworks/bfe/bfe_config/bfe_route_conf/route_rule_conf"

	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/model/itxn"
	"github.com/yf-networks/ai-gateway-api/model/iversion_control"
)

type BasicRouteRule struct {
	HostNames   []string
	Paths       []string
	ClusterName string
	ClusterID   int64
	Description string
}

type AdvanceRouteRule struct {
	Name          string
	Description   string
	Expression    string
	ClusterName   string
	ClusterID     int64
	RouteAction   *RouteAction    `json:"action"`
	ExtendActions []*ExtendAction `json:"extend_actions"`
}

type RouteRuleCase struct {
	Description   string
	URL           string
	Method        string
	Body          string
	Header        map[string]string
	ExpectCluster string
	ExpectAction  *RouteAction
}

type ProductRouteRule struct {
	DefaultRouteRule  *DefaultRouteRule
	AdvanceRouteRules []*AdvanceRouteRule
}

type HostUsedInfo struct {
	Type   string
	Detail string
}

type ExtendAction struct {
	Cmd    string   `json:"cmd"`
	Params []string `json:"params"`
}

type RouteAction struct {
	Forward           *ActionForward           `json:"forward,omitempty"`
	GoToAdvancedRules *ActionGoToAdvancedRules `json:"go_to_advanced_rules,omitempty"`
	Redirect          *ActionRedirect          `json:"redirect,omitempty"`
	Response          *ActionResponse          `json:"response,omitempty"`
}

type ActionResponse struct {
	StatusCode  string `json:"status_code"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

type ActionRedirect struct {
	URL string `json:"url"`
}

type ActionGoToAdvancedRules struct {
}

type ActionForward struct {
	ClusterName string `json:"cluster_name"`
	URL         string `json:"url"`
}

type DefaultRouteRule struct {
	Cmd           string
	Params        []string
	Description   string
	RouteAction   *RouteAction    `json:"action"`
	ExtendActions []*ExtendAction `json:"extend_actions"`
}

const (
	HOSTSET = "HOST_SET"
	PATHSET = "PATH_SET"

	DefaultRouteRuleCmdForward      = "FORWARD"
	DefaultRouteRuleCmdRedirect     = "REDIRECT"
	DefaultRouteRuleCmdReturn       = "RETURN"
	DefaultRouteRuleCmdSendResponse = "SEND_RESPONSE"
)

var (
	DefaultExpression = "default_t()"

	zeroVersion = "0000-00-00"
	testProduct = "test"
)

func (prr *ProductRouteRule) HostBeUsed(host string) *HostUsedInfo {
	keyword := fmt.Sprintf(`req_host_in("%s")`, host)

	for _, arr := range prr.AdvanceRouteRules {
		if strings.Contains(arr.Expression, keyword) {
			return &HostUsedInfo{
				Type:   "AdvanceConditionExpression",
				Detail: arr.Expression,
			}
		}
	}

	return nil
}

type ProductRouteRuleConvertResult struct {
	AdvancedRouteRuleFiles []route_rule_conf.AdvancedRouteRuleFile
	DefaultRouteRule       DefaultRouteRuleFile
	ReferClusterNames      []string
}

type DefaultRouteRuleFile struct {
	Cmd           *string
	Params        []string
	RouteAction   *RouteAction
	ExtendActions []*ExtendAction
}

type RouteRuleStorager interface {
	UpsertDefaultProductRule(ctx context.Context, product *ibasic.Product, rule *DefaultRouteRule) error
	UpsertAdvanceProductRule(ctx context.Context, product *ibasic.Product, rules []*AdvanceRouteRule) error
	FetchDefaultRouteRules(ctx context.Context, products []*ibasic.Product) ([]*DefaultRouteRule, error)
	FetchAdvanceRouteRules(ctx context.Context, products []*ibasic.Product, clusters []*icluster_conf.Cluster) ([]*AdvanceRouteRule, error)
	FetchProductRule(ctx context.Context, product *ibasic.Product,
		clusterList []*icluster_conf.Cluster) (*ProductRouteRule, error)
	FetchRouteRules(ctx context.Context, products []*ibasic.Product,
		clusterList []*icluster_conf.Cluster) (map[int64]*ProductRouteRule, error)
}

func NewRouteRuleManager(txn itxn.TxnStorager, storager RouteRuleStorager, clusterStorager icluster_conf.ClusterStorager,
	productStorager ibasic.ProductStorager, versionControlManager *iversion_control.VersionControlManager,
	domainStorager DomainStorager) *RouteRuleManager {
	return &RouteRuleManager{
		txn:                   txn,
		storager:              storager,
		clusterStorager:       clusterStorager,
		productStorager:       productStorager,
		versionControlManager: versionControlManager,
		domainStorager:        domainStorager,
	}
}

type RouteRuleManager struct {
	versionControlManager *iversion_control.VersionControlManager
	txn                   itxn.TxnStorager
	storager              RouteRuleStorager
	clusterStorager       icluster_conf.ClusterStorager
	productStorager       ibasic.ProductStorager
	domainStorager        DomainStorager
}

func (rm *RouteRuleManager) ExpressionVerify(ctx context.Context, expression string) (err error) {
	_, err = condition.Build(expression)
	return err
}

func (rm *RouteRuleManager) FetchDefaultRouteRules(ctx context.Context, products []*ibasic.Product) (rules []*DefaultRouteRule, err error) {
	err = rm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		rules, err = rm.storager.FetchDefaultRouteRules(ctx, products)
		return err
	})

	return
}

func (rm *RouteRuleManager) UpsertDefaultProductRule(ctx context.Context, product *ibasic.Product, rule *DefaultRouteRule) (err error) {
	err = rm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		return rm.storager.UpsertDefaultProductRule(ctx, product, rule)
	})

	return
}

func (rm *RouteRuleManager) UpsertAdvanceProductRule(ctx context.Context, product *ibasic.Product, rules []*AdvanceRouteRule) (err error) {
	err = rm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		return rm.storager.UpsertAdvanceProductRule(ctx, product, rules)
	})

	return
}

func (rm *RouteRuleManager) ClusterDeleteChecker(ctx context.Context, product *ibasic.Product, cluster *icluster_conf.Cluster) error {
	rules, err := rm.storager.FetchAdvanceRouteRules(ctx, []*ibasic.Product{product}, []*icluster_conf.Cluster{cluster})
	if err != nil {
		return err
	}

	if len(rules) > 0 {
		return xerror.WrapModelErrorWithMsg("Rule %s Refer To This Cluster", rules[0].Name)
	}

	defaultRules, err := rm.storager.FetchDefaultRouteRules(ctx, []*ibasic.Product{product})
	if err != nil {
		return err
	}

	if len(defaultRules) == 0 {
		return nil
	}

	return xerror.WrapModelErrorWithMsg("Default rule Refer To This Cluster")
}

func (rm *RouteRuleManager) FetchProductRule(ctx context.Context, product *ibasic.Product) (prr *ProductRouteRule, err error) {
	err = rm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		clusters, err := rm.clusterStorager.FetchClusterList(ctx, &icluster_conf.ClusterFilter{
			Product: product,
		})
		if err != nil {
			return err
		}

		clusters = icluster_conf.AppendAdvancedRuleCluster(clusters)

		m, err := rm.storager.FetchRouteRules(ctx, []*ibasic.Product{product}, clusters)
		if err != nil {
			return err
		}
		prr = m[product.ID]

		return nil
	})

	return
}
