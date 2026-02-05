# API-Key

## 1 创建API-Key

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	创建API-Key | | 
| 端点 |	/products/{product_name}/api-keys | |
| method |	POST | - |
| Content-Type | application/json | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名称 | Y | |

#### BODY 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| name | string | api-key名称。  | Y | 产品线内api-key名称不能重复。name参数不允许更新。 |
| enable | bool | 是否启用。  | Y | false：不启用；true：启用。 |
| key | string | api-key具体字符串 | Y | api-key格式为：产品线名称+多个随机生成的码段。允许的字符为大小写字母、数字以及-。 |
| is_limit | bool | 是否限额 | Y | false：不限额；true：有限额。|
| total_quota | int | 具体限额 | N | is_limit为true时，total_quota必填。取值范围：0-100000000。单位为个。|
| expired_time | string | 过期时间 | N | 空字符串: 永不过期；时间字符串:2025-01-01 01:01:01。时区以服务器时间为准。|
| allowed_models | []string | 允许的模型 | N | 不填代表不限制允许模型。 |
| allowed_subnets | []string | 允许的网段 | N | 不填和空数组代表不限制网段。 |

##### 请求示例
```shell
curl -X POST "http://api-gateway-server:port/open-api/v1/products/productname1/api-keys" -d data.json -H "Authorization:Token TOKEN_STRING" -H "Content-Type:application/json"
```

data.json如下：
```json
{
    "name": "test_key",
    "enable": true,
    "key": "productname1-7x9a2b4c8d3e5f",
    "is_limit": true,
    "total_quota": 50000,
    "expired_time": "2025-12-31 23:59:59",
    "allowed_models": ["gpt-3", "gpt-4"],
    "allowed_subnets": ["192.168.1.0/24"]
}
```

### 返回数据(Data内容)
无

#### 返回数据
状态码200为成功。

## 2 更新API-Key
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	更新API-Key || 
| 端点 |	/products/{product_name}/api-keys/{api_key_name} ||
| method |	PATCH | - |
| Content-Type | application/json | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名称 | Y | |
| api_key_name | string | API-Key名称|  Y | - |

#### Body参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| enable | bool | 是否启用。  | Y | false：不启用；true：启用。默认为false。 |
| is_limit | bool | 是否限额 | Y | false：不限额；true：有限额。|
| total_quota | int | 具体限额 | N | is_limit为true时，total_quota必填。取值范围：0-100000000。单位为个。|
| expired_time | string | 过期时间 | N | 空字符串: 永不过期；时间字符串:2025-01-01 01:01:01。时区以服务器时间为准。|
| allowed_models | []string | 允许的模型 | N | 不填代表不限制允许模型。 |
| allowed_subnets | []string | 允许的网段 | N | 不填和空数组代表不限制网段。 |

##### 请求示例
```shell
curl -X PATCH "http://api-server:port/open-api/v1/products/productname1/api-keys/test_key" -d data.json -H "Authorization:Token TOKEN_STRING" -H "Content-Type:application/json"
```
data.json如下：
```json
{
    "enable": true,
    "is_limit": true,
    "total_quota": 50000,
    "expired_time": "2025-12-31 23:59:59",
    "allowed_models": ["gpt-3", "gpt-4"],
    "allowed_subnets": ["192.168.1.0/24"]
}
```

### 返回数据(Data内容)
无

#### 返回数据  
状态码200为成功。

## 3 删除API-Key
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	删除API-Key || 
| 端点 |	/products/{product_name}/api-keys/{api_key_name} ||
| method |	DELETE | - |
| Content-Type | application/x-www-form-urlencoded | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名称 | Y | |
| api_key_name | string | API-Key名称|  Y | - |

#### Body参数
无

##### 请求示例
```shell
curl -X DELETE "http://api-server:port/open-api/v1/products/productname1/api-keys/test_key" -H "Authorization:Token TOKEN_STRING" -H "Content-Type:application/x-www-form-urlencoded"
```

### 返回数据(Data内容)
无

#### 返回数据  
状态码200为成功。

## 4 读取API-Key列表
### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 含义 |	读取API-Key列表 || 
| 端点 |	/products/{product_name}/api-keys ||
| method |	GET | - |
| Content-Type | application/x-www-form-urlencoded | - |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名称 | Y | |

#### Body参数
无

##### 请求示例
```shell
curl -X DELETE "http://api-server:port/open-api/v1/products/productname1/api-keys" -H "Authorization:Token TOKEN_STRING" -H "Content-Type:application/x-www-form-urlencoded"
```

### 返回数据(Data内容)
返回数据为列表。字段同 创建API-Key BODY参数。

#### 返回数据  
状态码200为成功。