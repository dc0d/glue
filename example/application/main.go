package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/dc0d/glue"

	"github.com/dc0d/onexit"
)

func main() {
	go app()
	<-onexit.Done()
}

func app() {
	responses := make(chan *glue.Request, 3)
	go func() {
		r, w, sender := glue.StartClient(func(req *glue.Request) *glue.Request {
			responses <- req
			return nil
		})
		cmd := exec.Command("../service1/service1")
		cmd.Stdout = w
		cmd.Stdin = r
		onexit.Register(func() { cmd.Process.Kill() })
		go func() {
			sender.Request([]byte("DATA"))
		}()
		cmd.Run()
	}()
	go func() {
		r, w, sender := glue.StartClient(func(req *glue.Request) *glue.Request {
			responses <- req
			return nil
		})
		cmd := exec.Command("../service2/service2")
		cmd.Stdout = w
		cmd.Stdin = r
		onexit.Register(func() { cmd.Process.Kill() })
		go func() {
			sender.Request([]byte("DATA"))
		}()
		cmd.Run()
	}()

	rcvd := make(map[string]int64)
OUT:
	for {
		select {
		case r := <-responses:
			rcvd[string(r.Payload)] = r.ID
		case <-time.After(time.Millisecond * 300):
			break OUT
		}
	}

	hits := 0
	for k, v := range rcvd {
		if k == "<<SERVICE 1 PROCESSED DATA>>:DATA" && v == 1 {
			hits++
		}
		if k == "<<SERVICE 2 PROCESSED DATA>>:DATA" && v == 1 {
			hits++
		}
	}
	if hits != 2 {
		panic("not what expected")
	}
	fmt.Println("ran successfully! :)")
	onexit.ForceExit(0)
}
