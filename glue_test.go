package glue

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSmoke(t *testing.T) {
	require := require.New(t)

	rcvd := make(chan struct{})
	start := make(chan struct{})
	var id int64
	commandStdin, commandStdout, sender := StartClient(func(req *Request) *Request {
		<-start
		defer close(rcvd)
		got := string(req.Payload)
		require.Equal("<<PROC>>:DATA", got)
		require.Equal(int64(1), atomic.LoadInt64(&id))
		return nil
	})

	StartServer(commandStdin, commandStdout, func(req *Request) *Request {
		got := string(req.Payload)
		require.Equal("DATA", got)
		req.Payload = append([]byte("<<PROC>>:"), req.Payload...)
		return req
	})

	samplePayload := "DATA"
	expectedID := sender.Request([]byte(samplePayload))
	atomic.StoreInt64(&id, expectedID)
	close(start)
	<-rcvd
}
