// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package config

import (
	"reflect"
	"testing"
)

func TestConfig_EnvironmentVariables(t *testing.T) {
	tests := []struct {
		name                         string
		nativeSessionDeletionEnabled bool
		queueSize                    int
		exposedVariables             map[string]bool
		want                         []EnvironmentVariable
	}{
		{
			name:                         "success",
			nativeSessionDeletionEnabled: false,
			queueSize:                    90,
			exposedVariables: map[string]bool{
				"NATIVE_SYNC_DELETION_ENABLED": true,
				"QUEUE_SIZE":                   true,
			},
			want: []EnvironmentVariable{
				{
					Name:         "NATIVE_SYNC_DELETION_ENABLED",
					Description:  "Set to false to disable deleting native session from 3rd party",
					DefaultValue: "true",
					ActualValue:  "false",
				},
				{
					Name:         "QUEUE_SIZE",
					Description:  "number of messages fits in the queue",
					DefaultValue: "100",
					ActualValue:  "90",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				NativeSessionDeletionEnabled: tt.nativeSessionDeletionEnabled,
				QueueSize:                    tt.queueSize,
			}
			if got := config.EnvironmentVariables(tt.exposedVariables); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.EnvironmentVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}
