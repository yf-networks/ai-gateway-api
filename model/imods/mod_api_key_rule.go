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

package imods

import (
	"github.com/yf-networks/ai-gateway-api/model/iai_route"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/model/itxn"
	"github.com/yf-networks/ai-gateway-api/model/iversion_control"
)

type APIKeyRuleManager struct {
	txn                   itxn.TxnStorager
	versionControlManager *iversion_control.VersionControlManager
	apiKeyStorager        icluster_conf.APIKeyStorager
	aiRouteStorager       iai_route.AIRouteRuleStorager
}

const (
	UnlimitedQuota = -1
)

type Action struct {
	Cmd    string   `json:"cmd" validate:"omitempty"`    // command of action
	Params []string `json:"params" validate:"omitempty"` // params of action
}

type APIKeyRule struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Cond        string   `json:"cond"`
	Actions     []Action `json:"actions"`
	ProductName string   `json:"-"`
	IDX         int64    `json:"-"`
}

type APIKeyRspData struct {
	Rules []*APIKeyRule `json:"rules"`
}

type APIKey struct {
	Rules []*APIKeyRule `json:"rules"`
}

const (
	APIKeyActionCMD = "CHECK_TOKEN"
)

func NewAPIKeyRuleManager(txn itxn.TxnStorager,
	versionControlManager *iversion_control.VersionControlManager,
	apiKeyStorager icluster_conf.APIKeyStorager,
	aiRouteStorager iai_route.AIRouteRuleStorager) *APIKeyRuleManager {
	return &APIKeyRuleManager{
		txn:                   txn,
		versionControlManager: versionControlManager,
		apiKeyStorager:        apiKeyStorager,
		aiRouteStorager:       aiRouteStorager,
	}
}
