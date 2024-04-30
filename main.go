package main

import "github.com/stcraft/dragonfly/server"

func main() {
	srv, _ := server.New()
	srv.Start()
}
