// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK UsersV2Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/controllers/users"
)

type UserClient interface {
	users.UsersV2Client
}

func NewUserClient(config config.Config) UserClient {
	return users.NewUsersV2(config)
}
