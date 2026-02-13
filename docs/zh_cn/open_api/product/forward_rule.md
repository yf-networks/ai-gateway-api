
# 转发规则 

## 1 转发规则的模型
BFE的转发规则分为两张表: 
- AI大模型路由规则
- 默认转发规则

转发规则的匹配顺序:
- 先匹配AI大模型路由规则
- 匹配失败使用默认转发规则

说明:
- 转发规则中使用到的集群，必须是已经就绪的状态。否则 API 服务器拒绝本次提交

### AI大模型路由规则
AI大模型路由规则
- AI大模型路由规则表是一个有序的列表，规则的先后顺序，与转发引擎执行的顺序相同
    - 规则匹配的时间复杂度是O(N)
- AI大模型路由规则表中每一条规则，都是用Condition元语编写的条件表达式，以及集群
- 默认转发规则，必须是:default_t()，它表示缺省的转发规则
   - 当其他所有规则都不匹配时，无条件匹配default_t()的规则 

| 条件 | 目标集群 | 说明 | 
| - | - | - |
| req_host_in("www.xyz.com") &&req_path_in(“/path1”) &&req_cookie_value_in("key1", "value1", false) | Cluster1 | 请求域名是 www.xyz.com，请求 path 是 /path1 ， 并 且 带 有 cookie: key1=value1，则转发到 Cluster1 |
 | req_host_in("www.xyz.com") &&req_path_in(“/path1”) | Cluster2 | 请求域名是 www.xyz.com，请求 path 是 /path1，则转发到 Cluster2 |
 | default_t() | Cluster3 | 其他情况下，一律进入 Cluster3| 


### 组合使用
- AI大模型路由规则表可以非常自由的指定流量特征，实现更复杂的转发
- AI大模型路由规则表中都不匹配时，使用默认转发规则


典型场景:TODO

## 2 更新转发规则

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点| /products/{product_name}/routes ||
| method | PATCH ||
| 含义 | 整体更新转发列表 |- |

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | - |

#### Body参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| default_forward_rule | object | 默认转发规则 | Y |  |
| default_forward_rule.action | object | 动作 | Y |  |
| default_forward_rule.action.forward | object | 转发动作 | Y |  |
| default_forward_rule.action.forward.cluster_name | string | 转发的集群名称 | Y |  |

#### Body 请求示例
```json
{
    "default_forward_rule":{
        "action": {
            "forward":{
                "cluster_name": "clustername"
            }
        }
    }
}
```

### 返回数据(Data内容)	

#### 返回数据示例
同 请求数据
 
## 3 获取转发规则列表

### 基本信息
| 项目  | 值  | 说明 | 
| - | - | - |
| 端点  | /products/{product_name}/routes | |
| method | GET | |
| 含义 | 获取产品线的转发规则列表 ||

### 输入参数

#### URI 参数
| 参数名 | 类型 |参数含义 | 必填 | 补充描述 |
| - | -  | - | - | - | 
| product_name | string | 产品线名字 | Y | - |
  
### 返回数据(Data内容)	
同 更新转发规则