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

const tAIRouteRuleTableName = "ai_route_rules"

type TAIRouteRule struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Basic       string    `db:"basic"`
	IDX         int64     `db:"idx"`
	ProductName string    `db:"product_name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TAIRouteRuleOne Query One
// return nil, nil if record not existed
func TAIRouteRuleOne(dbCtx lib.DBContexter, where *TAIRouteRuleParam) (*TAIRouteRule, error) {
	t := &TAIRouteRule{}
	err := internal.QueryOne(dbCtx, tAIRouteRuleTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TAIRouteRuleList Query Multiple
func TAIRouteRuleList(dbCtx lib.DBContexter, where *TAIRouteRuleParam) ([]*TAIRouteRule, error) {
	t := []*TAIRouteRule{}
	err := internal.QueryList(dbCtx, tAIRouteRuleTableName, where, &t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

type TAIRouteRuleParam struct {
	ID *int64 `db:"id"`

	Name        *string    `db:"name"`
	Basic       *string    `db:"basic"`
	IDX         *int64     `db:"idx"`
	ProductName *string    `db:"product_name"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TAIRouteRuleCreate One/Multiple
func TAIRouteRuleCreate(dbCtx lib.DBContexter, data ...*TAIRouteRuleParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tAIRouteRuleTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tAIRouteRuleTableName, list...)
}

// TAIRouteRuleUpdate Update One
func TAIRouteRuleUpdate(dbCtx lib.DBContexter, val, where *TAIRouteRuleParam) (int64, error) {
	return internal.Update(dbCtx, tAIRouteRuleTableName, where, val)
}

// TAIRouteRuleDelete Delete One/Multiple
func TAIRouteRuleDelete(dbCtx lib.DBContexter, where *TAIRouteRuleParam) (int64, error) {
	return internal.Delete(dbCtx, tAIRouteRuleTableName, where)
}
