package service

import (
	"testing"
)

var service CommandService

func setup() {
	service = GetCommandServiceImpl()
}

func TestCommandServiceImpl_Set(t *testing.T) {
	setup()
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 01",
			args: struct {
				key   string
				value string
			}{key: "foo", value: "123"},
			wantErr: false,
		}, {
			name: "Test 02",
			args: struct {
				key   string
				value string
			}{key: "bar", value: "456"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandServiceImpl_Get(t *testing.T) {
	setup()
	service.Set("foo", "123")
	tests := []struct {
		name      string
		key       string
		wantValue string
		wantErr   bool
	}{
		{name: "Test 01", key: "foo", wantValue: "123", wantErr: false},
		{name: "Test 02", key: "bar", wantValue: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, err := service.Get(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotValue != tt.wantValue {
				t.Errorf("Get() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestCommandServiceImpl_Delete(t *testing.T) {
	setup()
	service.Set("foo", "123")
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{name: "Test 01", key: "foo", wantErr: false},
		{name: "Test 02", key: "foo", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.Delete(tt.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandServiceImpl_Count(t *testing.T) {
	setup()
	service.Set("foo", "123")
	service.Set("bar", "456")
	service.Set("baz", "123")
	tests := []struct {
		name  string
		value string
		want  int
	}{
		{name: "Test 01", value: "432", want: 0},
		{name: "Test 02", value: "123", want: 2},
		{name: "Test 03", value: "456", want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.Count(tt.value); got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandServiceImpl_Commit(t *testing.T) {
	setup()
	service.Begin()
	service.Set("foo", "123")
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Should commit transaction", wantErr: false},
		{name: "Should fail since transaction has been committed, no transaction", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.Commit(); (err != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandServiceImpl_Rollback(t *testing.T) {
	setup()
	service.Set("foo", "123")
	service.Begin()
	service.Set("foo", "456")
	service.Set("bar", "123")
	service.Begin()
	service.Set("foo", "456")
	service.Set("bar", "def")
	tests := []struct {
		name          string
		key           string
		expectedValue string
		wantErr       bool
	}{
		{name: "Test 01", key: "foo", expectedValue: "456", wantErr: false},
		{name: "Test 02", key: "foo", expectedValue: "123", wantErr: false},
		{name: "Test 03", key: "foo", expectedValue: "123", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := service.Rollback(); (err != nil) != tt.wantErr {
				t.Errorf("Rollback() error = %v, wantErr %v", err, tt.wantErr)
			}
			if value, _ := service.Get(tt.key); value != tt.expectedValue {
				t.Errorf("Get() = %v, wantErr %v after Rollback()", value, tt.expectedValue)
			}
		})
	}
}
