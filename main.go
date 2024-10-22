package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"go.opentelemetry.io/otel/trace"
	"net"
	"os"
	"runtime"
	"session-dsm-grpc-plugin/pkg/client/awsgamelift"
	"session-dsm-grpc-plugin/pkg/client/gcpvm"
	"session-dsm-grpc-plugin/pkg/common"
	"session-dsm-grpc-plugin/pkg/config"
	sessiondsm "session-dsm-grpc-plugin/pkg/pb"
	serverDemo "session-dsm-grpc-plugin/pkg/server/demo"
	serverGamelift "session-dsm-grpc-plugin/pkg/server/gamelift"
	serverGCP "session-dsm-grpc-plugin/pkg/server/gcpvm"
	"strings"

	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/factory"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/iam"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os/signal"

	sdkAuth "github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/utils/auth"
	promgrpc "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	prometheusCollectors "github.com/prometheus/client_golang/prometheus/collectors"
)

const (
	id              = int64(1)
	environment     = "production"
	metricsEndpoint = "/metrics"
	metricsPort     = 8080
	grpcPort        = 6565
)

var (
	serviceName = common.GetEnv("OTEL_SERVICE_NAME", "RevocationServiceGoServerDocker")
	logLevelStr = common.GetEnv("LOG_LEVEL", logrus.InfoLevel.String())
)

func main() {
	go func() {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(10)
	}()

	logrus.Infof("starting app server..")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logrusLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		logrusLevel = logrus.InfoLevel
	}
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrusLevel)

	loggingOptions := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadSent),
		logging.WithFieldsFromContext(func(ctx context.Context) logging.Fields {
			if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
				return logging.Fields{"traceID", span.TraceID().String()}
			}

			return nil
		}),
		logging.WithLevels(logging.DefaultClientCodeToLevel),
		logging.WithDurationField(logging.DurationToDurationField),
	}

	srvMetrics := promgrpc.NewServerMetrics()
	unaryServerInterceptors := []grpc.UnaryServerInterceptor{
		srvMetrics.UnaryServerInterceptor(),
		logging.UnaryServerInterceptor(common.InterceptorLogger(logrusLogger), loggingOptions...),
	}
	streamServerInterceptors := []grpc.StreamServerInterceptor{
		srvMetrics.StreamServerInterceptor(),
		logging.StreamServerInterceptor(common.InterceptorLogger(logrusLogger), loggingOptions...),
	}

	if strings.ToLower(common.GetEnv("PLUGIN_GRPC_SERVER_AUTH_ENABLED", "false")) == "true" {
		configRepo := sdkAuth.DefaultConfigRepositoryImpl()
		tokenRepo := sdkAuth.DefaultTokenRepositoryImpl()
		common.OAuth = &iam.OAuth20Service{
			Client:           factory.NewIamClient(configRepo),
			ConfigRepository: configRepo,
			TokenRepository:  tokenRepo,
		}

		common.OAuth.SetLocalValidation(true)

		unaryServerInterceptors = append(unaryServerInterceptors, common.UnaryAuthServerIntercept)
		streamServerInterceptors = append(streamServerInterceptors, common.StreamAuthServerIntercept)
		logrus.Infof("added auth interceptors")
	}

	// Create gRPC Server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryServerInterceptors...),
		grpc.ChainStreamInterceptor(streamServerInterceptors...),
	)

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

	err = env.Parse(cfg)
	if err != nil {
		logrus.Error("unable to parse environment variables: ", err)

		return
	}

	grpcServer = grpc.NewServer()

	switch cfg.DsProvider {
	case "GAMELIFT":
		logrus.Infof("Session Dsms Grpc Plugin: %v", cfg.DsProvider)

		clientGamelift := awsgamelift.New(nil, cfg.GameliftRegion)
		dsmServiceGamelift := &serverGamelift.SessionDSM{
			UnimplementedSessionDsmServer: sessiondsm.UnimplementedSessionDsmServer{},
			ClientGamelift:                clientGamelift,
		}

		sessiondsm.RegisterSessionDsmServer(grpcServer, dsmServiceGamelift)

	case "GCP":
		logrus.Infof("Session Dsms Grpc Plugin: %v", cfg.DsProvider)

		clientGCPVM := gcpvm.New(cfg)
		dsmServiceGCP := &serverGCP.SessionDSM{
			UnimplementedSessionDsmServer: sessiondsm.UnimplementedSessionDsmServer{},
			ClientGCP:                     clientGCPVM,
		}

		sessiondsm.RegisterSessionDsmServer(grpcServer, dsmServiceGCP)

	default:
		logrus.Infof("Session Dsms Grpc Plugin: %v", cfg.DsProvider)

		dsmServiceDemo := &serverDemo.SessionDSM{
			UnimplementedSessionDsmServer: sessiondsm.UnimplementedSessionDsmServer{},
		}

		sessiondsm.RegisterSessionDsmServer(grpcServer, dsmServiceDemo)
	}

	// Enable gRPC Reflection
	reflection.Register(grpcServer)
	logrus.Infof("gRPC reflection enabled")

	// Enable gRPC Health Check
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	// Register Prometheus Metrics
	srvMetrics.InitializeMetrics(grpcServer)
	prometheusRegistry := prometheus.NewRegistry()
	prometheusRegistry.MustRegister(
		prometheusCollectors.NewGoCollector(),
		prometheusCollectors.NewProcessCollector(prometheusCollectors.ProcessCollectorOpts{}),
		srvMetrics,
	)

	go func() {
		http.Handle(metricsEndpoint, promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{}))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", metricsPort), nil))
	}()
	logrus.Infof("serving prometheus metrics at: (:%d%s)", metricsPort, metricsEndpoint)

	// Set Tracer Provider
	tracerProvider, err := common.NewTracerProvider(serviceName, environment, id)
	if err != nil {
		logrus.Fatalf("failed to create tracer provider: %v", err)

		return
	}
	otel.SetTracerProvider(tracerProvider)
	defer func(ctx context.Context) {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			logrus.Fatal(err)
		}
	}(ctx)
	logrus.Infof("set tracer provider: (name: %s environment: %s id: %d)", serviceName, environment, id)

	// Set Text Map Propagator
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			b3.New(),
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	logrus.Infof("set text map propagator")

	// Start gRPC Server
	logrus.Infof("starting gRPC server..")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logrus.Fatalf("failed to listen to tcp:%d: %v", grpcPort, err)

		return
	}
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			logrus.Fatalf("failed to run gRPC server: %v", err)

			return
		}
	}()
	logrus.Infof("gRPC server started")
	logrus.Infof("app server started")

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}
