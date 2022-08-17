package command

type Get struct {
	Key string `json:"key"`
}

func (c Get) isCommand() {}
