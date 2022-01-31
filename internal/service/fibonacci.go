package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	rds "github.com/dmitrykharchenko95/fibonacci/internal/rds"
	"github.com/go-redis/redis/v8"
)

var (
	ErrTimeoutExit = errors.New("timeout exit")
)

// GetFibonacci при успешном завершении возвращает срез чисел Фибоначчи с порядковыми номерами от x до y и ошибку nil.
// Через аргумент timeout передается максимальное время работы функции. После истечения timeout функция вернет срез
// чисел Фибоначчи, которые успела вычислить, и ошибку вида "timeout exit: returned <N> values from <M>", где
// N - количество вычисленных чисел Фибоначчи;
// M - ожидаемое количество чисел Фибоначчи.
// При возникновении maxRedisErrors ошибок в работе Redis GetFibonacci производит вычисление каждого значения через
// функцию fibonacci
func GetFibonacci(x, y int, timeout time.Duration, rdb *rds.Client) ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var (
		res         = make([]int64, 0, y-x)
		stopCh      = make(chan struct{})
		redisAtWork = 0
	)

	for i := x; i <= y; i++ {
		var once sync.Once
		var num int

		if redisAtWork <= rdb.MaxErrors {
			val, err := rdb.Cl.Get(ctx, strconv.Itoa(i)).Result()
			switch {
			case errors.Is(err, redis.Nil):
				num = fibonacci(ctx, i, stopCh, &once)
				err = rdb.Cl.Set(ctx, strconv.Itoa(i), strconv.Itoa(num), rdb.Expiration).Err()
				if err != nil {
					redisAtWork += 1
					log.Printf("Redis Set error #%v: %v\n", redisAtWork, err)
				} /*else { 							// раскомментировать для логирования при добавлении значения в Redis
					log.Printf("val %v with key %v set in Redis", num, i)
				}*/
			case err != nil:
				log.Printf("Redis Get error #%v: %v\n", redisAtWork, err)
				redisAtWork += 1
				num = fibonacci(ctx, i, stopCh, &once)
			default:
				num, err = strconv.Atoi(val)
				if err != nil {
					log.Printf("wrong value with key '%v' in Redis: %v\n", val, err)
					num = fibonacci(ctx, i, stopCh, &once)
				} /*else {							// раскомментировать для логирования при получении значения из Redis
					log.Printf("val %v with key %v got from Redis", num, i)
				}*/
			}
		} else {
			num = fibonacci(ctx, i, stopCh, &once)
		}

		select {
		case <-stopCh:
			log.Printf("timeout exit: returned %v values from %v\n", len(res), y-x+1)
			return res, fmt.Errorf("%w: returned %v values from %v", ErrTimeoutExit, len(res), y-x+1)
		default:
			res = append(res, int64(num))
		}
	}
	return res, nil
}

// fibonacci вычисляет число Фибоначчи под порядковым номером n. Выполнение функции fibonacci можно прервать через
// ctx. При преждевременном завершении функции через ctx закрывается сигнальный канал stopCh. Аргумент once передается
// для возможности корректного преждевременного завершения функции и закрытия канала stopCh.
func fibonacci(ctx context.Context, n int, stopCh chan struct{}, once *sync.Once) int {
	select {
	case <-ctx.Done():
		once.Do(func() {
			close(stopCh)
		})
		return 0
	default:
		switch {
		case n == 0:
			return 0
		case n == 1 || n == -1:
			return 1
		case n < 0:
			return fibonacci(ctx, n+2, stopCh, once) - fibonacci(ctx, n+1, stopCh, once)
		default:
			return fibonacci(ctx, n-1, stopCh, once) + fibonacci(ctx, n-2, stopCh, once)
		}
	}
}
