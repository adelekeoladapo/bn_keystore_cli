package command

type Set struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (c Set) isCommand() {}
