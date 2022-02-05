package api

import (
	"context"
	"errors"
	v12 "github.com/didopimentel/go-saga-poc/protogen/payments/api/v1"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Settings struct {
	Addr            string
	Server          v12.PaymentsAPIServer
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxConnAgeGrace time.Duration
	MaxConnAge      time.Duration
	Logger          *zap.Logger
}

type ProtoErrorHandler struct {
	Logger *zap.Logger
}

func New(ctx context.Context, s Settings) (*http.Server, error) {
	if s.Logger == nil {
		s.Logger = zap.L()
	}
	errorHandler := ProtoErrorHandler{
		Logger: s.Logger,
	}
	keepAlive := grpc.KeepaliveParams(
		keepalive.ServerParameters{
			MaxConnectionAgeGrace: s.MaxConnAgeGrace,
			MaxConnectionAge:      s.MaxConnAge,
		})

	recoveryHandler := func(p interface{}) (err error) {
		errorHandler.Logger.Error("panic recovery", zap.Any("parameter", p))

		err = status.Errorf(codes.Internal, "%v", "something wrong happened")

		return err
	}

	grpcMux := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			grpcrecovery.UnaryServerInterceptor(grpcrecovery.WithRecoveryHandler(recoveryHandler)),
			ErrorInterceptor(s.Logger),
		),
		keepAlive,
	)
	v12.RegisterPaymentsAPIServer(grpcMux, s.Server)

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L55
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcMux.ServeHTTP(w, r)

			return
		}
		// add http method info to server gateway service
		r.Header.Add("Grpc-Metadata-HTTP-Method", r.Method)
	})

	return &http.Server{
		Addr:         s.Addr,
		Handler:      AllowCORS(h2c.NewHandler(rootHandler, &http2.Server{})),
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
	}, nil
}

func (h *ProtoErrorHandler) handle(
	_ context.Context,
	_ *runtime.ServeMux,
	m runtime.Marshaler,
	w http.ResponseWriter,
	_ *http.Request,
	err error,
) {
	e := FromGRPCError(err)
	if errors.Is(err, runtime.ErrNotMatch) {
		e = NewNotFoundError("")
	}

	// Set custom response
	w.Header().Set("Content-type", m.ContentType(nil))
	w.WriteHeader(e.HTTPStatus())
	_, err = w.Write([]byte(e.HTTPError()))
	if err != nil {
		h.Logger.Error("handle failed to write", zap.Error(err))
	}
}
