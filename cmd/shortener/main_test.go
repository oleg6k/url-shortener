package main

import (
	"github.com/oleg6k/url-shortener/internal/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	type storage map[string]string

	tests := []struct {
		name string
		args storage
	}{
		{
			name: "positive test #1",
			args: storage{"hashedUrl": "origUrl"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := app.NewService(test.args)
			assert.IsType(t, &app.Service{}, service)
		})
	}
}
