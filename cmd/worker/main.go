package main

import (
	"fmt"

	"github.com/conmeo200/Golang-V1/internal/bootstrap"
)

func main() {
	fmt.Println("Init Worker......")

	err := bootstrap.RunAPI()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server Worker Running ......")
}
