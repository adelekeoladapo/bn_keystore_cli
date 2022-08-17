package main

import (
	"bn_keystore/pkg/controller"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	ctrl := controller.GetCommandController()
	r := bufio.NewReader(os.Stdin)
	var s string
	for {
		fmt.Print("> ")
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "EXIT" {
			break
		}
		if s == "" {
			continue
		}
		if response, err := ctrl.Process(s); err != nil {
			fmt.Println(err)
		} else if len(response) > 0 {
			fmt.Println(response)
		}
	}
}
