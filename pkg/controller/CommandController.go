package controller

import (
	command2 "bn_keystore/pkg/command"
	service2 "bn_keystore/pkg/service"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CommandController struct {
	service service2.CommandService
}

func (o *CommandController) Process(str string) (response string, err error) {
	var cmd command2.Command
	cmd, err = o.validate(str)
	if err != nil {
		//log.Printf("an error occurred. %s\n", err)
		response = err.Error()
		return
	}
	switch c := cmd.(type) {
	case command2.Set:
		err = o.service.Set(c.Key, c.Value)
	case command2.Get:
		response, err = o.service.Get(c.Key)
	case command2.Delete:
		err = o.service.Delete(c.Key)
	case command2.Count:
		response = strconv.Itoa(o.service.Count(c.Key))
	case command2.Begin:
		o.service.Begin()
	case command2.Commit:
		err = o.service.Commit()
	case command2.Rollback:
		err = o.service.Rollback()
	}
	return
}

func (o *CommandController) validate(str string) (cmd command2.Command, err error) {
	segments := strings.Split(str, " ")
	for i, segment := range segments {
		segments[i] = strings.TrimSpace(segment)
	}
	if len(segments) > 3 {
		err = o.returnErrorStr(str)
		return
	}
	if len(segments) == 3 {
		if segments[0] == "SET" {
			cmd = command2.Set{
				Key:   segments[1],
				Value: segments[2],
			}
		} else {
			err = o.returnErrorStr(str)
		}
	} else if len(segments) == 2 {
		if segments[0] == "GET" {
			cmd = command2.Get{Key: segments[1]}
		} else if segments[0] == "COUNT" {
			cmd = command2.Count{Key: segments[1]}
		} else if segments[0] == "DELETE" {
			cmd = command2.Delete{Key: segments[1]}
		} else {
			err = o.returnErrorStr(str)
		}
	} else if len(segments) == 1 {
		if segments[0] == "BEGIN" {
			cmd = command2.Begin{}
		} else if segments[0] == "COMMIT" {
			cmd = command2.Commit{}
		} else if segments[0] == "ROLLBACK" {
			cmd = command2.Rollback{}
		} else {
			err = o.returnErrorStr(str)
		}
	}
	return
}

func (o *CommandController) returnErrorStr(error string) error {
	//log.Printf("Invalid command, %s", error)
	return errors.New(fmt.Sprintf("invalid command, %s", error))
}

func GetCommandController() CommandController {
	return CommandController{service: service2.GetCommandServiceImpl()}
}
