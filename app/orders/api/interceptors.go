package api

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"net/http"
)

func ErrorInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		res, err := handler(ctx, req)
		if err != nil {
			// handles gRPC errors in a different flow
			if _, ok := status.FromError(err); ok {
				return res, err
			}

			// translating the error to our own subset of errors, in api pkg
			e := FromError(err)
			grpcErr := e.GRPCError()

			if e.HTTPStatus() >= http.StatusInternalServerError {
				log.Error("request failed", zap.Error(err), zap.Int("http_status", e.HTTPStatus()))
			}

			return nil, grpcErr
		}

		return res, nil
	}
}

func AllowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Expose-Headers", "X-Access-Token")
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Content-Type,X-EM-Service")
				w.Header().Set("Access-Control-Allow-Methods", "DELETE,GET,HEAD,PATCH,POST,PUT")

				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
