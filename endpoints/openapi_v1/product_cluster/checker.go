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

package product_cluster

import (
	"fmt"
	"strconv"

	"github.com/yf-networks/ai-gateway-api/lib/xerror"
	"github.com/yf-networks/ai-gateway-api/model/icluster_conf"
)

const (
	maxServiceNameLen = 255
	maxGroupLen       = 255
)

func checkLLMConfig(llmConfig *icluster_conf.LLMConfig) error {
	if llmConfig == nil {
		return nil
	}

	if llmConfig.ServiceName != nil && len(*llmConfig.ServiceName) > maxServiceNameLen {
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("llm_config.service_name length must be lower than %s", strconv.Itoa(maxServiceNameLen)))
	}

	if llmConfig.Group != nil && len(*llmConfig.Group) > maxGroupLen {
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("llm_config.group length must be lower than %s", strconv.Itoa(maxGroupLen)))
	}

	if llmConfig.ModelEndpoint != nil {
		switch llmConfig.ModelEndpoint.Schema {
		case "http":
		case "https":
		default:
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("llm_config.model_endpoint.schema must be http or https"))
		}
	}

	if llmConfig.Enable == nil {
		return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Must set llm_config.enable"))
	}

	if llmConfig.Enable != nil && *llmConfig.Enable {
		if llmConfig.ServiceName == nil || *llmConfig.ServiceName == "" {
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Must set llm_config.service_name"))
		}

		if llmConfig.ModelEndpoint == nil {
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Must set llm_config.model_endpoint"))
		}

		if llmConfig.ModelEndpoint.URI == "" {
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Must set llm_config.model_endpoint.uri"))
		}

		if len(llmConfig.Models) == 0 {
			return xerror.WrapParamErrorWithMsg(fmt.Sprintf("Must set llm_config.models"))
		}
	}

	return nil
}
