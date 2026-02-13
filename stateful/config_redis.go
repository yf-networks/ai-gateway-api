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

package stateful

import (
	"fmt"

	"github.com/bfenetworks/bfe/bfe_util/bns"
	"github.com/bfenetworks/bfe/bfe_util/redis_client"
)

// redis conf
type Redis struct {
	Bns            string // bns name for redis proxy
	ConnectTimeout int    // connect timeout (ms)
	ReadTimeout    int    // read timeout (ms)
	WriteTimeout   int    // write timeout(ms)

	// max idle connections in pool
	MaxIdle int

	// redis password，ignore if not set
	Password string

	// max active connections in pool,
	// when set 0, there is no connection num limit
	MaxActive int

	// redis cluster mode, value in {"official", "proxy"}.
	// if "official", use opensource redis cluster
	ClusterMode string
}

func (r *Redis) Init() error {
	// new Redis Client
	AccessLogger.Info("redis conf:%v", *r)

	options := &redis_client.Options{
		ServiceConf:    r.Bns,
		MaxIdle:        r.MaxIdle,
		MaxActive:      r.MaxActive,
		Wait:           false,
		ConnTimeoutMs:  r.ConnectTimeout,
		ReadTimeoutMs:  r.ReadTimeout,
		WriteTimeoutMs: r.WriteTimeout,
		Password:       r.Password,
	}

	err := bns.LoadLocalNameConf("./conf/name_conf.data")
	if err != nil {
		AccessLogger.Error("load redis conf ./conf/name_conf.data error:%s", err.Error())
		return err
	}

	client := redis_client.NewRedisClient(options)
	if DefaultClientSet == nil {
		DefaultClientSet = new(ClientSet)
	}

	DefaultClientSet.RedisClient = client

	return nil
}

func AIUsedQuotaKey(key string, updatetime int64) string {
	return fmt.Sprintf("usedquota_%s:%d", key, updatetime)
}
