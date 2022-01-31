package service

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

func Test_fibonacci(t *testing.T) {
	type args struct {
		ctx    context.Context
		n      int
		stopCh chan struct{}
		once   *sync.Once
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "0",
			args: args{
				ctx:    context.Background(),
				n:      0,
				stopCh: make(chan struct{}),
				once:   &sync.Once{},
			},
			want: 0,
		},
		{
			name: "1",
			args: args{
				ctx:    context.Background(),
				n:      1,
				stopCh: make(chan struct{}),
				once:   &sync.Once{},
			},
			want: 1,
		},
		{
			name: "-1",
			args: args{
				ctx:    context.Background(),
				n:      -1,
				stopCh: make(chan struct{}),
				once:   &sync.Once{},
			},
			want: 1,
		},
		{
			name: "8",
			args: args{
				ctx:    context.Background(),
				n:      8,
				stopCh: make(chan struct{}),
				once:   &sync.Once{},
			},
			want: 21,
		},
		{
			name: "-8",
			args: args{
				ctx:    context.Background(),
				n:      -8,
				stopCh: make(chan struct{}),
				once:   &sync.Once{},
			},
			want: -21,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibonacci(tt.args.ctx, tt.args.n, tt.args.stopCh, tt.args.once); got != tt.want {
				t.Errorf("fibonacci() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFibonacci(t *testing.T) {
	type args struct {
		x       int
		y       int
		timeout time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "0 10",
			args: args{
				x:       0,
				y:       10,
				timeout: time.Second * 3,
			},
			want:    []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55},
			wantErr: false,
		},
		{
			name: "-10 0",
			args: args{
				x:       -10,
				y:       0,
				timeout: time.Second * 3,
			},
			want:    []int{-55, 34, -21, 13, -8, 5, -3, 2, -1, 1, 0},
			wantErr: false,
		},
		{
			name: "-5 5",
			args: args{
				x:       -5,
				y:       5,
				timeout: time.Second * 3,
			},
			want:    []int{5, -3, 2, -1, 1, 0, 1, 1, 2, 3, 5},
			wantErr: false,
		},
		{
			name: "timeout exit",
			args: args{
				x:       100000,
				y:       100000,
				timeout: time.Second * 1,
			},
			want:    []int{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFibonacci(tt.args.x, tt.args.y, tt.args.timeout, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFibonacci() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFibonacci() got = %v, want %v", got, tt.want)
			}
		})
	}
}
