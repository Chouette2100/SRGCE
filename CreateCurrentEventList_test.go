package main

import (
	"net/http"

	"reflect"
	"testing"
)

func TestCreateCurrentEventList(t *testing.T) {
	type args struct {
		client *http.Client
		status int
	}
	tests := []struct {
		name    string
		args    args
		wantTop *OngoingEvent
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test1",
			args: args{status: 3},
			wantTop: &OngoingEvent{
				0,
				0,
				"",
				nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTop, err := CreateCurrentEventList(tt.args.client, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCurrentEventList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTop, tt.wantTop) {
				t.Errorf("CreateCurrentEventList() = %v, want %v", gotTop, tt.wantTop)
			}
		})
	}
}
