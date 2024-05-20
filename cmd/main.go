package main

import (
	"fmt"
	"os"

	"github.com/himuraxkenji/mocker-go/internal/server"
)

func main() {
	fmt.Println("Variable - ", os.Getenv("HELLO"))
	server.InitServer()
}
