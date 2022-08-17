package repository

import (
	command2 "bn_keystore/pkg/command"
	"errors"
)

type CommandRepositoryImpl struct {
	store map[string][]command2.Command
}

func (o *CommandRepositoryImpl) Create(key string, command command2.Command) (err error) {
	o.store[key] = append(o.store[key], command)
	return
}

func (o *CommandRepositoryImpl) Update(key string, commands []command2.Command) (err error) {
	o.store[key] = commands
	return
}

func (o *CommandRepositoryImpl) Delete(key string) (err error) {
	if _, ok := o.store[key]; !ok {
		err = errors.New("key not set")
		return
	} else {
		delete(o.store, key)
	}
	return
}

func (o *CommandRepositoryImpl) GetCommands(key string) (commands []command2.Command, err error) {
	list, ok := o.store[key]
	if !ok {
		err = errors.New("key not set")
		return
	}
	commands = list
	return
}

func (o *CommandRepositoryImpl) GetKeys() (keys []string) {
	for key, _ := range o.store {
		keys = append(keys, key)
	}
	return
}

func (o *CommandRepositoryImpl) Count(value string) (number int) {
	for _, commands := range o.store {
		for _, cmd := range commands {
			if c, ok := cmd.(command2.Set); ok && c.Value == value {
				number++
			}
		}
	}
	return
}

func GetCommandRepositoryImpl() CommandRepository {
	return &CommandRepositoryImpl{
		store: map[string][]command2.Command{},
	}
}
