// Copyright (c) 2018-2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package config

import (
	"fmt"
	"reflect"
)

// Config specifies configurable options through env vars
//
//nolint:lll
type Config struct {
	GRPCPort       int    `env:"GRPC_PORT" envDocs:"The Port gRPC listens to" envDefault:"9090"`
	GameliftRegion string `env:"GAMELIFT_REGION" envDocs:"Region where the gamelift is" envDefault:""`
}

// HelpDocs returns documentation of Config based on field tags.
func (envVar Config) HelpDocs() []string {
	environmentVariables := envVar.EnvironmentVariables(nil)
	doc := make([]string, 1+len(environmentVariables))
	doc[0] = "Environment variables config:"
	for i := 1; i <= len(environmentVariables); i++ {
		doc[i+1] = fmt.Sprintf("  %v\t %v (default: %v)", environmentVariables[i].Name, environmentVariables[i].Description, environmentVariables[i].DefaultValue)
	}

	return doc
}

// EnvironmentVariables method to get a list of environment variables.
func (envVar Config) EnvironmentVariables(exposedVariables map[string]bool) []EnvironmentVariable {
	environmentVariables := make([]EnvironmentVariable, 0)
	reflectValue := reflect.ValueOf(envVar)
	reflectType := reflectValue.Type()

	for i := 0; i < reflectValue.NumField(); i++ {
		environmentVariable := newEnvironmentVariable(reflectValue, reflectType, i)
		if exposedVariables != nil {
			if _, ok := exposedVariables[environmentVariable.Name]; !ok {
				continue
			}
		}

		environmentVariables = append(environmentVariables, environmentVariable)
	}

	return environmentVariables
}

// EnvironmentVariable struct which contains env tags in config field.
type EnvironmentVariable struct {
	Name         string
	Description  string
	DefaultValue string
	ActualValue  string
}

func newEnvironmentVariable(reflectValue reflect.Value, reflectType reflect.Type, index int) EnvironmentVariable {
	field := reflectType.Field(index)

	return EnvironmentVariable{
		Name:         field.Tag.Get("env"),
		Description:  field.Tag.Get("envDocs"),
		DefaultValue: field.Tag.Get("envDefault"),
		ActualValue:  fmt.Sprintf("%v", reflectValue.Field(index).Interface()),
	}
}
