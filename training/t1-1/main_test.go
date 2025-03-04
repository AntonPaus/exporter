// sum_test.go
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	tests := []struct { // добавляем слайс тестов
		name  string
		value float64
		want  float64
	}{
		{
			name:  "1",
			value: -3.14,
			want:  3.14,
		},
		{
			name:  "2",
			value: 3.1,
			want:  3.1,
		},
		{
			name:  "3",
			value: -0,
			want:  0,
		},
		{
			name:  "4",
			value: -0.000000003,
			want:  0.000000003,
		},
	}
	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, Abs(test.value))
		})
	}
}
