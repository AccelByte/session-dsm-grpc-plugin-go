// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package gcpvm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"session-dsm-grpc-plugin/pkg/config"
	"session-dsm-grpc-plugin/pkg/constants"
	sessiondsm "session-dsm-grpc-plugin/pkg/pb"
	"session-dsm-grpc-plugin/pkg/session"
	sessionClient "session-dsm-grpc-plugin/pkg/session"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
	"strings"
	"time"

	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/iam"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

var (
	gcpRegionZones = map[string][]string{
		"us-east1":                {"us-east1-b", "us-east1-c", "us-east1-d"},
		"us-east4":                {"us-east4-a", "us-east4-b", "us-east4-c"},
		"us-west1":                {"us-west1-a", "us-west1-b", "us-west1-c"},
		"us-west2":                {"us-west2-a", "us-west2-b", "us-west2-c"},
		"northamerica-northeast1": {"northamerica-northeast1-a", "northamerica-northeast1-b", "northamerica-northeast1-c"},
		"southamerica-east1":      {"southamerica-east1-a", "southamerica-east1-b", "southamerica-east1-c"},
		"europe-west3":            {"europe-west3-a", "europe-west3-b", "europe-west3-c"},
		"europe-west1":            {"europe-west1-b", "europe-west1-c", "europe-west1-d"},
		"europe-west2":            {"europe-west2-a", "europe-west2-b", "europe-west2-c"},
		"europe-west9":            {"europe-west9-a", "europe-west9-b", "europe-west9-c"},
		"europe-north1":           {"europe-north1-a", "europe-north1-b", "europe-north1-c"},
		"me-west1":                {"me-west1-a", "me-west1-b", "me-west1-c"},
		"africa-north1":           {"africa-north1-a", "africa-north1-b", "africa-north1-c"},
		"asia-east2":              {"asia-east2-a", "asia-east2-b", "asia-east2-c"},
		"asia-south1":             {"asia-south1-a", "asia-south1-b", "asia-south1-c"},
		"asia-northeast2":         {"asia-northeast2-a", "asia-northeast2-b", "asia-northeast2-c"},
		"asia-northeast3":         {"asia-northeast3-a", "asia-northeast3-b", "asia-northeast3-c"},
		"asia-southeast1":         {"asia-southeast1-a", "asia-southeast1-b", "asia-southeast1-c"},
		"australia-southeast1":    {"australia-southeast1-a", "australia-southeast1-b", "australia-southeast1-c"},
		"asia-northeast1":         {"asia-northeast1-a", "asia-northeast1-b", "asia-northeast1-c"},
	}

	awsToGCPRegionMap = map[string]string{
		"us-east-1":      "us-east1",
		"us-east-2":      "us-east4",
		"us-west-1":      "us-west1",
		"us-west-2":      "us-west2",
		"ca-central-1":   "northamerica-northeast1",
		"sa-east-1":      "southamerica-east1",
		"eu-central-1":   "europe-west3",
		"eu-west-1":      "europe-west1",
		"eu-west-2":      "europe-west2",
		"eu-west-3":      "europe-west9",
		"eu-north-1":     "europe-north1",
		"me-south-1":     "me-west1",
		"af-south-1":     "africa-north1",
		"ap-east-1":      "asia-east2",
		"ap-south-1":     "asia-south1",
		"ap-northeast-3": "asia-northeast2",
		"ap-northeast-2": "asia-northeast3",
		"ap-southeast-1": "asia-southeast1",
		"ap-southeast-2": "australia-southeast1",
		"ap-northeast-1": "asia-northeast1",
	}
)

type GCPVM struct {
	credential    *compute.Service
	config        *config.Config
	sessionClient *sessionClient.SessionClient
	iamClient     *iam.OAuth20Service
}

