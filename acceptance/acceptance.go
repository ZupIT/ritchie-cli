package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

type Inputs struct {
	args   []string
	prompt []string
}

func commandInit(cmd *exec.Cmd) (stdin io.WriteCloser, err error, out io.Reader) {
	stdin, err = cmd.StdinPipe()
	if err != nil {
		return nil, err, nil
	}

	stdout, _ := cmd.StdoutPipe()

	err = cmd.Start()
	if err != nil {
		return nil, err, nil

	}

	return stdin, nil, stdout
}

func (i *Inputs) RunRit() (string, error) {
	args := i.args
	cmd := exec.Command("rit", args...)

	stdin, err, out := commandInit(cmd)

	if err != nil {
		log.Panic(err)
	}

	for _, s := range i.prompt {
		time.Sleep(1000 * time.Millisecond)
		_, err = io.Copy(stdin, bytes.NewBuffer([]byte(s)))
		if err != nil {
			log.Printf("Error when giving inputs: %q", err)
			os.Exit(1)
		}
	}

	defer stdin.Close()

	resp := ""
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		resp = fmt.Sprint(resp, m, "\n")
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("Error while running: %q", err)
		os.Exit(1)
	}

	return resp, nil
}
