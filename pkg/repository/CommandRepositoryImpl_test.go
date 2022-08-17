package repository

import (
	command2 "bn_keystore/pkg/command"
	"reflect"
	"testing"
)

func TestCommandRepositoryImpl_Create(t *testing.T) {
	type args struct {
		key     string
		command command2.Command
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 01",
			args: args{
				key: "key",
				command: command2.Set{
					Key:   "key",
					Value: "value",
				},
			},
			wantErr: false,
		}, {
			name: "Test 02",
			args: args{
				key:     "key",
				command: command2.Begin{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := GetCommandRepositoryImpl()
			if err := o.Create(tt.args.key, tt.args.command); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandRepositoryImpl_Update(t *testing.T) {
	type fields struct {
		store map[string][]command2.Command
	}
	type args struct {
		key      string
		commands []command2.Command
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		length  int
		wantErr bool
	}{
		{
			name:   "Test 01",
			fields: struct{ store map[string][]command2.Command }{store: map[string][]command2.Command{}},
			args: struct {
				key      string
				commands []command2.Command
			}{key: "key", commands: []command2.Command{
				command2.Set{Key: "key", Value: "value"},
				command2.Begin{},
			}},
			length:  2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CommandRepositoryImpl{
				store: tt.fields.store,
			}
			err := o.Update(tt.args.key, tt.args.commands)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.fields.store[tt.args.key]) != tt.length {
				t.Errorf("Update() length should be %v, got %v", tt.length, len(tt.fields.store[tt.args.key]))
			}
		})
	}
}

func TestCommandRepositoryImpl_Delete(t *testing.T) {
	store := map[string][]command2.Command{
		"key":  {command2.Set{Key: "key", Value: "value"}, command2.Begin{}, command2.Set{Key: "key", Value: "value2"}, command2.Commit{}},
		"key2": {command2.Begin{}, command2.Set{Key: "key2", Value: "value2"}, command2.Commit{}}}
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "Test 01",
			key:     "key",
			wantErr: false,
		}, {
			name:    "Test 02",
			key:     "key",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CommandRepositoryImpl{
				store: store,
			}
			if err := o.Delete(tt.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandRepositoryImpl_GetCommands(t *testing.T) {
	store := map[string][]command2.Command{
		"key":  {command2.Set{Key: "key", Value: "value"}, command2.Begin{}, command2.Set{Key: "key", Value: "value2"}, command2.Commit{}},
		"key2": {command2.Begin{}, command2.Set{Key: "key2", Value: "value2"}, command2.Commit{}}}

	tests := []struct {
		name         string
		key          string
		wantCommands []command2.Command
		wantErr      bool
	}{
		{
			name:         "Test 01",
			key:          "key",
			wantCommands: []command2.Command{command2.Set{Key: "key", Value: "value"}, command2.Begin{}, command2.Set{Key: "key", Value: "value2"}, command2.Commit{}},
			wantErr:      false,
		}, {
			name:         "Test 02",
			key:          "key8",
			wantCommands: []command2.Command{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CommandRepositoryImpl{
				store: store,
			}
			gotCommands, err := o.GetCommands(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (len(gotCommands) != 0 || len(tt.wantCommands) != 0) && !reflect.DeepEqual(gotCommands, tt.wantCommands) {
				t.Errorf("GetCommands() gotCommands = %v, want %v", gotCommands, tt.wantCommands)
			}
		})
	}
}

func TestCommandRepositoryImpl_GetKeys(t *testing.T) {
	store1 := map[string][]command2.Command{
		"key":  {command2.Set{Key: "key", Value: "value2"}, command2.Commit{}},
		"key2": {command2.Begin{}}}
	store2 := map[string][]command2.Command{
		"key":  {command2.Set{Key: "key", Value: "value"}},
		"key2": {command2.Set{Key: "key", Value: "value"}},
		"key3": {command2.Begin{}, command2.Set{Key: "key2", Value: "value2"}, command2.Commit{}}}
	tests := []struct {
		name     string
		fields   map[string][]command2.Command
		wantKeys []string
	}{
		{name: "Test 01", fields: store1, wantKeys: []string{"key", "key2"}},
		{name: "Test 02", fields: store2, wantKeys: []string{"key", "key2", "key3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CommandRepositoryImpl{
				store: tt.fields,
			}
			if gotKeys := o.GetKeys(); !reflect.DeepEqual(gotKeys, tt.wantKeys) {
				t.Errorf("GetKeys() = %v, want %v", gotKeys, tt.wantKeys)
			}
		})
	}
}

func TestCommandRepositoryImpl_Count(t *testing.T) {
	store := map[string][]command2.Command{
		"foo": {command2.Set{Key: "foo", Value: "123"}, command2.Begin{}, command2.Set{Key: "foo", Value: "456"}, command2.Commit{}},
		"bar": {command2.Begin{}, command2.Set{Key: "bar", Value: "456"}, command2.Commit{}}}
	tests := []struct {
		name       string
		value      string
		wantNumber int
	}{
		{name: "Test 01", value: "456", wantNumber: 2},
		{name: "Test 02", value: "123", wantNumber: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CommandRepositoryImpl{
				store: store,
			}
			if gotNumber := o.Count(tt.value); gotNumber != tt.wantNumber {
				t.Errorf("Count() = %v, want %v", gotNumber, tt.wantNumber)
			}
		})
	}
}
