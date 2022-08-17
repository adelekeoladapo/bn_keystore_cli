package controller

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCommandController_Process(t *testing.T) {
	tests := []struct {
		name           string
		commandStrings []string
		wantResponse   []string
		wantErr        bool
	}{
		{
			name: "Should test Set and Get",
			commandStrings: []string{
				"SET foo 123",
				"GET foo",
			},
			wantResponse: []string{
				"",
				"123",
			}, wantErr: false,
		}, {
			name: "Should delete a value",
			commandStrings: []string{
				"SET foo 123",
				"DELETE foo",
				"GET foo",
			},
			wantResponse: []string{
				"",
				"",
				"key not set",
			}, wantErr: false,
		}, {
			name: "Should count the number of occurrences of a value",
			commandStrings: []string{
				"SET foo 123",
				"SET bar 456",
				"SET baz 123",
				"COUNT 123",
				"COUNT 456",
			},
			wantResponse: []string{
				"",
				"",
				"",
				"2",
				"1",
			}, wantErr: false,
		}, {
			name: "Should commit a transaction",
			commandStrings: []string{
				"BEGIN",
				"SET foo 456",
				"COMMIT",
				"ROLLBACK",
				"GET foo",
			},
			wantResponse: []string{
				"",
				"",
				"",
				"no transaction",
				"456",
			}, wantErr: false,
		}, {
			name: "Should rollback a transaction",
			commandStrings: []string{
				"SET foo 123",
				"SET bar abc",
				"BEGIN",
				"SET foo 456",
				"GET foo",
				"SET bar def",
				"GET bar",
				"ROLLBACK",
				"GET foo",
				"GET bar",
				"COMMIT",
			},
			wantResponse: []string{
				"",
				"",
				"",
				"",
				"456",
				"",
				"def",
				"",
				"123",
				"abc",
				"no transaction",
			}, wantErr: false,
		}, {
			name: "Should perform nexted transactions",
			commandStrings: []string{
				"SET foo 123",
				"BEGIN",
				"SET bar 456",
				"SET foo 456",
				"BEGIN",
				"COUNT 456",
				"GET foo",
				"SET foo 789",
				"GET foo",
				"ROLLBACK",
				"GET foo",
				"ROLLBACK",
				"GET foo",
			},
			wantResponse: []string{
				"",
				"",
				"",
				"",
				"",
				"2",
				"456",
				"",
				"789",
				"",
				"456",
				"",
				"123",
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := GetCommandController()
			var actualResponse = make([]string, len(tt.commandStrings))
			for i, cmdStr := range tt.commandStrings {
				if gotResponse, err := ctrl.Process(cmdStr); err != nil {
					actualResponse[i] = fmt.Sprint(err)
				} else {
					actualResponse[i] = gotResponse
				}
			}
			if !reflect.DeepEqual(actualResponse, tt.wantResponse) {
				t.Errorf("Process() gotResponse = %v, want %v", actualResponse, tt.wantResponse)
			}
		})
	}
}
