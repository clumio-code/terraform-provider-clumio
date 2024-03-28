// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK TasksV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkTasks "github.com/clumio-code/clumio-go-sdk/controllers/tasks"
)

type TaskClient interface {
	sdkTasks.TasksV1Client
}

func NewTaskClient(config config.Config) TaskClient {
	return sdkTasks.NewTasksV1(config)
}
