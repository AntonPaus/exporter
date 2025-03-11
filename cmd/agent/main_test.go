// sum_test.go
package main

import (
	"testing"
	"time"
)

func TestSum(t *testing.T) {
	tests := []struct { // добавляем слайс тестов
		name   string
		values []int
		want   int
	}{
		{
			name:   "simple test #1", // описываем каждый тест:
			values: []int{1, 2},      // значения, которые будет принимать функция,
			want:   3,                // и ожидаемый результат
		},
		{
			name:   "one",
			values: []int{1},
			want:   1,
		},
		{
			name:   "with negative values",
			values: []int{-1, -2, 3},
			want:   0,
		},
		{
			name:   "with negative zero",
			values: []int{-0, 3},
			want:   3,
		},
		{
			name: "a lot of values",
			values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
				14, 15, 16, 17, 18, 18},
			want: 189,
		},
	}
	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			if sum := Sum(test.values...); sum != test.want {
				t.Errorf("Sum() = %d, want %d", sum, test.want)
			}
		})
	}
}

func Test_pollStats(t *testing.T) {
	type args struct {
		stats    chan Metrics
		interval time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pollStats(tt.args.stats, tt.args.interval)
		})
	}
}
