package command

type Delete struct {
	Key string `json:"key"`
}

func (c Delete) isCommand() {}
