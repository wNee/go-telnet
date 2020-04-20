package telnet

import (
	"bytes"
	"fmt"
	"github.com/reiver/go-oi"
	"io"
	"time"
)

// StandardCaller is a simple TELNET client which sends to the server any data it gets from os.Stdin
// as TELNET (and TELNETS) data, and writes any TELNET (or TELNETS) data it receives from
// the server to os.Stdout, and writes any error it has to os.Stderr.

type CmdCaller struct {
	cmds []string
	out  string
	err  error
	loginDelaySecond int
	cmdDelaySecond   int
}

func (caller *CmdCaller) SetLoginDelay(x int) *CmdCaller {
	caller.loginDelaySecond = x
	return caller
}

func (caller *CmdCaller) SetCmdDelay(x int) *CmdCaller {
	caller.cmdDelaySecond = x
	return caller
}

func (caller *CmdCaller) GetStdOut() string {
	return caller.out
}

func (caller *CmdCaller) GetStdError() string {
	return caller.out
}

func (caller *CmdCaller) AddCmds(cmds []string) {
	caller.cmds = append(caller.cmds, cmds...)
}

func (caller *CmdCaller) CallTELNET(ctx Context, w Writer, r Reader) {
	caller.cmdCallerCallTELNET(ctx, w, r)
}

func (caller *CmdCaller) cmdCallerCallTELNET(ctx Context, w Writer, r Reader) {
	writeBuf := bytes.NewBuffer(nil)
	go func(writer io.Writer, reader io.Reader) {
		var readBuffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
		readP := readBuffer[:]

		for {
			// Read 1 byte.
			n, err := reader.Read(readP)
			if n <= 0 && nil == err {
				continue
			} else if n <= 0 && nil != err {
				break
			}
			oi.LongWrite(writer, readP)
		}
	}(writeBuf, r)

	var buffer bytes.Buffer
	var p []byte

	var crlfBuffer [2]byte = [2]byte{'\r','\n'}
	crlf := crlfBuffer[:]

	for i, cmd := range caller.cmds {
		delay := caller.loginDelaySecond
		if delay == 0 {
			delay = 1
		}
		time.Sleep(time.Duration(delay) * time.Second)
		buffer.Write([]byte(cmd))
		buffer.Write(crlf)

		p = buffer.Bytes()

		n, err := oi.LongWrite(w, p)
		if nil != err {
			break
		}
		if expected, actual := int64(len(p)), n; expected != actual {
			err := fmt.Errorf("Transmission problem: tried sending %d bytes, but actually only sent %d bytes.", expected, actual)
			caller.err = err
			fmt.Println(err.Error())
			return
		}

		buffer.Reset()

		if i == len(caller.cmds) - 1 {
			cmdDelay := caller.cmdDelaySecond
			if cmdDelay == 0 {
				cmdDelay = 30
			}
			time.Sleep(time.Duration(cmdDelay) * time.Second)
		}
	}

	// Wait a bit to receive data from the server (that we would send to io.Stdout).
	time.Sleep(3 * time.Millisecond)
	fmt.Println(writeBuf.String())
	caller.out = writeBuf.String()
}


