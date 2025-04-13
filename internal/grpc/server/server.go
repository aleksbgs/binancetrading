package server

import (
	"binancetrading/internal/application/service"
	pb "binancetrading/internal/grpc/proto"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// candlestickServer implements the CandlestickService gRPC interface
type candlestickServer struct {
	pb.UnimplementedCandlestickServiceServer
	svc *service.CandlestickService
}

func (s *candlestickServer) StreamCandlesticks(req *pb.StreamCandlesticksRequest, stream pb.CandlestickService_StreamCandlesticksServer) error {
	symbol := req.Symbol
	candleChan := s.svc.CandlestickChan()

	for candle := range candleChan {
		if candle.Symbol != symbol {
			continue
		}
		pbCandle := &pb.Candlestick{
			Symbol:    candle.Symbol,
			Open:      candle.Open,
			High:      candle.High,
			Low:       candle.Low,
			Close:     candle.Close,
			Volume:    candle.Volume,
			Timestamp: candle.Timestamp.Format("2006-01-02T15:04:05Z"),
		}
		if err := stream.Send(pbCandle); err != nil {
			log.Printf("Failed to stream candlestick for %s: %v", symbol, err)
			return err
		}
	}
	return nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}

// StartGRPCServer starts the gRPC server with the given CandlestickService
func StartGRPCServer(svc *service.CandlestickService) error {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		return err
	}

	creds, err := loadTLSCredentials()
	var s *grpc.Server
	if err != nil {
		log.Printf("Failed to load TLS credentials: %v, using plaintext", err)
		s = grpc.NewServer()
	} else {
		s = grpc.NewServer(grpc.Creds(creds))
	}

	pb.RegisterCandlestickServiceServer(s, &candlestickServer{svc: svc})
	reflection.Register(s)

	log.Println("Starting gRPC server on :50051")
	return s.Serve(lis)
}
