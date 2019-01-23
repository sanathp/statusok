package requests

import (
	oi "github.com/reiver/go-oi"
	telnet "github.com/reiver/go-telnet"
)

type caller struct{}

func (c caller) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {
	oi.LongWrite(w, []byte("\r"))
}

func DoTelnet(addr string) error {
	err := telnet.DialToAndCall(addr, caller{})
	if err != nil {
		return err
	}
	return nil
}
