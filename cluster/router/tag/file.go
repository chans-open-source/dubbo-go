/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tag

import (
	"encoding/base64"
	"net/url"
	"strconv"
	"sync"
)

import (
	"github.com/RoaringBitmap/roaring"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/chans-open-source/dubbo-go/cluster/router"
	"github.com/chans-open-source/dubbo-go/common"
	"github.com/chans-open-source/dubbo-go/common/constant"
	"github.com/chans-open-source/dubbo-go/protocol"
)

// FileTagRouter Use for parse config file of Tag router
type FileTagRouter struct {
	parseOnce  sync.Once
	router     *tagRouter
	routerRule *RouterRule
	url        *common.URL
	// force      bool
}

// NewFileTagRouter Create file tag router instance with content (from config file)
// todo fix this router, now it is useless, tag router is nil
func NewFileTagRouter(content []byte) (*FileTagRouter, error) {
	fileRouter := &FileTagRouter{}
	rule, err := getRule(string(content))
	if err != nil {
		return nil, perrors.Errorf("yaml.Unmarshal() failed , error:%v", perrors.WithStack(err))
	}
	fileRouter.routerRule = rule
	notify := make(chan struct{})
	fileRouter.url = common.NewURLWithOptions(
		common.WithProtocol(constant.TAG_ROUTE_PROTOCOL),
		common.WithParams(url.Values{}),
		common.WithParamsValue(constant.RULE_KEY, base64.URLEncoding.EncodeToString(content)),
		common.WithParamsValue(constant.ForceUseTag, strconv.FormatBool(rule.Force)),
		common.WithParamsValue(constant.RouterPriority, strconv.Itoa(rule.Priority)),
		common.WithParamsValue(constant.ROUTER_KEY, constant.TAG_ROUTE_PROTOCOL))
	fileRouter.router, err = NewTagRouter(fileRouter.URL(), notify)
	return fileRouter, err
}

// URL Return URL in file tag router n
func (f *FileTagRouter) URL() *common.URL {
	return f.url
}

// Priority Return Priority in listenable router
func (f *FileTagRouter) Priority() int64 {
	return f.router.priority
}

func (f *FileTagRouter) Route(invokers *roaring.Bitmap, cache router.Cache, url *common.URL, invocation protocol.Invocation) *roaring.Bitmap {
	return invokers
}
