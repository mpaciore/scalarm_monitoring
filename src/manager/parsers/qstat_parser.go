package parsers

import (
	"bytes"
	"os/exec"
	"errors"
)

func SplitOutput(c exec.Cmd) ([]byte, []byte, error) {
	if c.Stdout != nil {
 		return nil, nil, errors.New("exec: Stdout already set")
	}
	if c.Stderr != nil {
		return nil, nil, errors.New("exec: Stderr already set")
 	}
	var o bytes.Buffer
	var e bytes.Buffer
	c.Stdout = &o
 	c.Stderr = &e
 	err := c.Run()
 	return o.Bytes(), e.Bytes(), err
 }