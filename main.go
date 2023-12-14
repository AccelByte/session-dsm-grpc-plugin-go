package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"session-dsm-grpc-plugin/pkg/awsgamelift"
	"session-dsm-grpc-plugin/pkg/config"
	sessiondsm "session-dsm-grpc-plugin/pkg/pb"
	"session-dsm-grpc-plugin/pkg/server"
	"session-dsm-grpc-plugin/pkg/utils"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
)

var (
	//nolint:gochecknoglobals
	buildDate = "unknown"

	//nolint:gochecknoglobals
	revisionID = "unknown"

	//nolint:gochecknoglobals
	gitHash = "unknown"

	//nolint:gochecknoglobals
	rolesSeedingVersion = "unknown"
)

const serviceName = "session-dsm-grpc-plugin"

func main() {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill) //nolint:staticcheck
	rootCtx, cancel := context.WithCancel(context.Background())
	scope := envelope.NewRootScope(rootCtx, serviceName, utils.MakeTraceID(serviceName))
	defer scope.Finish()

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("Justice Session Service")
	logrus.Infof("RevisionID: %s, Build Date: %s, Git Hash: %s Roles Seeding Version: %s\n", revisionID, buildDate, gitHash, rolesSeedingVersion)

	cfg := &config.Config{}

	flag.Usage = func() {
		flag.CommandLine.SetOutput(os.Stdout)
		for _, val := range cfg.HelpDocs() {
			//nolint:forbidigo
			fmt.Println(val)
		}

		//nolint:forbidigo
		fmt.Println("")
		flag.PrintDefaults()
	}
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		logrus.Error("unable to parse environment variables: ", err)

		return
	}

	clientGamelift := awsgamelift.New(nil, cfg.GameliftRegion)
	grpcServer := grpc.NewServer()
	dsmService := &server.SessionDSM{
		UnimplementedSessionDsmServer: sessiondsm.UnimplementedSessionDsmServer{},
		ClientGamelift:                clientGamelift,
	}

	sessiondsm.RegisterSessionDsmServer(grpcServer, dsmService)

	gRPCListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.GRPCPort))
	if err != nil {
		logrus.Fatalf("unable to create gRPC listener: %v", err)
	}

	go func() {
		logrus.Infof("Serving gRPC, listens at 0.0.0.0:%d", cfg.GRPCPort)
		if errServeGRPC := grpcServer.Serve(gRPCListener); err != nil {
			logrus.Fatalf("failed to serve gRPC: %v", errServeGRPC)
		}
	}()

	//nolint:gosimple
	select {
	case <-sigCh:
		cancel()
	}
}
