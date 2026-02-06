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
	"github.com/bfenetworks/bfe/bfe_basic/condition"
	"github.com/yf-networks/ai-gateway-api/model/iroute_conf"
)

const (
	MatchModePrefix = "prefix_match"
	MatchModeExact  = "exact_match"
	MatchModeSuffix = "suffix_match"
)

// Rule
type Rule struct {
	Name         string              `json:"name"`
	Basic        *BasicInfo          `json:"basic"`
	ConditionVar condition.Condition `json:"-" uri:"-"`
	ProductName  string              `json:"-"`
}

// BasicInfo
type BasicInfo struct {
	Domain        *string                  `json:"domain,omitempty"`
	PathFilter    *PathFilter              `json:"path_filter"`
	Method        *string                  `json:"method,omitempty"`
	HeaderFilters []*BasicHeaderFilter     `json:"header_filters,omitempty"`
	ModelFilter   *ModelFilter             `json:"model_filter,omitempty"`
	ExpectAction  *iroute_conf.RouteAction `json:"expect_action"`
}

// BasicHeaderFilter
type BasicHeaderFilter struct {
	Key        *string `json:"key,omitempty"`
	Value      *string `json:"value,omitempty"`
	MatchMode  *string `json:"match_mode,omitempty"`
	IgnoreCase *bool   `json:"ignore_case"`
}

// ModelFilter
type ModelFilter struct {
	Name       *string `json:"name,omitempty"`
	Pattern    *string `json:"pattern,omitempty"`
	IgnoreCase *bool   `json:"ignore_case"`
}

// PathFilter
type PathFilter struct {
	MatchMode  *string `json:"match_mode,omitempty"`
	IgnoreCase *bool   `json:"ignore_case"`
	Path       *string `json:"path,omitempty"`
}

// HeaderMap
type HeaderMap map[string]string
