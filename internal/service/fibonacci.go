package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	rds "github.com/dmitrykharchenko95/fibonacci/internal/rds"
	"github.com/go-redis/redis/v8"
)

var (
	ErrTimeoutExit = errors.New("timeout exit")
	redisAtWork    = 0
	UseRedis       = true
)

// GetFibonacci при успешном завершении возвращает срез чисел Фибоначчи, форматированных в строки, с порядковыми
// номерами от x до y и ошибку nil.Через аргумент timeout передается максимальное время работы функции. После истечения
// timeout функция вернет срез чисел Фибоначчи, которые успела вычислить, и ошибку вида
// "timeout exit: returned <N> values from <M>", где
// N - количество вычисленных чисел Фибоначчи;
// M - ожидаемое количество чисел Фибоначчи.
// При возникновении maxRedisErrors ошибок в работе Redis GetFibonacci производит вычисление каждого значения через
// функцию fibonacci
func GetFibonacci(x, y int, timeout time.Duration, rdb *rds.Client) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var (
		res    = make([]string, 0, y-x)
		stopCh = make(chan struct{})
	)

	for i := x; i <= y; i++ {
		I := big.NewInt(int64(i))
		var num = new(big.Int)

		if UseRedis && rdb.MaxErrors != 0 {
			val, err := rdb.Cl.Get(ctx, I.Text(10)).Result()
			switch {
			case errors.Is(err, redis.Nil):
				num = fibonacci(ctx, I, stopCh)
				err = rdb.Cl.Set(ctx, I.Text(10), num.Text(10), rdb.Expiration).Err()
				if err != nil {
					redisAtWork += 1
					if redisAtWork >= rdb.MaxErrors {
						UseRedis = false
						log.Println("Redis disabled")
					}
					log.Printf("Redis Set error #%v: %v\n", redisAtWork+1, err)
				} /*else { 							// раскомментировать для логирования при добавлении значения в Redis
					log.Printf("val %v with key %v set in Redis", num, i)
				}*/
			case err != nil:
				log.Printf("Redis Get error #%v: %v\n", redisAtWork+1, err)
				redisAtWork += 1
				if redisAtWork >= rdb.MaxErrors {
					UseRedis = false
					log.Println("Redis disabled")
				}
				num = fibonacci(ctx, I, stopCh)
			default:
				ok := true
				num, ok = num.SetString(val, 10)
				if !ok {
					log.Printf("wrong value with key '%v' in Redis\n", I.Text(10))
					num = fibonacci(ctx, I, stopCh)
				} /*else {							// раскомментировать для логирования при получении значения из Redis
					log.Printf("val %v with key %v got from Redis", num, i)
				}*/
			}
		} else {
			num = fibonacci(ctx, I, stopCh)
		}

		select {
		case <-stopCh:
			log.Printf("timeout exit: returned %v values from %v\n", len(res), y-x+1)
			return res, fmt.Errorf("%w: returned %v values from %v", ErrTimeoutExit, len(res), y-x+1)
		default:
			res = append(res, num.Text(10))
		}
	}
	return res, nil
}

// fibonacci вычисляет число Фибоначчи под порядковым номером n. Выполнение функции fibonacci можно прервать через
// ctx. При преждевременном завершении функции через ctx закрывается сигнальный канал stopCh.
func fibonacci(ctx context.Context, n *big.Int, stopCh chan struct{}) *big.Int {
	negative := false

	f2 := big.NewInt(0)
	f1 := big.NewInt(1)

	if n.Sign() < 0 {
		n.Mul(n, big.NewInt(-1))
		negative = true
	}

	switch {
	case n.Cmp(big.NewInt(0)) == 0:
		return f2
	case n.CmpAbs(big.NewInt(1)) == 0:
		return f1
	default:
		for i := 2; n.Cmp(big.NewInt(int64(i))) >= 0; i++ {
			select {
			case <-ctx.Done():
				close(stopCh)
				return f1
			default:
				next := big.NewInt(0)
				next.Add(f2, f1)
				f2 = f1
				f1 = next
			}
		}
	}

	if negative && big.NewInt(0).Rem(n, big.NewInt(2)).Sign() == 0 {
		f1 = f1.Mul(f1, big.NewInt(-1))
	}

	return f1
}