func New(cfg *config.Config, sessionClient *sessionClient.SessionClient, iamClient *iam.OAuth20Service) *GCPVM {
	ctx := context.Background()

	client, err := compute.NewService(ctx, option.WithCredentialsFile("./service-account-key.json"))
	if err != nil {
		log.Fatalf("Failed to create Compute Engine service client: %v", err)
	}

	return &GCPVM{
		config:        cfg,
		credential:    client,
		sessionClient: sessionClient,
		iamClient:     iamClient,
	}
}

func (g *GCPVM) CreateGameSession(rootScope *envelope.Scope, fleetAlias, sessionID, sessionData, location string, maxPlayer int, namespace string, deployment string) (*sessiondsm.ResponseCreateGameSession, error) {
	scope := rootScope.NewChildScope("gcpvm.CreateGameSession")
	defer scope.Finish()

	locationGCP, err := mapAWSRegionToGCPRegion(location)
	if err != nil {
		return nil, err
	}

	zone, err := selectRandomZone(locationGCP)
	if err != nil {
		return nil, err
	}

	instanceName := getInstanceName(sessionID, namespace)

	instance := &compute.Instance{
		Name:        instanceName,
		MachineType: "zones/" + zone + "/machineTypes/" + g.config.GCPMachineType, // Replace <ZONE> and <MACHINE_TYPE> with your desired zone and machine type.
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				DeviceName: sessionID,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					DiskSizeGb:  10,
					DiskType:    "projects/accelbyte-152423/zones/" + zone + "/diskTypes/pd-balanced",
					SourceImage: "projects/cos-cloud/global/images/cos-stable-113-18244-85-5",
				},
				Mode: "READ_WRITE",
				Type: "PERSISTENT",
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				AccessConfigs: []*compute.AccessConfig{
					{
						Name:        "External NAT",
						NetworkTier: "PREMIUM",
					},
				},
				StackType:  "IPV4_ONLY",
				Subnetwork: "projects/" + g.config.GCPProjectID + "/regions/" + locationGCP + "/subnetworks/" + g.config.GCPNetwork,
			},
		},
		// ServiceAccounts: []*compute.ServiceAccount{ if using service account can adding in here
		// 	{
		// 		Email: "",
		// 		Scopes: []string{
		// 			"https://www.googleapis.com/auth/cloud-platform",
		// 		},
		// 	},
		// },
		ShieldedInstanceConfig: &compute.ShieldedInstanceConfig{
			EnableIntegrityMonitoring: true,
			EnableSecureBoot:          false,
			EnableVtpm:                true,
		},

		ReservationAffinity: &compute.ReservationAffinity{
			ConsumeReservationType: "ANY_RESERVATION",
		},
		ConfidentialInstanceConfig: &compute.ConfidentialInstanceConfig{
			EnableConfidentialCompute: false,
		},

		Tags: &compute.Tags{
			Items: []string{
				"http-server",
				"https-server",
			},
		},

		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   "gce-container-declaration",
					Value: createContainerDeclaration(instanceName, g.config.GCPRepository, deployment),
				},
			},
		},
	}

	_, err = g.credential.Instances.Insert(g.config.GCPProjectID, zone, instance).Context(scope.Ctx).Do()
	if err != nil {
		scope.Log.Errorf("Failed to create VM instance: %v", err)
		return nil, err
	}

	// Call the Compute Engine API to get the instance details.
	instance, err = g.credential.Instances.Get(g.config.GCPProjectID, zone, instanceName).Context(scope.Ctx).Do()
	if err != nil {
		scope.Log.Errorf("Failed to get instance details: %v", err)
		return nil, err
	}

	if instance.Status != "RUNNING" {
		for i := 0; i < g.config.GCPRetry; i++ {
			instance, err = g.credential.Instances.Get(g.config.GCPProjectID, zone, instanceName).Context(scope.Ctx).Do()
			if err != nil {
				scope.Log.Errorf("Failed to get instance details: %v", err)
				return nil, err
			}

			if instance.Status == "RUNNING" {
				break
			}

			time.Sleep(time.Duration(g.config.GCPWaitGetIP) * time.Second)
		}
	}

	// Retrieve the external IP address of the instance.
	var externalIP string
	for _, iface := range instance.NetworkInterfaces {
		if len(iface.AccessConfigs) > 0 {
			externalIP = iface.AccessConfigs[0].NatIP
			break
		}
	}

	if externalIP == "" {
		_, err := g.TerminateGameSession(rootScope, sessionID, namespace, zone)
		if err != nil {
			scope.Log.Errorf("Failed to remove vm %s", instanceName)
		}

		scope.Log.Errorf("Failed to find external IP address for instance %s", instanceName)
		return nil, errors.New("failed to find external ip address")
	}

	return &sessiondsm.ResponseCreateGameSession{
		SessionId:     sessionID,
		Namespace:     namespace,
		SessionData:   sessionData,
		Status:        constants.ServerStatusReady,
		Ip:            externalIP,
		Port:          8080, //can change this port with your docker port open container
		ServerId:      instanceName,
		Source:        constants.GameServerSourceGCP,
		Deployment:    deployment,
		Region:        location,
		ClientVersion: "",
		GameMode:      "",
		CreatedRegion: zone,
	}, nil
}

