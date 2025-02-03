package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger (
	ctx context.Context, 
	req any, 
	info *grpc.UnaryServerInfo, 
	handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		resp,err = handler(ctx,req)
		duration := time.Since(start)
		
		statusCode := codes.Unknown
		if st,ok := status.FromError(err); ok {
			statusCode = st.Code()
		}

		logger := log.Info()
		if err != nil {
			logger = log.Error().Err(err)
		}

		logger.Str("protocol","grpc").
		Str("method",info.FullMethod).
		Int("status",int(statusCode)).
		Str("status_text",statusCode.String()).
		Dur("duration",duration).
		Msg("received a gRPC request")
		return
		

	}

	type ResponseRecorder struct {
		http.ResponseWriter
		StatusCode int 
		Body []byte
	}

	func (r *ResponseRecorder) WriteHeader(statusCode int) {
		r.StatusCode = statusCode
		r.ResponseWriter.WriteHeader(statusCode)
	}

	func (r *ResponseRecorder) Write(b []byte) (int, error) {
		if r.StatusCode == 0 {
			r.StatusCode = http.StatusOK
		}
		r.Body = b
		return r.ResponseWriter.Write(b)
	}

	func HttpLogger(handler http.Handler) (http.Handler) {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			recorder := &ResponseRecorder{ResponseWriter: w, StatusCode: http.StatusOK}
			handler.ServeHTTP(recorder, r)
			duration := time.Since(start)
			logger := log.Info()

			if recorder.StatusCode != http.StatusOK {
				logger = log.Error()
			}


			logger.Str("protocol","http").
		Str("method",r.Method).
		Str("path",r.RequestURI).
		Int("status",recorder.StatusCode).
		Str("status_text",http.StatusText(recorder.StatusCode)).
		Dur("duration",duration).
		Msg("received a HTTP request")	
	})
	}