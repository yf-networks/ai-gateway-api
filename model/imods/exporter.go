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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bfenetworks/bfe/bfe_modules/mod_ai_token_auth"
	"github.com/yf-networks/ai-gateway-api/lib"
	"github.com/yf-networks/ai-gateway-api/model/iai_route"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/model/iversion_control"
	"github.com/yf-networks/ai-gateway-api/stateful"
)

// ConfigTopicProductAPIKeyRule is the configuration topic for API key rules
const ConfigTopicProductAPIKeyRule = "mod_api_key_rule"

// ModAPIKeyRuleConf defines the configuration structure for API key rules module
type ModAPIKeyRuleConf struct {
	Version *string                             `json:"version"`
	Config  map[string][]*ExportAPIKeyRule      `json:"config"`
	Tokens  map[string]map[string]ExportContent `json:"tokens"`
}

// ExportContent defines the structure for API key information exported to BFE
type ExportContent struct {
	Key            string `json:"key"`             // API key value
	Status         int    `json:"status"`          // Key status (enabled/disabled/expired/exhausted)
	Name           string `json:"name"`            // API key name
	UpdatedTime    int64  `json:"update_time"`     // Last update timestamp
	ExpiredTime    int64  `json:"expired_time"`    // Expiration timestamp
	UnlimitedQuota bool   `json:"unlimited_quota"` // Whether quota is unlimited
	RemainQuota    int64  `json:"remain_quota"`    // Remaining quota
	Models         string `json:"models"`          // Allowed models (comma-separated)
	Subnet         string `json:"subnet"`          // Allowed subnets (comma-separated)
}

// ExportAPIKeyRule defines the structure for API key routing rules exported to BFE
type ExportAPIKeyRule struct {
	Cond   *string             `json:"Cond"`   // Routing condition
	Action *ExportAPIKeyAction `json:"action"` // Routing action
}

// ExportAPIKeyAction defines the action for API key routing rules
type ExportAPIKeyAction struct {
	CMD *string `json:"cmd"` // Command to execute
}

// UpdateVersion updates the configuration version
func (conf *ModAPIKeyRuleConf) UpdateVersion(version string) error {
	conf.Version = &version
	return nil
}

// ConfigExport exports API key rule configuration for BFE
func (rcm *APIKeyRuleManager) ConfigExport(ctx context.Context, lastVersion string) (*ModAPIKeyRuleConf, error) {
	// Export configuration using version control manager
	rst, err := rcm.versionControlManager.ExportConfig(ctx, ConfigTopicProductAPIKeyRule, rcm.APIKeyRuleGenerator)
	if err != nil {
		return nil, err
	}

	if rst.DataWithoutVersion == nil {
		return nil, fmt.Errorf("APIKeyRuleGenerator.DataWithoutVersion is nil")
	}

	// Convert exported data to configuration structure
	conf, ok := rst.DataWithoutVersion.(*ModAPIKeyRuleConf)
	if ok {
		// Return nil if version hasn't changed
		if *conf.Version == lastVersion {
			return nil, nil
		}

		return conf, nil
	}

	return nil, fmt.Errorf("convert APIKeyRuleGenerator.DataWithoutVersion to ModAPIKeyRuleConf is error")
}

func (rlm *APIKeyRuleManager) FormatAIRouteAPIKeyRules(ctx context.Context, productName string) ([]*APIKeyRule, error) {
	aiRouteRules, err := rlm.aiRouteStorager.FetchAIRouteRules(ctx, &iai_route.AIRouteFilter{
		ProductName: &productName,
	})
	if err != nil {
		return nil, err
	}
	if len(aiRouteRules) == 0 {
		return nil, nil
	}

	var ruleResult []*APIKeyRule
	for _, rule := range aiRouteRules {
		cond := iai_route.BuildAIRouteCond(ctx, rule.Basic)

		ruleResult = append(ruleResult, &APIKeyRule{
			Cond: cond,
			Actions: []Action{
				{
					Cmd: APIKeyActionCMD,
				},
			},
			ProductName: stateful.DefaultConfig.RunTime.AIRouteInnerProductName,
		})
	}
	return ruleResult, nil
}

func (rlm *APIKeyRuleManager) buildAIRouteAPIKeyRules(ctx context.Context) (map[string]*APIKey, error) {
	aiProductName := stateful.DefaultConfig.RunTime.AIRouteInnerProductName
	rules, err := rlm.FormatAIRouteAPIKeyRules(ctx, aiProductName)
	if err != nil {
		return nil, fmt.Errorf("read ai route api key rule is error:%s", err.Error())
	}

	product2config := make(map[string]*APIKey)
	if len(rules) > 0 {
		product2config[aiProductName] = &APIKey{
			Rules: rules,
		}
	}

	return product2config, nil
}

