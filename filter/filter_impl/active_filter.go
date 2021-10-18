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

package filter_impl

import (
	"context"
	"strconv"
)

import (
	"github.com/chans-open-source/dubbo-go/common/extension"
	"github.com/chans-open-source/dubbo-go/common/logger"
	"github.com/chans-open-source/dubbo-go/filter"
	"github.com/chans-open-source/dubbo-go/protocol"
	invocation2 "github.com/chans-open-source/dubbo-go/protocol/invocation"
)

const (
	active               = "active"
	dubboInvokeStartTime = "dubboInvokeStartTime"
)

func init() {
	extension.SetFilter(active, GetActiveFilter)
}

// ActiveFilter tracks the requests status
type ActiveFilter struct {
}

// Invoke starts to record the requests status
func (ef *ActiveFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	logger.Infof("invoking active filter. %v,%v", invocation.MethodName(), len(invocation.Arguments()))
	invocation.(*invocation2.RPCInvocation).SetAttachments(dubboInvokeStartTime, strconv.FormatInt(protocol.CurrentTimeMillis(), 10))
	protocol.BeginCount(invoker.GetURL(), invocation.MethodName())
	return invoker.Invoke(ctx, invocation)
}

// OnResponse update the active count base on the request result.
func (ef *ActiveFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	startTime, err := strconv.ParseInt(invocation.(*invocation2.RPCInvocation).AttachmentsByKey(dubboInvokeStartTime, "0"), 10, 64)
	if err != nil {
		result.SetError(err)
		logger.Errorf("parse dubbo_invoke_start_time to int64 failed")
		return result
	}
	elapsed := protocol.CurrentTimeMillis() - startTime
	protocol.EndCount(invoker.GetURL(), invocation.MethodName(), elapsed, result.Error() == nil)
	return result
}

// GetActiveFilter creates ActiveFilter instance
func GetActiveFilter() filter.Filter {
	return &ActiveFilter{}
}
