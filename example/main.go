package main

import (
	telnet "github.com/reiver/go-telnet"
)

func main() {
	c := &telnet.CmdCaller{}
	c.AddCmds([]string{"root", "", "touch b.txt"})

	c.SetCmdDelay(2).SetLoginDelay(1)
	telnet.DialToAndCall("10.169.231.19:23", telnet.Caller(c))

	//var caller telnet.Caller = telnet.StandardCaller
	//telnet.DialToAndCall("10.169.231.19:23", caller)
}
