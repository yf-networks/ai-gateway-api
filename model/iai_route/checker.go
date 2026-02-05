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
	"fmt"
	"strings"

	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/model/iroute_conf"
)

// validatePathFilter validates the path filter
func validatePathFilter(pathFilter *PathFilter, ruleName string) error {
	if pathFilter == nil {
		return nil
	}

	// Validate match modes
	validMatchModes := map[string]bool{
		MatchModePrefix: true,
		MatchModeExact:  true,
		MatchModeSuffix: true,
	}

	if pathFilter.MatchMode != nil && *pathFilter.MatchMode != "" && !validMatchModes[*pathFilter.MatchMode] {
		return fmt.Errorf("Invalid path_filter.match_mode value for rule [%s]: %s", ruleName, *pathFilter.MatchMode)
	}

	// Validate path length
	if pathFilter.Path != nil && len(*pathFilter.Path) > 2048 {
		return fmt.Errorf("Path length for rule [%s] exceeds 2048 characters limit", ruleName)
	}

	return nil
}

// validateExpectForward validates the expected forward action
func validateExpectForward(forward *iroute_conf.ActionForward, ruleName string) error {
	if forward.ClusterName == "" {
		return fmt.Errorf("cluster_name cannot be empty for expect_action.forward in rule [%s]", ruleName)
	}
	if forward.URL != "" && !isValidURL(forward.URL) {
		return fmt.Errorf("Invalid URL format for expect_action.forward.url in rule [%s]: %s", ruleName, forward.URL)
	}
	return nil
}

// validateExpectAction validates the expected action
func validateExpectAction(action *iroute_conf.RouteAction, ruleName string) error {
	actionCount := 0
	if action.Forward != nil {
		actionCount++
		if err := validateExpectForward(action.Forward, ruleName); err != nil {
			return err
		}
	}

	if actionCount != 1 {
		return fmt.Errorf("Rule [%s] must contain exactly one of forward", ruleName)
	}

	return nil
}

// validateBasicInfo validates basic information (updated version)
func validateBasicInfo(basic *BasicInfo, ruleName string) error {
	if (basic.Domain == nil || *basic.Domain == "") &&
		len(basic.HeaderFilters) == 0 &&
		(basic.Method == nil || *basic.Method == "") &&
		basic.ModelFilter == nil &&
		basic.PathFilter == nil {
		return xerror.WrapParamErrorWithMsg("Must set only one parameter: basic.domain, basic.method, basic.header_filters, basic.model_filter or basic.path_filter")
	}

	if err := validatePathFilter(basic.PathFilter, ruleName); err != nil {
		return err
	}

	// Validate header filter array
	for i, headerFilter := range basic.HeaderFilters {
		if err := validateBasicHeaderFilter(headerFilter, i, ruleName); err != nil {
			return err
		}
	}

	// Validate expected action
	if basic.ExpectAction == nil {
		return fmt.Errorf("expect_action cannot be empty for rule [%s]", ruleName)
	}

	if err := validateExpectAction(basic.ExpectAction, ruleName); err != nil {
		return err
	}

	// Validate request method
	if basic.Method != nil {
		validMethods := map[string]bool{
			"GET": true, "POST": true, "DELETE": true,
			"PATCH": true, "PUT": true, "OPTIONS": true,
		}

		if !validMethods[strings.ToUpper(*basic.Method)] {
			return fmt.Errorf("Invalid method value for rule [%s]: %s", ruleName, *basic.Method)
		}
	}

	// Validate domain format
	if basic.Domain != nil && *basic.Domain != "" && !isValidDomain(*basic.Domain) {
		return fmt.Errorf("Invalid domain format for rule [%s]: %s", ruleName, *basic.Domain)
	}

	// Validate model filter
	if basic.ModelFilter != nil && basic.ModelFilter.Name != nil && *basic.ModelFilter.Name == "" {
		return fmt.Errorf("model_filter.name cannot be empty for rule [%s]", ruleName)
	}

	return nil
}

func isValidDomain(domain string) bool {
	// Simple domain validation
	// Actual implementation should be more complex
	return len(domain) > 0 && len(domain) <= 255
}

func isValidURL(url string) bool {
	// Check if contains scheme
	return strings.Contains(url, "://")
}

// Helper validation function
func isValidHeaderKey(key string) bool {
	// Only contains printable ASCII characters (0x21-0x7E)
	for _, char := range key {
		if char < 0x21 || char > 0x7E {
			return false
		}
		// Cannot contain spaces, colons (:), parentheses, etc.
		if char == ' ' || char == ':' || char == '(' || char == ')' {
			return false
		}
	}
	return true
}

// validateBasicHeaderFilter validates basic header filter (within BasicInfo)
func validateBasicHeaderFilter(filter *BasicHeaderFilter, index int, ruleName string) error {
	if filter == nil {
		return nil
	}

	if filter.Key == nil && filter.Value == nil {
		return nil
	}

	if filter.Key == nil || filter.Value == nil {
		return fmt.Errorf("Key or value is empty for header_filters[%d] in rule [%s]", index+1, ruleName)
	}

	// Validate Header Key
	if filter.Key != nil && *filter.Key == "" {
		return fmt.Errorf("Key cannot be empty for header_filters[%d] in rule [%s]", index+1, ruleName)
	}

	if !isValidHeaderKey(*filter.Key) {
		return fmt.Errorf("Invalid key format for header_filters[%d] in rule [%s]: %s", index+1, ruleName, *filter.Key)
	}

	// Validate Header Value
	if filter.Value != nil && *filter.Value == "" {
		return fmt.Errorf("Value cannot be empty for header_filters[%d] in rule [%s]", index+1, ruleName)
	}

	// Validate match mode

	validMatchModes := map[string]bool{
		MatchModePrefix: true,
		MatchModeExact:  true,
		MatchModeSuffix: true,
	}
	if filter.MatchMode != nil && *filter.MatchMode != "" && !validMatchModes[*filter.MatchMode] {
		return fmt.Errorf("Invalid match_mode value for header_filters[%d] in rule [%s]: %s", index+1, ruleName, *filter.MatchMode)
	}

	// Validate header value length
	if filter.Value != nil && len(*filter.Value) > 8192 { // Single header value ≤ 8KB
		return fmt.Errorf("Header value length for header_filters[%d] in rule [%s] exceeds 8KB limit", index+1, ruleName)
	}

	return nil
}

// ValidateRule validates a single rule (using concrete types)
func ValidateRule(rule *Rule, index int) error {
	if rule == nil {
		return fmt.Errorf("Rule %d cannot be empty", index+1)
	}

	// Validate rule name
	if rule.Name == "" {
		return fmt.Errorf("Rule name cannot be empty")
	}
	if !isValidRuleName(rule.Name) {
		return fmt.Errorf("Invalid rule name format: %s", rule.Name)
	}

	// Validate basic information
	if rule.Basic == nil {
		return fmt.Errorf("Basic information cannot be empty for rule [%s]", rule.Name)
	}
	if err := validateBasicInfo(rule.Basic, rule.Name); err != nil {
		return err
	}

	return nil
}

// isValidRuleName validates rule name format
func isValidRuleName(name string) bool {
	if len(name) < 1 || len(name) > 128 {
		return false
	}

	// Must start with letter or digit
	firstChar := name[0]
	if !((firstChar >= 'a' && firstChar <= 'z') ||
		(firstChar >= 'A' && firstChar <= 'Z') ||
		(firstChar >= '0' && firstChar <= '9')) {
		return false
	}

	// Allow digits, letters, underscores, and hyphens
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}

	return true
}
