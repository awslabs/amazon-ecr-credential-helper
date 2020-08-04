// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api (interfaces: ClientFactory,Client)

package mock_api

import (
	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	api "github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	gomock "github.com/golang/mock/gomock"
)

// Mock of ClientFactory interface
type MockClientFactory struct {
	ctrl     *gomock.Controller
	recorder *_MockClientFactoryRecorder
}

// Recorder for MockClientFactory (not exported)
type _MockClientFactoryRecorder struct {
	mock *MockClientFactory
}

func NewMockClientFactory(ctrl *gomock.Controller) *MockClientFactory {
	mock := &MockClientFactory{ctrl: ctrl}
	mock.recorder = &_MockClientFactoryRecorder{mock}
	return mock
}

func (_m *MockClientFactory) EXPECT() *_MockClientFactoryRecorder {
	return _m.recorder
}

func (_m *MockClientFactory) NewClient(_param0 *session.Session, _param1 *aws.Config) api.Client {
	ret := _m.ctrl.Call(_m, "NewClient", _param0, _param1)
	ret0, _ := ret[0].(api.Client)
	return ret0
}

func (_mr *_MockClientFactoryRecorder) NewClient(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewClient", arg0, arg1)
}

func (_m *MockClientFactory) NewClientFromRegion(_param0 string) api.Client {
	ret := _m.ctrl.Call(_m, "NewClientFromRegion", _param0)
	ret0, _ := ret[0].(api.Client)
	return ret0
}

func (_mr *_MockClientFactoryRecorder) NewClientFromRegion(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewClientFromRegion", arg0)
}

func (_m *MockClientFactory) NewClientWithDefaults() api.Client {
	ret := _m.ctrl.Call(_m, "NewClientWithDefaults")
	ret0, _ := ret[0].(api.Client)
	return ret0
}

func (_mr *_MockClientFactoryRecorder) NewClientWithDefaults() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewClientWithDefaults")
}

func (_m *MockClientFactory) NewClientWithFipsEndpoint(_param0 string) (api.Client, error) {
	ret := _m.ctrl.Call(_m, "NewClientWithFipsEndpoint", _param0)
	ret0, _ := ret[0].(api.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientFactoryRecorder) NewClientWithFipsEndpoint(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewClientWithFipsEndpoint", arg0)
}

func (_m *MockClientFactory) NewClientWithExplicitProfile(_param0 string) (api.Client, error) {
	ret := _m.ctrl.Call(_m, "NewClientWithExplicitProfile", _param0)
	ret0, _ := ret[0].(api.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientFactoryRecorder) NewClientWithExplicitProfile(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewClientWithExplicitProfile", arg0)
}

func (_m *MockClientFactory) NewClientWithOptions(_param0 api.Options) api.Client {
	ret := _m.ctrl.Call(_m, "NewClientWithOptions", _param0)
	ret0, _ := ret[0].(api.Client)
	return ret0
}

func (_mr *_MockClientFactoryRecorder) NewClientWithOptions(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NewClientWithOptions", arg0)
}

// Mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *_MockClientRecorder
}

// Recorder for MockClient (not exported)
type _MockClientRecorder struct {
	mock *MockClient
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &_MockClientRecorder{mock}
	return mock
}

func (_m *MockClient) EXPECT() *_MockClientRecorder {
	return _m.recorder
}

func (_m *MockClient) GetCredentials(_param0 string) (*api.Auth, error) {
	ret := _m.ctrl.Call(_m, "GetCredentials", _param0)
	ret0, _ := ret[0].(*api.Auth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetCredentials(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetCredentials", arg0)
}

func (_m *MockClient) GetCredentialsByRegistryID(_param0 string) (*api.Auth, error) {
	ret := _m.ctrl.Call(_m, "GetCredentialsByRegistryID", _param0)
	ret0, _ := ret[0].(*api.Auth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) GetCredentialsByRegistryID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetCredentialsByRegistryID", arg0)
}

func (_m *MockClient) ListCredentials() ([]*api.Auth, error) {
	ret := _m.ctrl.Call(_m, "ListCredentials")
	ret0, _ := ret[0].([]*api.Auth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) ListCredentials() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ListCredentials")
}
