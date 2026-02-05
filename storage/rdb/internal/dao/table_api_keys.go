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

package dao

import (
	"time"

	"github.com/yf-networks/ai-gateway-api/lib"
	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/storage/rdb/internal/dao/internal"
)

const tAPIKeyTableName = "api_keys"

type TAPIKey struct {
	ID            int64     `db:"id"`
	Name          string    `db:"name"`
	Enable        bool      `db:"enable"`
	Key           string    `db:"api_key"`
	IsLimit       bool      `db:"is_limit"`
	ProductName   string    `db:"product_name"`
	Limit         int64     `db:"total_quota"`
	ExpiredTime   string    `db:"expired_time"`
	AllowedModels string    `db:"allowed_models"`
	AllowedCIDR   string    `db:"allowed_cidr"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// TAPIKeyOne Query One
// return nil, nil if record not existed
func TAPIKeyOne(dbCtx lib.DBContexter, where *TAPIKeyParam) (*TAPIKey, error) {
	t := &TAPIKey{}
	err := internal.QueryOne(dbCtx, tAPIKeyTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TAPIKeyList Query Multiple
func TAPIKeyList(dbCtx lib.DBContexter, where *TAPIKeyParam) ([]*TAPIKey, error) {
	t := []*TAPIKey{}
	err := internal.QueryList(dbCtx, tAPIKeyTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

type TAPIKeyParam struct {
	ID *int64 `db:"id"`

	Name          *string    `db:"name"`
	Enable        *bool      `db:"enable"`
	Key           *string    `db:"api_key"`
	IsLimit       *bool      `db:"is_limit"`
	ProductName   *string    `db:"product_name"`
	Limit         *int64     `db:"total_quota"`
	ExpiredTime   *string    `db:"expired_time"`
	AllowedModels *string    `db:"allowed_models"`
	AllowedCIDR   *string    `db:"allowed_cidr"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TAPIKeyCreate One/Multiple
func TAPIKeyCreate(dbCtx lib.DBContexter, data ...*TAPIKeyParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tAPIKeyTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tAPIKeyTableName, list...)
}

// TAPIKeyUpdate Update One
func TAPIKeyUpdate(dbCtx lib.DBContexter, val, where *TAPIKeyParam) (int64, error) {
	return internal.Update(dbCtx, tAPIKeyTableName, where, val)
}

// TAPIKeyDelete Delete One/Multiple
func TAPIKeyDelete(dbCtx lib.DBContexter, where *TAPIKeyParam) (int64, error) {
	return internal.Delete(dbCtx, tAPIKeyTableName, where)
}
