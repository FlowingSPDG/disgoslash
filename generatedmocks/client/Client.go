// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	models "github.com/wafer-bw/disgoslash/models"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// CreateApplicationCommand provides a mock function with given fields: guildID, command
func (_m *Client) CreateApplicationCommand(guildID string, command *models.ApplicationCommand) error {
	ret := _m.Called(guildID, command)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *models.ApplicationCommand) error); ok {
		r0 = rf(guildID, command)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteApplicationCommand provides a mock function with given fields: guildID, commandID
func (_m *Client) DeleteApplicationCommand(guildID string, commandID string) error {
	ret := _m.Called(guildID, commandID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(guildID, commandID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListApplicationCommands provides a mock function with given fields: guildID
func (_m *Client) ListApplicationCommands(guildID string) ([]*models.ApplicationCommand, error) {
	ret := _m.Called(guildID)

	var r0 []*models.ApplicationCommand
	if rf, ok := ret.Get(0).(func(string) []*models.ApplicationCommand); ok {
		r0 = rf(guildID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.ApplicationCommand)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(guildID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