// APIKeyRuleGenerator generates API key rules and token information for BFE configuration
func (rlm *APIKeyRuleManager) APIKeyRuleGenerator(ctx context.Context) (*iversion_control.ExportData, error) {
	// Fetch API key rules from storage
	product2config, err := rlm.buildAIRouteAPIKeyRules(ctx)
	if err != nil {
		return nil, err
	}

	// Convert rules to BFE format
	productName2Config := make(map[string][]*ExportAPIKeyRule)
	for productName, productConfig := range product2config {
		if len(productConfig.Rules) > 0 {
			productName2Config[productName] = convertAPIKeyRulesToBfeRules(productConfig.Rules)
		}
	}

	// Fetch API key list from storage
	apiKeyList, err := rlm.apiKeyStorager.FetchAPIKeyList(ctx, &icluster_conf.APIKeyFilter{})
	if err != nil {
		return nil, err
	}

	// Build token configuration for each product
	apiKey2Config := make(map[string]map[string]ExportContent)
	for _, one := range apiKeyList {
		// Initialize product map if not exists
		if _, ok := apiKey2Config[*one.ProductName]; !ok {
			items := make(map[string]ExportContent)
			apiKey2Config[*one.ProductName] = items
		}

		// Parse expiration time
		expiredTime := int64(UnlimitedQuota) // Default to unlimited
		if one.ExpiredTime != nil && *one.ExpiredTime != "" {
			t, err := time.ParseInLocation(lib.FormatTimeYYMMDD_HHMMSS, *one.ExpiredTime, time.Local)
			if err != nil {
				return nil, fmt.Errorf("parse expired time:%s is error:%s", *one.ExpiredTime, err.Error())
			}
			expiredTime = t.Unix()
		}

		// Determine key status and quota
		limit := int64(0)
		status := mod_ai_token_auth.TokenStatusDisabled
		if one.Enable != nil && *one.Enable {
			if *one.IsLimit {
				limit = *one.Limit
			}

			status = mod_ai_token_auth.TokenStatusEnabled
			remainingQuota, err := icluster_conf.GetRemainingQuota(one)
			if err != nil {
				return nil, err
			}

			if remainingQuota != nil {
				if *remainingQuota == 0 {
					status = mod_ai_token_auth.TokenStatusExhausted
				} else {
					// Check if key has expired
					if expiredTime != int64(UnlimitedQuota) && time.Now().Local().Unix() >= expiredTime {
						status = mod_ai_token_auth.TokenStatusExpired
					}
				}
			} else {
				status = mod_ai_token_auth.TokenStatusDisabled
			}
		}

		// Build export content
		items := apiKey2Config[*one.ProductName]
		ec := ExportContent{
			Key:         *one.Key,
			Status:      status,
			Name:        *one.Name,
			ExpiredTime: expiredTime,
			UpdatedTime: one.KeyCreateAt.Unix(),
		}

		// Set quota information
		if *one.IsLimit {
			ec.UnlimitedQuota = false
			ec.RemainQuota = limit
		} else {
			ec.UnlimitedQuota = true
		}

		// Convert allowed models to comma-separated string
		if len(one.AllowedModels) > 0 {
			ec.Models = strings.Join(one.AllowedModels, ",")
		}

		// Convert allowed CIDR to comma-separated string
		if len(one.AllowedCIDR) > 0 {
			ec.Subnet = strings.Join(one.AllowedCIDR, ",")
		}

		// Add to configuration
		items[*one.Key] = ec
		apiKey2Config[*one.ProductName] = items
	}

	// Build final configuration
	conf := &ModAPIKeyRuleConf{
		Config: productName2Config,
		Tokens: apiKey2Config,
	}

	// Set version to zero (will be updated by version control system)
	conf.UpdateVersion(iversion_control.ZeroVersion)

	return &iversion_control.ExportData{
		Topic:              ConfigTopicProductAPIKeyRule,
		DataWithoutVersion: conf,
	}, nil
}

// convertAPIKeyRulesToBfeRules converts internal API key rules to BFE format
func convertAPIKeyRulesToBfeRules(oldRules []*APIKeyRule) []*ExportAPIKeyRule {
	exportRules := make([]*ExportAPIKeyRule, len(oldRules))
	for i, rule := range oldRules {
		newRule := &ExportAPIKeyRule{
			Cond: &rule.Cond,
		}

		// Include action if present
		if len(rule.Actions) > 0 {
			newRule.Action = &ExportAPIKeyAction{
				CMD: &rule.Actions[0].Cmd,
			}
		}

		exportRules[i] = newRule
	}
	return exportRules
}
