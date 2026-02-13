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

package icluster_conf

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yf-networks/ai-gateway-api/lib"
	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/model/itxn"
	"github.com/yf-networks/ai-gateway-api/stateful"
)

// APIKeyParam defines the parameters for API key operations
type APIKeyParam struct {
	Name        *string    `json:"name"`
	Enable      *bool      `json:"enable"`
	Status      *string    `json:"status,omitempty"`
	UpdatedTime *string    `json:"updated_time,omitempty"`
	KeyCreateAt *time.Time `json:"-"`

	// Key is the actual API key string, format: product line name + multiple randomly generated segments
	Key *string `json:"key"`

	// IsLimit indicates whether quota limitation is enabled
	IsLimit *bool `json:"is_limit"`

	// Limit is the specific quota limit, required when IsLimit is true, range: 0-100000000
	Limit *int64 `json:"total_quota,omitempty"`

	// ExpiredTime defines the expiration time with formats:
	// Empty string: Never expires
	// "1m": One month later
	// "1d": One day later
	// "1h": One hour later
	// Or timestamp string: "2025-01-01 01:01:01"
	ExpiredTime    *string  `json:"expired_time,omitempty"`
	AllowedModels  []string `json:"allowed_models,omitempty"`
	AllowedCIDR    []string `json:"allowed_subnets,omitempty"`
	ProductName    *string  `json:"-"`
	ID             *int64   `json:"-"`
	RemainingQuota *int64   `json:"remaining_quota,omitempty"`
}

// APIKeyTokenParam defines parameters for API key token operations
type APIKeyTokenParam struct {
	Key       *string
	CreatedAt *time.Time
}

// APIKeyTokenFilter defines filters for querying API key tokens
type APIKeyTokenFilter struct {
	Key *string
	ID  *int64
}

// APIKeyFilter defines filters for querying API keys
type APIKeyFilter struct {
	ProductName  *string
	ProductNames []string
	Name         *string
	ALBGroupName *string
	ID           *int64
}

// APIKeyStorager interface defines storage operations for API keys
type APIKeyStorager interface {
	FetchAPIKeyList(ctx context.Context, filter *APIKeyFilter) ([]*APIKeyParam, error)
	CreateAPIKey(ctx context.Context, param *APIKeyParam) (int64, error)
	UpdateAPIKey(ctx context.Context, filter *APIKeyFilter, param *APIKeyParam) (int64, error)
	DeleteAPIKey(ctx context.Context, filter *APIKeyFilter) error

	CreateAPIKeyToken(ctx context.Context, param *APIKeyTokenParam) (int64, error)
	UpdateAPIKeyToken(ctx context.Context, filter *APIKeyTokenFilter, param *APIKeyTokenParam) error
	FetchAPIKeyTokenList(ctx context.Context, filter *APIKeyTokenFilter) ([]*APIKeyTokenParam, error)
}

// APIKeyManager manages API key operations with transaction support
type APIKeyManager struct {
	storager        APIKeyStorager
	txn             itxn.TxnStorager
	clusterStorager ClusterStorager
}

// NewAPIKeyManager creates a new APIKeyManager instance
func NewAPIKeyManager(txn itxn.TxnStorager, storager APIKeyStorager, clusterStorager ClusterStorager) *APIKeyManager {
	return &APIKeyManager{
		txn:             txn,
		storager:        storager,
		clusterStorager: clusterStorager,
	}
}

// GetRemainingQuota calculates the remaining quota for an API key
func GetRemainingQuota(param *APIKeyParam) (*int64, error) {
	// Retrieve used quota from Redis cache
	used, err := stateful.DefaultClientSet.RedisClient.GetInt64(stateful.AIUsedQuotaKey(*param.Key, param.KeyCreateAt.Unix()))
	if err != nil {
		if strings.Contains(err.Error(), "redigo: nil returned") {
			// If no usage record exists, return the full limit
			return lib.PInt64(*param.Limit), nil
		}

		return nil, fmt.Errorf("get %s-%d from cache is error:%s", *param.Key, param.KeyCreateAt.Unix(), err.Error())
	}

	// Calculate remaining quota
	if param.Limit != nil {
		if *param.Limit > used {
			return lib.PInt64(*param.Limit - used), nil
		}
		return nil, nil
	}

	// No limit set
	return nil, nil
}

