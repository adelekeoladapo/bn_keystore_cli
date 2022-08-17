package repository

import (
	"bn_keystore/pkg/command"
)

/*
*  This repository stores commands on each key
 */

type CommandRepository interface {
	Create(key string, command command.Command) error
	Update(key string, command []command.Command) error
	Delete(key string) error
	GetCommands(key string) ([]command.Command, error)
	GetKeys() []string
	Count(value string) int
}
