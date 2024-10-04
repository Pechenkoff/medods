// Code generated by mockery v2.46.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Producer is an autogenerated mock type for the Producer type
type Producer struct {
	mock.Mock
}

type Producer_Expecter struct {
	mock *mock.Mock
}

func (_m *Producer) EXPECT() *Producer_Expecter {
	return &Producer_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *Producer) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Producer_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type Producer_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *Producer_Expecter) Close() *Producer_Close_Call {
	return &Producer_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *Producer_Close_Call) Run(run func()) *Producer_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Producer_Close_Call) Return(_a0 error) *Producer_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Producer_Close_Call) RunAndReturn(run func() error) *Producer_Close_Call {
	_c.Call.Return(run)
	return _c
}

// SendMessage provides a mock function with given fields: topic, message
func (_m *Producer) SendMessage(topic string, message []byte) error {
	ret := _m.Called(topic, message)

	if len(ret) == 0 {
		panic("no return value specified for SendMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte) error); ok {
		r0 = rf(topic, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Producer_SendMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMessage'
type Producer_SendMessage_Call struct {
	*mock.Call
}

// SendMessage is a helper method to define mock.On call
//   - topic string
//   - message []byte
func (_e *Producer_Expecter) SendMessage(topic interface{}, message interface{}) *Producer_SendMessage_Call {
	return &Producer_SendMessage_Call{Call: _e.mock.On("SendMessage", topic, message)}
}

func (_c *Producer_SendMessage_Call) Run(run func(topic string, message []byte)) *Producer_SendMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]byte))
	})
	return _c
}

func (_c *Producer_SendMessage_Call) Return(_a0 error) *Producer_SendMessage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Producer_SendMessage_Call) RunAndReturn(run func(string, []byte) error) *Producer_SendMessage_Call {
	_c.Call.Return(run)
	return _c
}

// NewProducer creates a new instance of Producer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProducer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Producer {
	mock := &Producer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