// FetchAPIKeyList retrieves API keys based on filter criteria
func (rppm *APIKeyManager) FetchAPIKeyList(ctx context.Context,
	filter *APIKeyFilter) (list []*APIKeyParam, err error) {
	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err = rppm.storager.FetchAPIKeyList(ctx, filter)
		return err
	})

	return
}

// FetchAPIKey retrieves a single API key based on filter criteria
func (rppm *APIKeyManager) FetchAPIKey(ctx context.Context,
	filter *APIKeyFilter) (one *APIKeyParam, err error) {
	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		list, err := rppm.storager.FetchAPIKeyList(ctx, filter)
		if err != nil {
			return err
		}
		if len(list) > 0 {
			one = list[0]
		}
		return nil
	})

	return
}

// DeleteAPIKey deletes an API key based on filter criteria
func (rppm *APIKeyManager) DeleteAPIKey(ctx context.Context, filter *APIKeyFilter) error {
	return rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		// Verify the API key exists before deletion
		list, err := rppm.storager.FetchAPIKeyList(ctx, filter)
		if err != nil {
			return err
		}
		if len(list) == 0 {
			return xerror.WrapRecordNotExist("APIKey")
		}

		return rppm.storager.DeleteAPIKey(ctx, filter)
	})
}

// UpdateAPIKey updates an existing API key
func (rppm *APIKeyManager) UpdateAPIKey(ctx context.Context, filter *APIKeyFilter, param *APIKeyParam) error {
	return rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		// Verify the API key exists before update
		list, err := rppm.storager.FetchAPIKeyList(ctx, filter)
		if err != nil {
			return err
		}
		if len(list) == 0 {
			return xerror.WrapRecordNotExist("API-Key")
		}

		one := list[0]
		// Skip quota update if the limit value remains unchanged
		if param.Enable != nil && *param.Enable && param.IsLimit != nil && *param.IsLimit &&
			param.Limit != nil && *param.Limit > 0 && one.Limit != nil && *param.Limit == *one.Limit {
			param.Limit = nil
		}

		_, err = rppm.storager.UpdateAPIKey(ctx, &APIKeyFilter{
			ID: list[0].ID,
		}, param)
		return err
	})
}

// CreateAPIKey creates a new API key
func (rppm *APIKeyManager) CreateAPIKey(ctx context.Context,
	param *APIKeyParam) (err error) {
	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		// Check for duplicate API key name within the same product
		list, err := rppm.storager.FetchAPIKeyList(ctx, &APIKeyFilter{
			Name:        param.Name,
			ProductName: param.ProductName,
		})
		if err != nil {
			return err
		}

		if len(list) > 0 {
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Duplicate name with product:%s", *param.ProductName))
		}

		// Check for existing API key tokens
		tokens, err := rppm.storager.FetchAPIKeyTokenList(ctx, &APIKeyTokenFilter{Key: param.Key})
		if err != nil {
			return err
		}
		if len(tokens) > 1 {
			return xerror.WrapDirtyDataErrorWithMsg(fmt.Sprintf("API-Key-Token:%s", *param.Key))
		}

		// Set updated time based on existing token or current time
		if len(tokens) > 0 {
			param.UpdatedTime = lib.PString(tokens[0].CreatedAt.Format(lib.FormatTimeYYMMDD_HHMMSS))
		} else {
			param.UpdatedTime = lib.PString(time.Now().Format(lib.FormatTimeYYMMDD_HHMMSS))
		}

		_, err = rppm.storager.CreateAPIKey(ctx, param)
		return err
	})

	return
}

// CreateAPIKeyToken creates a new API key token
func (rppm *APIKeyManager) CreateAPIKeyToken(ctx context.Context,
	param *APIKeyTokenParam) (id int64, err error) {
	err = rppm.txn.AtomExecute(ctx, func(ctx context.Context) error {
		id, err = rppm.storager.CreateAPIKeyToken(ctx, param)
		if err != nil {
			return err
		}

		// Update the token with the full key format: key-id
		return rppm.storager.UpdateAPIKeyToken(ctx, &APIKeyTokenFilter{ID: &id}, &APIKeyTokenParam{
			Key: lib.PString(fmt.Sprintf("%s-%d", *param.Key, id)),
		})
	})

	return
}
