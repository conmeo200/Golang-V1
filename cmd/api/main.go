package main

import (
	"fmt"

	"github.com/conmeo200/Golang-V1/internal/bootstrap"
)

func main() {
	fmt.Println("Init Api......")

	err := bootstrap.RunAPI()

	if err != nil {
		fmt.Println("Error running api: ", err)
	}

	fmt.Println("Server Api Running ......")
}