func (g *GCPVM) TerminateGameSession(rootScope *envelope.Scope, sessionID, namespace, zone string) (*sessiondsm.ResponseTerminateGameSession, error) {
	scope := rootScope.NewChildScope("gcpvm.TerminateGameSession")
	defer scope.Finish()
	res, err := g.credential.Instances.Delete(g.config.GCPProjectID, zone, getInstanceName(sessionID, namespace)).Context(scope.Ctx).Do()
	if err != nil {
		if strings.Contains(err.Error(), "Error 404") {
			return &sessiondsm.ResponseTerminateGameSession{
				SessionId: sessionID,
				Namespace: namespace,
				Success:   true,
				Reason:    "",
			}, nil
		}
		return nil, err
	}

	return &sessiondsm.ResponseTerminateGameSession{
		SessionId: sessionID,
		Namespace: namespace,
		Success:   true,
		Reason:    res.StatusMessage,
	}, nil
}

func (g *GCPVM) UpdateDSInformation(rootScope *envelope.Scope,
	request *session.UpdateGamesessionDSInformationRequest, namespace, sessionID string) (int, error) {
	scope := rootScope.NewChildScope("gcpvm.UpdateDSInformation")
	defer scope.Finish()

	return g.sessionClient.RequestAdminUpdateDSInformation(scope, request, namespace, sessionID)
}

func createContainerDeclaration(instanceName, repository, deploymentImage string) *string {
	declaration := fmt.Sprintf("spec:\n  containers:\n  - name: %s\n    image: %s/%s\n    env:\n    - name: SESSION_ID\n      value: %s\n    securityContext:\n      privileged: true\n    stdin: true\n    tty: true\n  restartPolicy: Never\n# This container declaration format is not public API and may change without notice. Please\n# use gcloud command-line tool or Google Cloud Console to run Containers on Google Compute Engine.", instanceName, repository, deploymentImage, instanceName)
	return &declaration
}

func mapAWSRegionToGCPRegion(awsRegion string) (string, error) {
	if gcpRegion, ok := awsToGCPRegionMap[awsRegion]; ok {
		return gcpRegion, nil
	}
	return "", fmt.Errorf("unknown AWS region: %s", awsRegion)
}

func selectRandomZone(region string) (string, error) {
	zones, ok := gcpRegionZones[region]
	if !ok {
		return "", fmt.Errorf("unknown GCP region: %s", region)
	}

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	randomIndex := rand.Intn(len(zones))
	return zones[randomIndex], nil
}

func getInstanceName(sessionID, namespace string) string {
	return namespace + "-" + sessionID
}
