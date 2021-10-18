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

package getty

import (
	"bytes"
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

import (
	hessian "github.com/chans-open-source/dubbo-go-hessian2"

	perrors "github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/chans-open-source/dubbo-go/common"
	. "github.com/chans-open-source/dubbo-go/common/constant"
	"github.com/chans-open-source/dubbo-go/protocol"
	"github.com/chans-open-source/dubbo-go/protocol/invocation"
	"github.com/chans-open-source/dubbo-go/remoting"
)

func TestRunSuite(t *testing.T) {
	initTestEnvironment(t)
	userUrl := initUserUrl(t)
	server := getServer(userUrl)
	client := getClient(userUrl)
	assert.NotNil(t, client)

	testRequestOneWay(t, client)
	testClientCall(t, client)
	testClientAsyncCall(t, client)
	server.Stop()
}

//////////////////////////////////
// init special url
//////////////////////////////////

func initUserUrl(t *testing.T) *common.URL {
	hessian.RegisterPOJO(&User{})
	remoting.RegistryCodec("dubbo", &DubboTestCodec{})

	methods, err := common.ServiceMap.Register("com.ikurento.user.UserProvider", "dubbo", "", "", &UserProvider{})
	assert.NoError(t, err)
	assert.Equal(t, "GetBigPkg,GetUser,GetUser0,GetUser1,GetUser2,GetUser3,GetUser4,GetUser5,GetUser6", methods)

	url, err := common.NewURL("dubbo://127.0.0.1:20060/com.ikurento.user.UserProvider?anyhost=true&" +
		"application=BDTService&category=providers&default.timeout=10000&dubbo=dubbo-provider-golang-1.0.0&" +
		"environment=dev&interface=com.ikurento.user.UserProvider&ip=127.0.0.1&methods=GetUser%2C&" +
		"module=dubbogo+user-info+server&org=ikurento.com&owner=ZX&pid=1447&revision=0.0.1&" +
		"side=provider&timeout=3000&timestamp=1556509797245&bean.name=UserProvider")
	assert.NoError(t, err)

	return url
}

//////////////////////////////////
// test cases
//////////////////////////////////

func testRequestOneWay(t *testing.T, client *Client) {

	request := remoting.NewRequest("2.0.2")
	invocation := createInvocation("GetUser", nil, nil, []interface{}{"1", "username"},
		[]reflect.Value{reflect.ValueOf("1"), reflect.ValueOf("username")})
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = false
	err := client.Request(request, 3*time.Second, nil)
	assert.NoError(t, err)
}

func testClientCall(t *testing.T, c *Client) {
	testGetBigPkg(t, c)
	testGetUser(t, c)
	testGetUser0(t, c)
	testGetUser1(t, c)
	testGetUser2(t, c)
	testGetUser3(t, c)
	testGetUser4(t, c)
	testGetUser5(t, c)
	testGetUser6(t, c)
	testGetUser61(t, c)
}

func testGetBigPkg(t *testing.T, c *Client) {
	user := &User{}
	request := remoting.NewRequest("2.0.2")
	invocation := createInvocation("GetBigPkg", nil, nil, []interface{}{[]interface{}{nil}, user},
		[]reflect.Value{reflect.ValueOf([]interface{}{nil}), reflect.ValueOf(user)})
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	pendingResponse := remoting.NewPendingResponse(request.ID)
	pendingResponse.Reply = user
	remoting.AddPendingResponse(pendingResponse)
	err := c.Request(request, 8*time.Second, pendingResponse)
	assert.NoError(t, err)
	assert.NotEqual(t, "", user.Id)
	assert.NotEqual(t, "", user.Name)
}

func testGetUser(t *testing.T, c *Client) {
	user := &User{}
	request := remoting.NewRequest("2.0.2")
	invocation := createInvocation("GetUser", nil, nil, []interface{}{"1", "username"},
		[]reflect.Value{reflect.ValueOf("1"), reflect.ValueOf("username")})
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	pendingResponse := remoting.NewPendingResponse(request.ID)
	pendingResponse.Reply = user
	remoting.AddPendingResponse(pendingResponse)
	err := c.Request(request, 3*time.Second, pendingResponse)
	assert.NoError(t, err)
	assert.Equal(t, User{Id: "1", Name: "username"}, *user)
}

func testGetUser0(t *testing.T, c *Client) {
	var (
		user *User
		err  error
	)
	user = &User{}
	request := remoting.NewRequest("2.0.2")
	invocation := createInvocation("GetUser0", nil, nil, []interface{}{"1", nil, "username"},
		[]reflect.Value{reflect.ValueOf("1"), reflect.ValueOf(nil), reflect.ValueOf("username")})
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	rsp := remoting.NewPendingResponse(request.ID)
	rsp.SetResponse(remoting.NewResponse(request.ID, "2.0.2"))
	remoting.AddPendingResponse(rsp)
	rsp.Reply = user
	err = c.Request(request, 3*time.Second, rsp)
	assert.NoError(t, err)
	assert.Equal(t, User{Id: "1", Name: "username"}, *user)
}

func testGetUser1(t *testing.T, c *Client) {
	var (
		err error
	)
	request := remoting.NewRequest("2.0.2")
	invocation := createInvocation("GetUser1", nil, nil, []interface{}{},
		[]reflect.Value{})
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	pendingResponse := remoting.NewPendingResponse(request.ID)
	user := &User{}
	pendingResponse.Reply = user
	remoting.AddPendingResponse(pendingResponse)
	err = c.Request(request, 3*time.Second, pendingResponse)
	assert.NoError(t, err)
}

func testGetUser2(t *testing.T, c *Client) {
	var (
		err error
	)
	request := remoting.NewRequest("2.0.2")
	invocation := createInvocation("GetUser2", nil, nil, []interface{}{},
		[]reflect.Value{})
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	pendingResponse := remoting.NewPendingResponse(request.ID)
	remoting.AddPendingResponse(pendingResponse)
	err = c.Request(request, 3*time.Second, pendingResponse)
	assert.EqualError(t, err, "error")
}

func testGetUser3(t *testing.T, c *Client) {
	var (
		err error
	)
	request := remoting.NewRequest("2.0.2")
	invocation := createInvocation("GetUser3", nil, nil, []interface{}{},
		[]reflect.Value{})
	attachment := map[string]string{
		INTERFACE_KEY: "com.ikurento.user.UserProvider",
	}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	pendingResponse := remoting.NewPendingResponse(request.ID)
	user2 := []interface{}{}
	pendingResponse.Reply = &user2
	remoting.AddPendingResponse(pendingResponse)
	err = c.Request(request, 3*time.Second, pendingResponse)
	assert.NoError(t, err)
	assert.Equal(t, &User{Id: "1", Name: "username"}, user2[0])
}

func testGetUser4(t *testing.T, c *Client) {
	var (
		err error
	)
	request := remoting.NewRequest("2.0.2")
	invocation := invocation.NewRPCInvocation("GetUser4", []interface{}{[]interface{}{"1", "username"}}, nil)
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	pendingResponse := remoting.NewPendingResponse(request.ID)
	user2 := []interface{}{}
	pendingResponse.Reply = &user2
	remoting.AddPendingResponse(pendingResponse)
	err = c.Request(request, 3*time.Second, pendingResponse)
	assert.NoError(t, err)
	assert.Equal(t, &User{Id: "1", Name: "username"}, user2[0])
}

func testGetUser5(t *testing.T, c *Client) {
	var (
		err error
	)
	request := remoting.NewRequest("2.0.2")
	invocation := invocation.NewRPCInvocation("GetUser5", []interface{}{map[interface{}]interface{}{"id": "1", "name": "username"}}, nil)
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	pendingResponse := remoting.NewPendingResponse(request.ID)
	user3 := map[interface{}]interface{}{}
	pendingResponse.Reply = &user3
	remoting.AddPendingResponse(pendingResponse)
	err = c.Request(request, 3*time.Second, pendingResponse)
	assert.NoError(t, err)
	assert.NotNil(t, user3)
	assert.Equal(t, &User{Id: "1", Name: "username"}, user3["key"])
}

func testGetUser6(t *testing.T, c *Client) {
	var (
		user *User
		err  error
	)
	user = &User{}
	request := remoting.NewRequest("2.0.2")
	invocation := invocation.NewRPCInvocation("GetUser6", []interface{}{0}, nil)
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	pendingResponse := remoting.NewPendingResponse(request.ID)
	pendingResponse.Reply = user
	remoting.AddPendingResponse(pendingResponse)
	err = c.Request(request, 3*time.Second, pendingResponse)
	assert.NoError(t, err)
	assert.Equal(t, User{Id: "", Name: ""}, *user)
}

func testGetUser61(t *testing.T, c *Client) {
	var (
		user *User
		err  error
	)
	user = &User{}
	request := remoting.NewRequest("2.0.2")
	invocation := invocation.NewRPCInvocation("GetUser6", []interface{}{1}, nil)
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	pendingResponse := remoting.NewPendingResponse(request.ID)
	pendingResponse.Reply = user
	remoting.AddPendingResponse(pendingResponse)
	err = c.Request(request, 3*time.Second, pendingResponse)
	assert.NoError(t, err)
	assert.Equal(t, User{Id: "1", Name: ""}, *user)
}

func testClientAsyncCall(t *testing.T, client *Client) {
	user := &User{}
	wg := sync.WaitGroup{}
	request := remoting.NewRequest("2.0.2")
	invocation := createInvocation("GetUser0", nil, nil, []interface{}{"4", nil, "username"},
		[]reflect.Value{reflect.ValueOf("4"), reflect.ValueOf(nil), reflect.ValueOf("username")})
	attachment := map[string]string{INTERFACE_KEY: "com.ikurento.user.UserProvider"}
	setAttachment(invocation, attachment)
	request.Data = invocation
	request.Event = false
	request.TwoWay = true
	rsp := remoting.NewPendingResponse(request.ID)
	rsp.SetResponse(remoting.NewResponse(request.ID, "2.0.2"))
	remoting.AddPendingResponse(rsp)
	rsp.Reply = user
	rsp.Callback = func(response common.CallbackResponse) {
		r := response.(remoting.AsyncCallbackResponse)
		rst := *r.Reply.(*remoting.Response).Result.(*protocol.RPCResult)
		assert.Equal(t, User{Id: "4", Name: "username"}, *(rst.Rest.(*User)))
		wg.Done()
	}
	wg.Add(1)
	err := client.Request(request, 3*time.Second, rsp)
	assert.NoError(t, err)
	assert.Equal(t, User{}, *user)
	wg.Wait()
}

//////////////////////////////////
// provider
//////////////////////////////////

type (
	User struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	UserProvider struct {
		//user map[string]User
	}
)

// size:4801228
func (u *UserProvider) GetBigPkg(ctx context.Context, req []interface{}, rsp *User) error {
	argBuf := new(bytes.Buffer)
	for i := 0; i < 400; i++ {
		argBuf.WriteString("击鼓其镗，踊跃用兵。土国城漕，我独南行。从孙子仲，平陈与宋。不我以归，忧心有忡。爰居爰处？爰丧其马？于以求之？于林之下。死生契阔，与子成说。执子之手，与子偕老。于嗟阔兮，不我活兮。于嗟洵兮，不我信兮。")
		argBuf.WriteString("击鼓其镗，踊跃用兵。土国城漕，我独南行。从孙子仲，平陈与宋。不我以归，忧心有忡。爰居爰处？爰丧其马？于以求之？于林之下。死生契阔，与子成说。执子之手，与子偕老。于嗟阔兮，不我活兮。于嗟洵兮，不我信兮。")
	}
	rsp.Id = argBuf.String()
	rsp.Name = argBuf.String()
	return nil
}

func (u *UserProvider) GetUser(ctx context.Context, req []interface{}, rsp *User) error {
	rsp.Id = req[0].(string)
	rsp.Name = req[1].(string)
	return nil
}

func (u *UserProvider) GetUser0(id string, k *User, name string) (User, error) {
	// fix testClientAsyncCall assertion
	time.Sleep(1 * time.Second)
	return User{Id: id, Name: name}, nil
}

func (u *UserProvider) GetUser1() error {
	return nil
}

func (u *UserProvider) GetUser2() error {
	return perrors.New("error")
}

func (u *UserProvider) GetUser3(rsp *[]interface{}) error {
	*rsp = append(*rsp, User{Id: "1", Name: "username"})
	return nil
}

func (u *UserProvider) GetUser4(ctx context.Context, req []interface{}) ([]interface{}, error) {

	return []interface{}{User{Id: req[0].([]interface{})[0].(string), Name: req[0].([]interface{})[1].(string)}}, nil
}

func (u *UserProvider) GetUser5(ctx context.Context, req []interface{}) (map[interface{}]interface{}, error) {
	return map[interface{}]interface{}{"key": User{Id: req[0].(map[interface{}]interface{})["id"].(string), Name: req[0].(map[interface{}]interface{})["name"].(string)}}, nil
}

func (u *UserProvider) GetUser6(id int64) (*User, error) {
	if id == 0 {
		return nil, nil
	}
	return &User{Id: "1"}, nil
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}

func (u User) JavaClassName() string {
	return "com.ikurento.user.User"
}
