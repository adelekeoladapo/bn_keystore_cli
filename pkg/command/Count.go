package command

type Count struct {
	Key string `json:"key"`
}

func (c Count) isCommand() {}
