package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/dmitrykharchenko95/fibonacci/internal/rds"
	"github.com/dmitrykharchenko95/fibonacci/internal/server/grpc/pb"
	"github.com/dmitrykharchenko95/fibonacci/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Server struct {
	srv *grpc.Server
	rdb *rds.Client
	addr string
	timeout time.Duration
	pb.UnimplementedFibonacciServer
}

func New(host, port string, timeout time.Duration, rdb *rds.Client)  *Server{
	return &Server{
		srv: grpc.NewServer(),
		rdb: rdb,
		addr: net.JoinHostPort(host, port),
		timeout: timeout,
	}
}

func (s *Server) Start () error {
	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	pb.RegisterFibonacciServer(s.srv, s)
	log.Printf("Start grpc server on %s...\n", s.addr)
	if err := s.srv.Serve(lsn); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop () {
	log.Println("Stop grpc server...")
	s.srv.Stop()
}

func (s *Server) GetFibonacci(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	var x, y int64
	if req.X > req.Y {
		x, y = req.Y, req.X
	} else {
		x, y = req.X, req.Y
	}

	ip, err := getClientIP(ctx)
	if err != nil {
		log.Printf("can not get client IP: %v", err)
	}

	resp := &pb.Response{}
	data, err := service.GetFibonacci(int(x), int(y), s.timeout,s.rdb)
	if err != nil {
		resp.Data, resp.Err = data, err.Error()
	} else {
		resp.Data = data
	}

	log.Printf("%v: sended %v numbers fibonacci\n", ip, len(resp.Data))
	return resp, nil
}

func getClientIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("couldn't parse client IP address")
	}
	return p.Addr.String(), nil
}