package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dmitrykharchenko95/fibonacci/internal/service"
)

var ErrWrongArgs = errors.New("request's body should has two int values through a comma")

type Response struct {
	Data []int64
	Err  string
}

// writeResponse осуществляет запись структуры Response в http.ResponseWriter в формате JSON.
func writeResponse(w http.ResponseWriter, resp *Response) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
	_, err = w.Write(resBuf)
	if err != nil {
		log.Printf("response write error: %s", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// parseArgs принимает в качестве аргумента in строку вида "A,B", где А и В - целые числа, и при успешном выполнении
// возвращает А и В (если А < B) или В и А (если А > В). При несоответствии in шаблону "A,B", возвращает 0,0 и
// ошибку ErrWrongArgs.
func parseArgs(in string) (int, int, error) {
	args := strings.Split(in, ",")

	if len(args) != 2 {
		return 0, 0, ErrWrongArgs
	}

	x, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, 0, err
	}

	y, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, 0, err
	}

	if x > y {
		x, y = y, x
	}

	return x, y, nil
}

// getFib обрабатывает запросы к серверу и отправляет клиенту структуру Response с результатами выполнения
// service.GetFibonacci в формате JSON. getFib обрабатывает только GET-запросы по адресу "host:port/". В теле запроса
// ожидаются два целых числа через запятую.
func (s *Server) getFib(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	resp := &Response{
		Data: make([]int64, 0),
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp.Err = fmt.Sprintf("method %s not supported on uri %s", r.Method, r.URL.Path)
		writeResponse(w, resp)
		log.Printf("%v: unsupported method <%v> on uri %v\n", r.RemoteAddr, r.Method, r.URL.Path)
		return
	}

	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		w.WriteHeader(http.StatusBadRequest)
		resp.Err = err.Error()
		writeResponse(w, resp)
		log.Printf("%v: reading request body failed: %v\n", r.RemoteAddr, err)
		return
	}

	x, y, err := parseArgs(string(buf))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Err = err.Error()
		writeResponse(w, resp)
		log.Printf("%v: wrong arguments: %v\n", r.RemoteAddr, err)
		return
	}
	data, err := service.GetFibonacci(x, y, s.timeout, s.rdb)
	if err != nil {
		resp.Data, resp.Err = data, err.Error()
	} else {
		resp.Data = data
	}

	writeResponse(w, resp)
	log.Printf("%v: sended %v numbers fibonacci\n", r.RemoteAddr, len(resp.Data))
}
