package main

import (
	"os"

	"github.com/dc0d/glue"
	"github.com/dc0d/onexit"
)

func main() {
	glue.StartServer(os.Stdin, os.Stdout, func(req *glue.Request) *glue.Request {
		req.Payload = append([]byte("<<SERVICE 1 PROCESSED DATA>>:"), req.Payload...)
		return req
	})
	<-onexit.Done()
}
