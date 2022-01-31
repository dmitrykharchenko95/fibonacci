package httpserver


import (
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/dmitrykharchenko95/fibonacci/internal/rds"
)

type Server struct {
	srv *http.Server
	rdb *rds.Client
	addr string
	timeout time.Duration
}

// New создает новый объект типа Server, который будет прослушивать адрес host:httpPort. Аргумент timeout устанавливает
// максимальное время работы функции service.GetFibonacci, вызываемой в хэндлере getFib.
func New(host, port string, timeout time.Duration, rdb *rds.Client) *Server {

	addr := net.JoinHostPort(host, port)
	return &Server{
		srv: &http.Server{
			Addr: addr,
		},
		rdb: rdb,
		addr: addr,
		timeout: timeout,
	}
}

// Start запускает http сервер.
func (s *Server) Start () error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.getFib)
	s.srv.Handler = mux

	log.Printf("Start http server on %s...\n", s.addr)
	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// Stop останавливает http сервер.
func (s *Server) Stop() error {
	log.Printf("Stop http server...")

	return s.srv.Close()
}

