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

package api_key

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/yf-networks/ai-gateway-api/lib"
	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
)

const (
	maxLimit   = 100000000 // Maximum allowed quota limit
	maxNameLen = 255       // Maximum length for API key name
)

// checkCreateAPIKey validates parameters for creating a new API key
func checkCreateAPIKey(param *icluster_conf.APIKeyParam, productName string) error {
	if err := checkName(param.Name); err != nil {
		return err
	}

	if err := checkAllowSubnet(param.AllowedCIDR); err != nil {
		return err
	}

	if param.IsLimit == nil {
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Must set is_limit"))
	}
	if err := checkLimit(param.Limit, param.IsLimit); err != nil {
		return err
	}

	if err := checkKey(param.Key, productName); err != nil {
		return err
	}

	if err := checkExpiredTime(param.ExpiredTime); err != nil {
		return err
	}

	return nil
}

// checkName validates the API key name
func checkName(name *string) error {
	if name == nil || *name == "" {
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Must set name"))
	}

	if len(*name) > maxNameLen {
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("name must between 0 and %s", strconv.Itoa(maxNameLen)))
	}

	return nil
}

// checkExpiredTime validates the expiration time format
func checkExpiredTime(expiredTime *string) error {
	if expiredTime == nil || *expiredTime == "" {
		return nil
	}

	_, err := time.Parse(lib.FormatTimeYYMMDD_HHMMSS, *expiredTime)
	return err
}

// checkUpdateAPIKey validates parameters for updating an existing API key
func checkUpdateAPIKey(param *icluster_conf.APIKeyParam, productName string) error {
	if err := checkAllowSubnet(param.AllowedCIDR); err != nil {
		return err
	}

	if err := checkLimit(param.Limit, param.IsLimit); err != nil {
		return err
	}

	if param.Key != nil {
		if !strings.HasPrefix(*param.Key, fmt.Sprintf("%s-", productName)) {
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("key must begin with prefix %s-", productName))
		}

		if !ValidateString(*param.Key) {
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Allowed Characters: Uppercase/Lowercase Letters, Numbers, and Hyphen (-)"))
		}
	}

	if err := checkExpiredTime(param.ExpiredTime); err != nil {
		return err
	}

	return nil
}

// checkKey validates the API key value
func checkKey(key *string, productName string) error {
	if key == nil || *key == "" {
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Must set key"))
	}

	if !strings.HasPrefix(*key, fmt.Sprintf("%s-", productName)) {
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("key must begin with prefix %s-", productName))
	}

	if !ValidateString(*key) {
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Allowed Characters: Uppercase/Lowercase Letters, Numbers, and Hyphen (-)"))
	}

	return nil
}

// ValidateString checks if a string contains only allowed characters
func ValidateString(s string) bool {
	for _, r := range s {
		if !isValidChar(r) {
			return false
		}
	}
	return true
}

// isValidChar checks if a rune is a valid character for API key
func isValidChar(c rune) bool {
	return (c >= 65 && c <= 90) || // A-Z
		(c >= 97 && c <= 122) || // a-z
		(c >= 48 && c <= 57) || // 0-9
		(c == 45) || // Hyphen (-)
		(c == '_') // Underscore (_)
}

// checkLimit validates the quota limit parameters
func checkLimit(limit *int64, isLimit *bool) error {
	if isLimit == nil || !*isLimit {
		// If not limiting, limit should be nil
		if limit == nil {
			return nil
		}
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Invalid total_quota"))
	}

	// If limiting, limit must be set and within valid range
	if limit == nil {
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Must set total_quota"))
	}

	if *limit >= 0 && *limit <= maxLimit {
		return nil
	}

	return xerror.WrapParamErrorWithMsg(fmt.Sprintf("total_quota must be between 0 and %d", maxLimit))
}

// checkAllowSubnet validates CIDR subnet format
func checkAllowSubnet(cidrs []string) error {
	for _, cidr := range cidrs {
		arr := strings.Split(cidr, "/")
		if len(arr) != 2 {
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("invalid subnet format:%s", cidr))
		}

		ip := net.ParseIP(arr[0])
		if ip == nil {
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("invalid subnet:%s", cidr))
		}

		_, _, err := net.ParseCIDR(ip.String() + "/" + arr[1])
		if err != nil {
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("invalid subnet:%s parse error:%s", cidr, err.Error()))
		}
	}

	return nil
}
