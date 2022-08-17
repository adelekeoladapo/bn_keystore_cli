package service

import (
	command2 "bn_keystore/pkg/command"
	repository2 "bn_keystore/pkg/repository"
	"errors"
)

type CommandServiceImpl struct {
	commands          []command2.Command
	transactionsDepth int
	repository        repository2.CommandRepository
}

func (o *CommandServiceImpl) Set(key, value string) (err error) {
	o.commands, err = o.repository.GetCommands(key)
	// if key does not exist or there's no command for the key or the last command is not a SET, create a SET command
	if err != nil || len(o.commands) == 0 || !o.isMostRecentCommandASet() {
		// if key is empty or there\'re no commands for key, check if there\'re running transactions as well as the depth, then create BEGIN commands
		if err != nil || len(o.commands) == 0 {
			for i := 0; i < o.transactionsDepth; i++ {
				if err := o.repository.Create(key, command2.Begin{}); err != nil {
					//log.Printf("error occurred, %s", err)
				}
			}
		}
		if err = o.repository.Create(key, command2.Set{
			Key:   key,
			Value: value,
		}); err != nil {
			return
		}
	}
	// if the most recent command is a SET, update it\'s value
	if o.isMostRecentCommandASet() {
		o.commands[len(o.commands)-1] = command2.Set{
			Key:   key,
			Value: value,
		}
		if err = o.repository.Update(key, o.commands); err != nil {
			return
		}
	}
	return
}

func (o *CommandServiceImpl) Get(key string) (value string, err error) {
	if o.commands, err = o.repository.GetCommands(key); err != nil {
		return
	}
	for i := len(o.commands) - 1; i >= 0; i-- {
		if cmd, ok := o.commands[i].(command2.Set); ok {
			value = cmd.Value
			return
		}
	}
	err = errors.New("key not set")
	return
}

func (o *CommandServiceImpl) Delete(key string) (err error) {
	if err = o.repository.Delete(key); err == nil {
		o.commands = []command2.Command{}
	}
	return
}

func (o *CommandServiceImpl) Count(key string) int {
	return o.repository.Count(key)
}

func (o *CommandServiceImpl) Begin() {
	for _, key := range o.repository.GetKeys() {
		if err := o.repository.Create(key, command2.Begin{}); err != nil {
			//log.Printf("error occurred, %s", err)
		}
	}
	o.transactionsDepth++
}

func (o *CommandServiceImpl) Commit() (err error) {
	// ensure there\'re running transactions before committing
	if o.transactionsDepth == 0 {
		err = errors.New("no transaction")
		return err
	}
	for _, key := range o.repository.GetKeys() {
		if err := o.repository.Create(key, command2.Commit{}); err != nil {
			//log.Printf("error occurred, %s", err)
		}
	}
	o.transactionsDepth--
	return
}

func (o *CommandServiceImpl) Rollback() (err error) {
	/* check if there's an uncommitted transaction
	   ensure there\'re running transactions before rolling back */
	if o.transactionsDepth == 0 {
		//log.Printf("can't rollback, no uncommitted transaction")
		err = errors.New("no transaction")
		return
	}
	for _, key := range o.repository.GetKeys() {
		commands, _ := o.repository.GetCommands(key)
		for i := len(commands) - 1; i >= 0; i-- {
			if _, ok := commands[i].(command2.Begin); ok {
				if e := o.repository.Update(key, commands[:i]); e != nil {
					//log.Printf("an error occurred, %s", e)
				}
				break
			}
		}
	}
	o.transactionsDepth--
	return
}

func (o *CommandServiceImpl) isMostRecentCommandASet() bool {
	if len(o.commands) == 0 {
		return false
	}
	_, ok := o.commands[len(o.commands)-1].(command2.Set)
	return ok
}

func GetCommandServiceImpl() CommandService {
	return &CommandServiceImpl{
		commands:          []command2.Command{},
		transactionsDepth: 0,
		repository:        repository2.GetCommandRepositoryImpl(),
	}
}
