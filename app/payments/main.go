package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	api2 "github.com/didopimentel/go-saga-poc/app/payments/api"
	v12 "github.com/didopimentel/go-saga-poc/app/payments/api/v1"
	"github.com/didopimentel/go-saga-poc/domain/payment"
	"github.com/didopimentel/go-saga-poc/gateways/persistence"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	log, err := logConfig.Build()
	if err != nil {
		fmt.Printf("can't initialize zap logger: %v", err)
		os.Exit(1)
	}

	startApp(log)

	defer func() { _ = log.Sync() }() //nolint:errcheck
}

func startApp(log *zap.Logger) {
	var AppVersion = "development"
	type config struct {
		SVAddr            string        `conf:"env:SV_ADDR,default:0.0.0.0:7010"`
		SVReadTimeout     time.Duration `conf:"env:SV_READ_TIMEOUT,default:30s"`
		SVWriteTimeout    time.Duration `conf:"env:SV_WRITE_TIMEOUT,default:30s"`
		SVMaxConnAge      time.Duration `conf:"env:SV_MAX_CONN_AGE,default:1m"`
		SVMaxConnAgeGrace time.Duration `conf:"env:SV_MAX_CONN_AGE_GRACE,default:5m"`
		PGAddr            string        `conf:"env:PG_ADDR,default:postgres://ps_user:ps_password@localhost:7002/go-saga-poc?sslmode=disable,mask"`
		PGPoolMinConn     int32         `conf:"env:PG_POOL_MIN_CONN,default:20"`
		PGPoolMaxConn     int32         `conf:"env:PG_POOL_MAX_CONN,default:100"`
		Version           conf.Version
	}

	cfg := config{}
	cfg.Version = conf.Version{
		SVN:  AppVersion,
		Desc: "Payments API",
	}

	if err := conf.Parse(os.Args[1:], "PAYMENTS", &cfg); err != nil {
		switch {
		case errors.Is(err, conf.ErrHelpWanted):
			var usageErr error
			usage, usageErr := conf.Usage("PAYMENTS", &cfg)
			if usageErr != nil {
				log.Fatal(fmt.Errorf("generating config usage: %w", usageErr).Error())
			}
			fmt.Println(usage)
		case errors.Is(err, conf.ErrVersionWanted):
			var versionErr error
			version, versionErr := conf.VersionString("PAYMENTS", &cfg)
			if versionErr != nil {
				log.Fatal(fmt.Errorf("generating config version: %w", versionErr).Error())
			}
			fmt.Println(version)
		}

		log.Fatal(fmt.Errorf("parsing config: %w", err).Error())
	}

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("generating config for output: %w", err).Error())
	}
	log.Info("Config", zap.String("config", out))

	ctx := context.Background()

	//
	// Gateways
	//

	//schedulerGRPCConn, err := createGRPCConn(ctx, cfg.SchedulerAddr, tracer)
	//if err != nil {
	//	log.Fatal("failed to connect to scheduler API",
	//		zap.String("address", cfg.SchedulerAddr),
	//		zap.Error(err),
	//	)
	//}
	//schedulerAPIClient := schedulerV1.NewSchedulerServiceClient(schedulerGRPCConn)

	// Postgres
	txManager, err := persistence.NewTxManager(ctx, cfg.PGAddr, cfg.PGPoolMinConn, cfg.PGPoolMaxConn)
	if err != nil {
		log.Fatal("failed to instantiate pg", zap.Error(err))
	}

	repository := getRepository(txManager)

	//
	// UseCases
	//

	createPaymentUseCase := payment.NewCreatePaymentUseCase(repository.Payments)
	deletePaymentUseCase := payment.NewDeletePaymentUseCase(repository.Payments)

	paymentsUseCases := &struct {
		*payment.CreatePaymentUseCase
		*payment.DeletePaymentUseCase
	}{
		createPaymentUseCase,
		deletePaymentUseCase,
	}
	paymentsAPI := &v12.API{
		PaymentsAPI: v12.NewPaymentsAPI(paymentsUseCases),
		Repository:  repository,
	}

	svs := api2.Settings{
		Addr:            cfg.SVAddr,
		Server:          paymentsAPI,
		ReadTimeout:     cfg.SVReadTimeout,
		WriteTimeout:    cfg.SVWriteTimeout,
		MaxConnAgeGrace: cfg.SVMaxConnAgeGrace,
		MaxConnAge:      cfg.SVMaxConnAge,
		Logger:          log,
	}

	sv, err := api2.New(ctx, svs)
	if err != nil {
		log.Fatal("failed to instantiate service")
	}

	defer func() { _ = sv.Shutdown(ctx) }() //nolint:errcheck
	// errch is used to signal when any of our Listen/Serve goroutines stop.
	// The program ends on the first error, nil or not.
	errch := make(chan error)
	go func() {
		errch <- fmt.Errorf("service's ListenAndServe failed. %w", sv.ListenAndServe())
	}()

	go handleInterrupt(ctx, log, sv)

	log.Info("orders service started")

	err = <-errch
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("%v", zap.Error(err))
	}
}

func getRepository(txManager *persistence.TxManager) *v12.Repository {
	return &v12.Repository{
		Payments: &persistence.Payments{
			Transactioner: txManager,
			Q:             txManager,
		},
		Health: &persistence.Health{Q: txManager},
	}
}

func handleInterrupt(ctx context.Context, log *zap.Logger, ss ...*http.Server) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	sig := <-signals
	log.Info("captured signal", zap.String("signal", sig.String()))
	signal.Stop(signals)
	for _, s := range ss {
		if err := s.Shutdown(ctx); err != nil {
			log.Error("Error on shutdown server", zap.Error(err))
		}
	}
}

func createGRPCConn(
	ctx context.Context,
	address string,
) (*grpc.ClientConn, error) {
	schedulerDNSAddr := fmt.Sprintf("dns:///%s", address)

	return grpc.DialContext(ctx, schedulerDNSAddr,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
}
