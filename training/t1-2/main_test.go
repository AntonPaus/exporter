// sum_test.go
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullName(t *testing.T) {
	type fields struct {
		F string
		L string
	}
	tests := []struct { // добавляем слайс тестов
		name   string
		values fields
		want   string
	}{
		{
			name:   "1",
			values: fields{"Misha", "Popov"},
			want:   "Misha Popov",
		},
		{
			name:   "2",
			values: fields{"Misha", "Popov"},
			want:   "Misha Popov",
		},
		{
			name:   "3",
			values: fields{"Pablo " + "Rui", "z Picasso"},
			want:   "Pablo Rui z Picasso",
		},
	}
	for _, test := range tests { // цикл по всем тестам
		t.Run(test.name, func(t *testing.T) {
			u := User{
				FirstName: test.values.F,
				LastName:  test.values.L,
			}
			assert.Equal(t, u.FullName(), test.want)
		})
	}
}
