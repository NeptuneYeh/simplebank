package grpcGateway

import (
	"context"
	"github.com/NeptuneYeh/simplebank/internal/grpcApp"
	"github.com/NeptuneYeh/simplebank/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	"time"
)

type Module struct {
	GrpcApi           *grpcApp.Module
	GrpcGatewayServer *http.Server
}

func NewModule() *Module {
	return &Module{
		GrpcApi: grpcApp.NewModule(),
	}
}

func (module *Module) Run(address string) {
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	// ServeMux 是一個請求多路復用器，它負責將 HTTP 請求匹配到特定的路徑模式並調用相應的處理器
	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background()) // cancel 是一個函數，當調用這個函數時，會發送取消信號，通知所有使用這個上下文的操作應該停止
	defer cancel()

	err := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, module.GrpcApi)
	if err != nil {
		log.Fatalf("cannot register handler server: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", enableCors(grpcMux))

	//fs := http.FileServer(http.Dir("./doc/swagger"))
	statikFS, err := fs.New()
	if err != nil {
		log.Fatalf("cannot create statik file system: %v", err)
	}
	// 對外顯示 FQDN/swagger/index.html, 系統內部去掉 /swagger/ 然後在 ./doc/swagger 內尋找 index.html
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	handler := grpcApp.HttpLogger(mux)

	module.GrpcGatewayServer = &http.Server{
		Addr:    address,
		Handler: handler,
	}

	go func() {
		log.Printf("Starting GrpcGatewayServer on %s\n", address)
		if err := module.GrpcGatewayServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run GrpcGatewayServer: %v", err)
		}
	}()
}

func (module *Module) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := module.GrpcGatewayServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to run GrpcGatewayServer shutdown: %v", err)
	}
	return nil
}

func enableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}
