# AI大模型路由规则配置

## 1 全量更新AI大模型路由规则

### 基本信息
| 项目  | 值  | 说明 |
| - | - | - |
| 端点|	/products/{product_name}/ai-route-rules ||
| 动作|	PATCH | |
| 含义|	全量更新AI大模型路由规则 |  |
| Content-Type | application/json | - |

#### 输入参数

##### URI 参数
无。

#### Body参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | - | - | - | - |
| rules | []Rule | 路由规则列表 | N | 为空代表清空规则。[表1：rule数据结构](#rule_data_structure)。字段domain、path_filter、method、header_filters与model_filter至少必填一项。 |

<a id="rule_data_structure">表1：rule 数据结构</a>

| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | - | - | - | - |
| name | string | 路由规则名称 | Y | 需要以字母或数字开头，允许数字、大小写字母、下划线、中划线组合且长度大于1，小于128。不能和其他的超时配置模版名称重复，创建后不允许修改。 |
| basic | object | 基础信息 | Y | |
| basic.domain | string | 域名 | N | |
| basic.path_filter | object | 路径匹配 | N | |
| basic.path_filter.match_mode | string | 路径匹配方式 | N | prefix_match：前缀匹配；exact_match：精确匹配；suffix_match：后缀匹配。 |
| basic.path_filter.ignore_case | bool | 是否忽略大小写 | N | true：忽略；false：不忽略。默认值为false。 |
| basic.path_filter.path | string | 请求路径 | N |  |
| basic.method | string | 请求方式 | N | 取值：GET，POST，DELETE，PATCH，PUT，OPTIONS。 |
| basic.header_filters | []object | Header匹配方式 | N | 单个 Header 值: ≤ 8KB；所有 Headers 总和: ≤ 16KB。 |
| basic.header_filters[].key | string | Header key | N | 只能包含可打印的 ASCII 字符（0x21-0x7E）；不能包含空格、冒号(:)、括号等特殊字符；不区分大小写，但约定使用首字母大写的连字符形式（如 "Content-Type"）。 |
| basic.header_filters[].value | string | Header value | N |  |
| basic.header_filters[].match_mode | string | Header匹配方式 | N | prefix_match：前缀匹配；exact_match：精确匹配；suffix_match：后缀匹配。header_filters不为空时必传。 |
| basic.header_filters[].ignore_case | bool | 是否忽略大小写 | N | true：忽略；false：不忽略。默认false。 |
| basic.model_filter | object | 模型匹配方式 | N | |
| basic.model_filter.name | string | 模型匹配名称 | N | 无需模型匹配条件时，留空（不生成模型匹配条件原语）|
| basic.model_filter.pattern | string | 模型匹配模式 | N | 无需模型匹配条件时，留空（不生成模型匹配条件原语）|
| basic.model_filter.ignore_case | bool | 是否忽略大小写 | N | true：忽略；false：不忽略。默认值为false。 |
| basic.expect_action | object            | 期望的动作 |  Y    | 详见[表：expect_action对象说明](#expect_action) |

<a id="expect_action">表：expect_action对象说明</a>

| 参数名       | 类型   | 参数含义     | 必填 | 补充描述           |
| ------------ | ------ | ------------ | ---- | ------------------ |
| forward         | object | 期望转发动作参数 | N    |       详见[表：expect_forward对象说明](#expect_forward)            |

<a id="expect_forward">表：expect_forward对象说明</a>

| 参数名       | 类型   | 参数含义     | 必填 | 补充描述           |
| ------------ | ------ | ------------ | ---- | ------------------ |
| cluster_name         | string | 期望的目标集群 | Y    |                    |
| url         | string | 期望的URL | N    |       必须是合法的URL，需要包括scheme(如http://)。如果url不为空，需要验证经过额外动作（如果有）后的URL是否是期望的URL；如果为url为空，则不验证。             |

##### 请求示例
curl -X PATCH "http://{api_server}/open-api/v1/products/productname1/ai-route-rules/actions/run-cases" -d data.json -H "Authorization:Token token_string" -H 'Content-Type:application/json'

##### 输入参数示例
data.json如下：
```json
{
    "rules":
    [
        {
            "name": "api_route_rule_001",
            "basic": {
                "domain": "api.example.com",
                "path_filter": {
                    "match_mode": "prefix_match",
                    "ignore_case": true,
                    "path": "/a"
                },
                "method": "POST",
                "header_filters": [{
                    "key": "X-API-Key",
                    "value": "prod_",
                    "match_mode": "prefix_match",
                    "ignore_case": false
                },{
                    "key": "Content-Type",
                    "value": "application/json",
                    "match_mode": "exact_match",
                    "ignore_case": true
                }],
                "model_filter": {
                    "name": "{\"name\":\"deepseek\"}",
                    "pattern": "deepseek",
                    "ignore_case": true
                },
                "expect_action": {
                    "forward": {
                        "cluster_name": "backend-cluster-1",
                        "url": "http://products-service.internal.com/api/v2/products/123?version=2"
                    }
                }
            }
        }
    ]
}
    
```

#### 返回数据(Data内容)	

##### 返回数据示例

```json
{
    "rules":
    [
        {
            "name": "api_route_rule_001",
            "basic": {
                "domain": "api.example.com",
                "path_filter": {
                    "match_mode": "prefix_match",
                    "ignore_case": true,
                    "path": "/a"
                },
                "method": "POST",
                "header_filters": [{
                    "key": "X-API-Key",
                    "value": "prod_",
                    "match_mode": "prefix_match",
                    "ignore_case": false
                },{
                    "key": "Content-Type",
                    "value": "application/json",
                    "match_mode": "exact_match",
                    "ignore_case": true
                }],
                "model_filter": {
                    "name": "{\"name\":\"deepseek\"}",
                    "pattern": "deepseek",
                    "ignore_case": true
                },
                "expect_action": {
                    "forward": {
                        "cluster_name": "backend-cluster-1",
                        "url": "http://products-service.internal.com/api/v2/products/123?version=2"
                    }
                }
            }
        }
    ]
}
```

##### 错误返回
| **错误码** | 错误信息 |
| ---------------------- | -------- |
| 422 | 参数不合法|
| 511 | 数据库异常|

## 2 获取AI大模型路由规则列表

### 基本信息
| 项目  | 值  | 说明 |
| - | - | - |
| 含义 |	获取AI大模型路由规则列表 ||
| 端点 |	/products/{product_name}/ai-route-rules  ||
| 动作 |	GET | - |
| Content-Type | application/x-www-form-urlencoded | - |

### 输入参数

#### Body参数
无。

#### 请求示例
curl 'http://api-server:port/open-api/v1/products/productname1/ai-route-rules' -H 'Authorization:Token token_string'  -H 'Content-Type:application/x-www-form-urlencoded'

### 返回数据(Data内容)	
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | - | - | - | - |
| rules | []Rule | 路由规则列表 | Y | 为空代表清空规则。 |

#### 返回数据示例
同 全量更新AI大模型路由规则 接口返回数据示例。

##### 错误返回
| **错误码** | 错误信息 |
| ---------------------- | -------- |
| 422 | 参数不合法|
| 511 | 数据库异常|

