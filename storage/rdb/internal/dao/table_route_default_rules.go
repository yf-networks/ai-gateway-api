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

const tRouteDefaultRuleTableName = "route_default_rules"

// TRouteDefaultRule Query Result
type TRouteDefaultRule struct {
	ID          int64     `db:"id"`
	Cmd         string    `db:"cmd"`
	Params      string    `db:"params"`
	ProductID   int64     `db:"product_id"`
	Description string    `db:"description"`
	RouteAction string    `db:"route_action"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// TRouteDefaultRuleParamCreate/Update/Where Data Carrier
// See: https://github.com/didi/gendry/blob/master/builder/README.md
type TRouteDefaultRuleParam struct {
	ID          *int64     `db:"id"`
	Cmd         *string    `db:"cmd"`
	Params      *string    `db:"params"`
	ProductID   *int64     `db:"product_id"`
	ProductIDs  []int64    `db:"product_id,in"`
	Description *string    `db:"description"`
	RouteAction *string    `db:"route_action"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`

	OrderBy *string `db:"_orderby"`
}

// TRouteDefaultRuleOne Query One
// return nil, nil if record not existed
func TRouteDefaultRuleOne(dbCtx lib.DBContexter, where *TRouteDefaultRuleParam) (*TRouteDefaultRule, error) {
	t := &TRouteDefaultRule{}
	err := internal.QueryOne(dbCtx, tRouteDefaultRuleTableName, where, t)
	if err == nil {
		return t, nil
	}
	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}
	return nil, err
}

// TRouteDefaultRuleList Query Multiple
func TRouteDefaultRuleList(dbCtx lib.DBContexter, where *TRouteDefaultRuleParam) ([]*TRouteDefaultRule, error) {
	t := []*TRouteDefaultRule{}
	err := internal.QueryList(dbCtx, tRouteDefaultRuleTableName, where, &t)
	if err == nil {
		return t, nil
	}

	if xerror.Cause(err) == internal.ErrRecordNotFound {
		return nil, nil
	}

	return nil, err
}

// TRouteDefaultRuleCreate One/Multiple
func TRouteDefaultRuleCreate(dbCtx lib.DBContexter, data ...*TRouteDefaultRuleParam) (int64, error) {
	if len(data) == 1 {
		if data[0].CreatedAt == nil {
			data[0].CreatedAt = internal.PTimeNow()
		}
		return internal.Create(dbCtx, tRouteDefaultRuleTableName, data[0])
	}

	list := make([]interface{}, len(data))
	for i, one := range data {
		if one.CreatedAt == nil {
			one.CreatedAt = internal.PTimeNow()
		}
		list[i] = one
	}

	return internal.Create(dbCtx, tRouteDefaultRuleTableName, list...)
}

// TRouteDefaultRuleUpdate Update One
func TRouteDefaultRuleUpdate(dbCtx lib.DBContexter, val, where *TRouteDefaultRuleParam) (int64, error) {
	return internal.Update(dbCtx, tRouteDefaultRuleTableName, where, val)
}

// TRouteDefaultRuleDelete Delete One/Multiple
func TRouteDefaultRuleDelete(dbCtx lib.DBContexter, where *TRouteDefaultRuleParam) (int64, error) {
	return internal.Delete(dbCtx, tRouteDefaultRuleTableName, where)
}
