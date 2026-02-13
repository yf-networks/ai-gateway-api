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
//limitations under the License. All rights reserved.

// Copyright (c) 2021 The BFE Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package subcluster

import (
	"net/http"

	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/lib/xreq"
	"github.com/yf-networks/ai-gateway-api/model/iauth"
	"github.com/yf-networks/ai-gateway-api/model/ibasic"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
	"github.com/yf-networks/ai-gateway-api/stateful/container"
)

// OneData Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneData struct {
	Name         string                   `json:"name" uri:"name"`
	InstancePool string                   `json:"instance_pool" uri:"instance_pool"`
	Description  string                   `json:"description" uri:"description"`
	Ready        bool                     `json:"ready" uri:"ready"`
	ProductName  string                   `json:"product_name,omitempty"`
	Instances    []icluster_conf.Instance `json:"instances"`
}

// OneParam Request Param
// AUTO GEN BY ctrl, MODIFY AS U NEED
type OneParam struct {
	SubClusterName string `json:"sub_cluster_name" uri:"sub_cluster_name" validate:"required,min=2"`
}

// OneRoute route
// AUTO GEN BY ctrl, MODIFY AS U NEED
var OneEndpoint = &xreq.Endpoint{
	Path:       "/products/{product_name}/sub-clusters/{sub_cluster_name}",
	Method:     http.MethodGet,
	Handler:    xreq.Convert(OneAction),
	Authorizer: iauth.FAP(iauth.FeatureSubCluster, iauth.ActionRead),
}

var _ xreq.Handler = OneAction

// OneAction action
// AUTO GEN BY ctrl, MODIFY AS U NEED
func OneAction(req *http.Request) (interface{}, error) {
	param, err := newOneParam4One(req)
	if err != nil {
		return nil, err
	}

	return oneActionProcess(req, param)
}

// AUTO GEN BY ctrl, MODIFY AS U NEED
func newOneParam4One(req *http.Request) (*OneParam, error) {
	param := &OneParam{}
	err := xreq.BindURI(req, param)
	return param, err
}

func oneActionProcess(req *http.Request, param *OneParam) (*OneData, error) {
	// get product info
	product, err := ibasic.MustGetProduct(req.Context())
	if err != nil {
		return nil, err
	}

	subCluster, err := one(req, &icluster_conf.SubClusterFilter{
		Names:   []string{param.SubClusterName},
		Product: product,
	})
	if err != nil {
		return nil, err
	}

	if subCluster == nil {
		return nil, xerror.WrapRecordNotExist()
	}

	return newOneData(subCluster), nil
}

func one(req *http.Request, subCluster *icluster_conf.SubClusterFilter) (*icluster_conf.SubCluster, error) {
	return container.SubClusterManager.FetchSubCluster(req.Context(), subCluster)
}

func newOneData(sc *icluster_conf.SubCluster) *OneData {
	if sc == nil {
		return nil
	}
	tmp := &OneData{
		Name:        sc.Name,
		Description: sc.Description,
		Ready:       sc.Ready,
		ProductName: sc.ProductName,
	}

	if sc.InstancePool != nil {
		tmp.InstancePool = sc.InstancePool.Name
		tmp.Instances = sc.InstancePool.Instances
	}

	return tmp
}
