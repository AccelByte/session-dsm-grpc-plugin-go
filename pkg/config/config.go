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
	// awsgamelift Config
	GRPCPort                    int    `env:"GRPC_PORT" envDocs:"The Port gRPC listens to" envDefault:"6565"`
	GameliftRegion              string `env:"GAMELIFT_REGION" envDocs:"Region where the gamelift is" envDefault:""`
	PluginGRPCServerAuthEnabled bool   `env:"PLUGIN_GRPC_SERVER_AUTH_ENABLED" envDocs:"Enable or disable access token and permission verification" envDefault:""`
	AWSAccessKeyId              string `env:"AWS_ACCESS_KEY_ID" envDocs:"AWS Access Key if using gamelift" envDefault:""`
	AWSSecretKeyId              string `env:"AWS_SECRET_ACCESS_KEY" envDocs:"AWS Secret Key if using gamelift" envDefault:""`
	// GCP Config
	GCPProjectID   string `env:"GCP_PROJECT_ID" envDocs:"Project ID in GCP VM" envDefault:""`
	GCPNetwork     string `env:"GCP_NETWORK" envDocs:"GCP Network for allow traffic and port" envDefault:""`
	GCPMachineType string `env:"GCP_MACHINE_TYPE" envDocs:"GCP Machine Type example e2-micro" envDefault:""`
	GCPRepository  string `env:"GCP_REPOSITORY" envDocs:"GCP Repository URL" envDefault:""`
	GCPRetry       int    `env:"GCP_RETRY" envDocs:"GCP Retry for get instance" envDefault:"3"`
	GCPWaitGetIP   int    `env:"GCP_WAIT_GET_IP" envDocs:"GCP Wait Get IP in seconds" envDefault:"1"`
	// AB Config
	ABBaseURL      string `env:"AB_BASE_URL" envDocs:"Base URL of AccelByte Gaming Services" envDefault:""`
	ABClientId     string `env:"AB_CLIENT_ID" envDocs:"Client ID from the Prerequisites section" envDefault:""`
	ABClientSecret string `env:"AB_CLIENT_SECRET" envDocs:"Client Secret from the Prerequisites section" envDefault:""`
	// Switching ds provider
	DsProvider string `env:"DS_PROVIDER" envDocs:"Ds Provider to choose (DEMO, GAMELIFT, GCP)" envDefault:"DEMO"`
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
