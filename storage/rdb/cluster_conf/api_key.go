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

package cluster_conf

import (
	"context"
	"encoding/json"

	"github.com/yf-networks/ai-gateway-api/lib"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/storage/rdb/internal/dao"
)

type APIKeyStorager struct {
	dbCtxFactory lib.DBContextFactory
}

func NewAPIKeyStorager(dbCtxFactory lib.DBContextFactory) *APIKeyStorager {
	return &APIKeyStorager{
		dbCtxFactory: dbCtxFactory,
	}
}

var _ icluster_conf.APIKeyStorager = &APIKeyStorager{}

func (rpps *APIKeyStorager) CreateAPIKey(ctx context.Context,
	param *icluster_conf.APIKeyParam) (int64, error) {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return 0, err
	}

	data := newAPIKeyDataToParam(param)

	allowedModels := make([]string, 0)
	if len(param.AllowedModels) > 0 {
		allowedModels = param.AllowedModels
	}
	allowedModelsValue, _ := json.Marshal(allowedModels)
	data.AllowedModels = lib.PString(string(allowedModelsValue))

	data.CreatedAt = lib.PTimeNow()

	allowedSubnets := make([]string, 0)
	if len(param.AllowedCIDR) > 0 {
		allowedSubnets = param.AllowedCIDR
	}
	allowedSubnetsValue, _ := json.Marshal(allowedSubnets)
	data.AllowedCIDR = lib.PString(string(allowedSubnetsValue))

	return dao.TAPIKeyCreate(dbCtx, data)
}

func newAPIKeyDataToParam(param *icluster_conf.APIKeyParam) *dao.TAPIKeyParam {
	data := &dao.TAPIKeyParam{
		Name:        param.Name,
		Enable:      param.Enable,
		Key:         param.Key,
		IsLimit:     param.IsLimit,
		Limit:       param.Limit,
		ExpiredTime: param.ExpiredTime,
		ProductName: param.ProductName,
		UpdatedAt:   lib.PTimeNow(),
	}

	return data
}

func newAPIKeyFilterToParam(filter *icluster_conf.APIKeyFilter) *dao.TAPIKeyParam {
	if filter == nil {
		return nil
	}

	return &dao.TAPIKeyParam{
		ProductName: filter.ProductName,
		Name:        filter.Name,
		ID:          filter.ID,
	}
}

func (rpps *APIKeyStorager) FetchAPIKeyList(ctx context.Context,
	filter *icluster_conf.APIKeyFilter) ([]*icluster_conf.APIKeyParam, error) {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TAPIKeyList(dbCtx, newAPIKeyFilterToParam(filter))
	if err != nil {
		return nil, err
	}

	var rst []*icluster_conf.APIKeyParam
	for _, one := range list {
		rst = append(rst, apiKeyParamToData(one))
	}

	return rst, nil
}

func apiKeyParamToData(one *dao.TAPIKey) *icluster_conf.APIKeyParam {
	allowedModels := make([]string, 0)
	if one.AllowedModels != "" {
		json.Unmarshal([]byte(one.AllowedModels), &allowedModels)
	}

	allowedSubnets := make([]string, 0)
	if one.AllowedCIDR != "" {
		json.Unmarshal([]byte(one.AllowedCIDR), &allowedSubnets)
	}

	return &icluster_conf.APIKeyParam{
		ID:            &one.ID,
		Name:          &one.Name,
		Enable:        &one.Enable,
		Key:           &one.Key,
		IsLimit:       &one.IsLimit,
		Limit:         &one.Limit,
		ExpiredTime:   &one.ExpiredTime,
		AllowedModels: allowedModels,
		AllowedCIDR:   allowedSubnets,
		ProductName:   &one.ProductName,
		KeyCreateAt:   &one.CreatedAt,
		UpdatedTime:   lib.PString(one.CreatedAt.Format(lib.FormatTimeYYMMDD_HHMMSS)),
	}
}

func (rpps *APIKeyStorager) DeleteAPIKey(ctx context.Context, filter *icluster_conf.APIKeyFilter) error {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TAPIKeyDelete(dbCtx, newAPIKeyFilterToParam(filter))

	return err
}

func (rpps *APIKeyStorager) UpdateAPIKey(ctx context.Context, filter *icluster_conf.APIKeyFilter, param *icluster_conf.APIKeyParam) (int64, error) {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return 0, err
	}

	data := newAPIKeyDataToParam(param)

	// 需要更新额度，需要重新生成AI模型的更新时间
	if data.Enable != nil && *data.Enable {
		if data.IsLimit != nil && *data.IsLimit && data.Limit != nil && *data.Limit > 0 {
			data.CreatedAt = lib.PTimeNow()
		}
	}

	allowedSubnets := make([]string, 0)
	if len(param.AllowedCIDR) > 0 {
		allowedSubnets = param.AllowedCIDR
	}
	allowedSubnetsValue, _ := json.Marshal(allowedSubnets)
	data.AllowedCIDR = lib.PString(string(allowedSubnetsValue))

	return dao.TAPIKeyUpdate(dbCtx, data, newAPIKeyFilterToParam(filter))
}

func (rpps *APIKeyStorager) CreateAPIKeyToken(ctx context.Context,
	param *icluster_conf.APIKeyTokenParam) (int64, error) {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return 0, err
	}

	return dao.TAPIKeyTokenCreate(dbCtx, &dao.TAPIKeyTokenParam{
		Key:       param.Key,
		CreatedAt: lib.PTimeNow(),
		UpdatedAt: lib.PTimeNow(),
	})
}

func (rpps *APIKeyStorager) UpdateAPIKeyToken(ctx context.Context, filter *icluster_conf.APIKeyTokenFilter, param *icluster_conf.APIKeyTokenParam) error {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return err
	}

	_, err = dao.TAPIKeyTokenUpdate(dbCtx, &dao.TAPIKeyTokenParam{
		Key: param.Key,
	}, &dao.TAPIKeyTokenParam{ID: filter.ID})
	return err
}

func (rpps *APIKeyStorager) FetchAPIKeyTokenList(ctx context.Context,
	filter *icluster_conf.APIKeyTokenFilter) ([]*icluster_conf.APIKeyTokenParam, error) {
	dbCtx, err := rpps.dbCtxFactory(ctx)
	if err != nil {
		return nil, err
	}

	list, err := dao.TAPIKeyTokenList(dbCtx, &dao.TAPIKeyTokenParam{
		Key: filter.Key,
	})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return newAPIKeyTokensDataToParam(list), nil
}

func newAPIKeyTokensDataToParam(list []*dao.TAPIKeyToken) []*icluster_conf.APIKeyTokenParam {
	results := make([]*icluster_conf.APIKeyTokenParam, len(list))
	for i, one := range list {
		results[i] = &icluster_conf.APIKeyTokenParam{
			Key:       &one.Key,
			CreatedAt: &one.CreatedAt,
		}
	}

	return results
}
