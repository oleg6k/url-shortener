package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewService(t *testing.T) {
	type args struct {
		storage *Storage
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewService(tt.args.storage), "NewService(%v)", tt.args.storage)
		})
	}
}

func TestService_Ping(t *testing.T) {
	type fields struct {
		storage *Storage
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				storage: tt.fields.storage,
			}
			tt.wantErr(t, service.Ping(), fmt.Sprintf("Ping()"))
		})
	}
}

func TestService_getHashByURL(t *testing.T) {
	type fields struct {
		storage *Storage
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				storage: tt.fields.storage,
			}
			got, err := service.getHashByURL(tt.args.url)
			if !tt.wantErr(t, err, fmt.Sprintf("getHashByURL(%v)", tt.args.url)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getHashByURL(%v)", tt.args.url)
		})
	}
}

func TestService_getURLByHash(t *testing.T) {
	type fields struct {
		storage *Storage
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				storage: tt.fields.storage,
			}
			assert.Equalf(t, tt.want, service.getURLByHash(tt.args.hash), "getURLByHash(%v)", tt.args.hash)
		})
	}
}
