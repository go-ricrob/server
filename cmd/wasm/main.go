//go:build wasm

package main

import (
	"fmt"

	"github.com/go-ricrob/server/internal/webui"
)

func main() {
	fmt.Println("Hello from the ricrob go web assembly")
	webui.New()
	done := make(chan struct{})
	<-done
}
