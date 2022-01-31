package httpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os/exec"
	"testing"
	"time"

	"github.com/dmitrykharchenko95/fibonacci/internal/rds"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

const (
	httpHost = "localhost"
	httpPort = "8080"

	redisHost       = "localhost"
	redisPort       = "6379"
	redisExpiration = time.Hour
	redisMaxErr     = 6
	timeout         = time.Second * 3
)

var (
	buf                              = &bytes.Buffer{}
	expectedResponse, actualResponse Response
	resBody                          = make([]byte, 0, 20)
)

func TestServer(t *testing.T) {
	rdb := &rds.Client{
		Cl: redis.NewClient(&redis.Options{
			Addr:     net.JoinHostPort(redisHost, redisPort),
			Password: "",
			DB:       0,
		}),
		Expiration: redisExpiration,
		MaxErrors:  redisMaxErr,
	}

	s := New(httpHost, httpPort, timeout, rdb)

	go func() {
		err := s.Start()
		require.NoError(t, err)
	}()

	defer func() {
		err := s.Stop()
		require.NoError(t, err)
	}()

	t.Run("0,10", func(t *testing.T) {
		cmd := exec.Command("curl", "-X", "GET", "-i", "localhost:8080/", `-d`, "0,10")

		cmd.Stdout = buf
		err := cmd.Run()
		require.NoError(t, err)

		for {
			resBody, err = buf.ReadBytes(10)
			if errors.Is(err, io.EOF) {
				break
			}
		}

		expectedResponse.Data = []int64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}

		err = json.Unmarshal(resBody, &actualResponse)
		require.NoError(t, err)

		require.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("-10,0", func(t *testing.T) {
		cmd := exec.Command("curl", "-X", "GET", "-i", "localhost:8080/", `-d`, "-10,0")

		cmd.Stdout = buf
		err := cmd.Run()
		require.NoError(t, err)

		for {
			resBody, err = buf.ReadBytes(10)
			if errors.Is(err, io.EOF) {
				break
			}
		}

		expectedResponse.Data = []int64{-55, 34, -21, 13, -8, 5, -3, 2, -1, 1, 0}

		err = json.Unmarshal(resBody, &actualResponse)
		require.NoError(t, err)

		require.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("-5,5", func(t *testing.T) {
		cmd := exec.Command("curl", "-X", "GET", "-i", "localhost:8080/", `-d`, "-5,5")

		cmd.Stdout = buf
		err := cmd.Run()
		require.NoError(t, err)

		for {
			resBody, err = buf.ReadBytes(10)
			if errors.Is(err, io.EOF) {
				break
			}
		}

		expectedResponse.Data = []int64{5, -3, 2, -1, 1, 0, 1, 1, 2, 3, 5}

		err = json.Unmarshal(resBody, &actualResponse)
		require.NoError(t, err)

		require.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("1,1", func(t *testing.T) {
		cmd := exec.Command("curl", "-X", "GET", "-i", "localhost:8080/", `-d`, "1,1")

		cmd.Stdout = buf
		err := cmd.Run()
		require.NoError(t, err)

		for {
			resBody, err = buf.ReadBytes(10)
			if errors.Is(err, io.EOF) {
				break
			}
		}

		expectedResponse.Data = []int64{1}

		err = json.Unmarshal(resBody, &actualResponse)
		require.NoError(t, err)

		require.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("wrong arguments", func(t *testing.T) {
		cmd := exec.Command("curl", "-X", "GET", "-i", "localhost:8080/", `-d`, "10 20")

		cmd.Stdout = buf
		err := cmd.Run()
		require.NoError(t, err)

		for {
			resBody, err = buf.ReadBytes(10)
			if errors.Is(err, io.EOF) {
				break
			}
		}

		expectedResponse = Response{
			Data: []int64{},
			Err:  ErrWrongArgs.Error(),
		}

		err = json.Unmarshal(resBody, &actualResponse)
		require.NoError(t, err)

		require.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("wrong syntax", func(t *testing.T) {
		cmd := exec.Command("curl", "-X", "GET", "-i", "localhost:8080/", `-d`, "test,test")

		cmd.Stdout = buf
		err := cmd.Run()
		require.NoError(t, err)

		for {
			resBody, err = buf.ReadBytes(10)
			if errors.Is(err, io.EOF) {
				break
			}
		}
		expectedResponse = Response{
			Data: []int64{},
			Err:  "strconv.Atoi: parsing \"test\": invalid syntax",
		}

		err = json.Unmarshal(resBody, &actualResponse)
		require.NoError(t, err)

		require.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("timeout exit", func(t *testing.T) {
		cmd := exec.Command("curl", "-X", "GET", "-i", "localhost:8080/", `-d`, "1000000,1000000")

		cmd.Stdout = buf

		now := time.Now()
		err := cmd.Run()
		require.NoError(t, err)

		execTime := time.Since(now)
		require.LessOrEqual(t, timeout.Milliseconds(), execTime.Milliseconds(), "timeout exit not executed")

		for {
			resBody, err = buf.ReadBytes(10)
			if errors.Is(err, io.EOF) {
				break
			}
		}

		expectedResponse = Response{
			Data: []int64{},
			Err:  "timeout exit: returned 0 values from 1",
		}

		err = json.Unmarshal(resBody, &actualResponse)
		require.NoError(t, err)

		require.Equal(t, expectedResponse, actualResponse)
	})
}
