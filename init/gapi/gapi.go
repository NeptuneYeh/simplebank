package gapi

import (
	"github.com/NeptuneYeh/simplebank/internal/grpcApp"
	"github.com/NeptuneYeh/simplebank/pb"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Module struct {
	GrpcApi    *grpcApp.Module
	GrpcServer *grpc.Server
	Listener   net.Listener
}

var MainGapi *Module

func NewModule() *Module {

	// TODO register logger Interceptor
	grpcLogger := grpc.UnaryInterceptor(grpcApp.GrpcLogger)

	gAPIModule := &Module{
		GrpcApi:    grpcApp.NewModule(),
		GrpcServer: grpc.NewServer(grpcLogger),
	}

	MainGapi = gAPIModule
	return gAPIModule
}

// Run grpc server
func (module *Module) Run(address string) {
	var err error
	module.Listener, err = net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	// Register your services here
	pb.RegisterSimpleBankServer(module.GrpcServer, module.GrpcApi)
	reflection.Register(module.GrpcServer) // 使用反射可以使用 grpcurl debug, prod 不建議用

	go func() {
		log.Printf("Starting gRPC server on %s\n", address)
		if err := module.GrpcServer.Serve(module.Listener); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve")
		}
	}()
}

func (module *Module) Shutdown() error {
	module.GrpcServer.GracefulStop()
	log.Info().Msg("gRPC server gracefully stopped")
	return nil
}
