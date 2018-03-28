package preconn

import (
	"bytes"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	head = "hello "
	text = "world"
	full = "hello world"
)

func TestPreConn(t *testing.T) {
	l, err := net.Listen("tcp", ":0")
	if !assert.NoError(t, err) {
		return
	}
	defer l.Close()
	go func() {
		for {
			conn, err := l.Accept()
			if err == nil {
				conn.Write([]byte(text))
				conn.Close()
			}
		}
	}()

	conn, err := net.Dial("tcp", l.Addr().String())
	if !assert.NoError(t, err) {
		return
	}
	defer conn.Close()
	pconn := Wrap(conn, []byte(head))
	var buf bytes.Buffer
	b := make([]byte, 2)
	for {
		n, err := pconn.Read(b)
		if err == io.EOF {
			break
		}
		if !assert.NoError(t, err) {
			return
		}
		buf.Write(b[:n])
	}
	assert.Equal(t, full, buf.String(), "Read() multiple times should get the full data")

	conn, err = net.Dial("tcp", l.Addr().String())
	if !assert.NoError(t, err) {
		return
	}
	defer conn.Close()
	b = make([]byte, len(full))
	pconn = Wrap(conn, []byte(head))
	n, err := pconn.Read(b)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, n, len(b), "Read() should return as many data as possible")
	assert.Equal(t, full, string(b))
}
