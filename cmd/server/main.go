package main

import (
	"github.com/Sarek-Wyona/proglog/internal/server"
	"log"
)

func main()  {
	srv := server.NewHTTPServer(":8080")
	log.Fatalln(srv.ListenAndServe())

	// TODO Add a server stop and add comments
}
