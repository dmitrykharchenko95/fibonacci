package config

import (
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {

	t.Run("base", func(t *testing.T) {
		got, err := New("../configs/fibonacci_config.json")
		fmt.Println(err)
		fmt.Println(got)
	})
}
