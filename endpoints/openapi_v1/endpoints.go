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

package openapi_v1

import (
	"github.com/gorilla/mux"

	"github.com/yf-networks/ai-gateway-api/endpoints/middleware"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/ai_route"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/api_key"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/auth"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/bfe_cluster"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/bfe_pool"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/certificate"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/domain"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/general"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/product"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/product_cluster"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/product_pool"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/route"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/subcluster"
	"github.com/yf-networks/ai-gateway-api/endpoints/openapi_v1/traffic"
	"github.com/yf-networks/ai-gateway-api/lib/xreq"
)

func RegisterEndpoints(router *mux.Router) *mux.Router {
	openAPIV1Router := router.PathPrefix("/open-api/v1").Subrouter()
	openAPIV1Router.Use(middleware.McProductProbe, middleware.McUserProbe)
	for _, one := range endpoints() {
		one.Register(openAPIV1Router)
	}
	return openAPIV1Router
}

func endpoints() []*xreq.Endpoint {
	return merge(
		product.Routers,
		product_cluster.Endpoints,
		certificate.Endpoints,
		product_pool.Endpoints,
		subcluster.Endpoints,
		bfe_pool.Endpoints,
		auth.Endpoints,
		traffic.Endpoints,
		bfe_cluster.Endpoints,
		route.Endpoints,
		domain.Endpoints,
		api_key.Endpoints,
		ai_route.Endpoints,
		general.Endpoints,
	)
}

func merge(rss ...[]*xreq.Endpoint) (rs []*xreq.Endpoint) {
	for _, r := range rss {
		rs = append(rs, r...)
	}
	return
}
