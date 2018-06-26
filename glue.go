package glue

import (
	"bufio"
	"encoding/json"
	"io"
	"sync/atomic"
)

// Request .
type Request struct {
	ID      int64  `json:"id,omitempty"`
	Payload []byte `json:"payload,omitempty"`
}

// Processor .
type Processor func(req *Request) *Request

type plug struct {
	stop          chan struct{}
	next          int64
	inRMain, outR io.Reader
	inW, outWMain io.Writer
	proc          Processor
}

func (p *plug) agent() {
	reader := p.inRMain
	writer := p.outWMain
	proc := p.proc
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		js := scanner.Bytes()
		var req Request
		json.Unmarshal(js, &req)

		res := proc(&req)
		if res == nil {
			continue
		}

		js, _ = json.Marshal(res)
		js = append(js, '\n')
		writer.Write(js)
	}
}

func (p *plug) Request(payload []byte) int64 {
	id := atomic.AddInt64(&p.next, 1)
	var req Request
	req.ID = id
	req.Payload = payload
	js, _ := json.Marshal(&req)
	js = append(js, '\n')
	p.outWMain.Write(js)
	return id
}

// StartServer reader and writer should be os.Stdin and os.Stdout,
func StartServer(reader io.Reader, writer io.Writer, proc Processor) {
	res := &plug{
		inRMain:  reader,
		outWMain: writer,
		proc:     proc,
	}
	go res.agent()
}

// StartClient the result values should be used as:
//	cmd.Stdout = writer
//	cmd.Stdin = reader
func StartClient(proc Processor) (reader io.Reader, writer io.Writer, sender interface{ Request(payload []byte) int64 }) {
	res := &plug{
		proc: proc,
	}
	res.inRMain, res.inW = io.Pipe()
	res.outR, res.outWMain = io.Pipe()
	go res.agent()
	return res.outR, res.inW, res
}
