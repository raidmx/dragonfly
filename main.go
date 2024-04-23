package main

import "github.com/STCraft/dragonfly/server"

func main() {
	srv, _ := server.New()
	srv.Start()
}
