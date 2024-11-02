package exec

import (
	"bufio"
	_ "embed"
	"io"
	"os/exec"
	"sync"
)

type Data struct {
	IsStderr bool
	Value    string
	EOF      bool
}
type Cmd struct {
	ch                     chan Data
	cmd                    *exec.Cmd
	scanner, scannerStderr *bufio.Scanner
	stdin                  io.WriteCloser
}

func NewCmd(name string, arg ...string) (*Cmd, error) {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Cmd{
		ch:            make(chan Data),
		cmd:           cmd,
		scanner:       bufio.NewScanner(stdout),
		scannerStderr: bufio.NewScanner(stderr),
		stdin:         stdin,
	}, nil
}

func (command *Cmd) Run() (int, error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for command.scannerStderr.Scan() {
			command.ch <- Data{Value: command.scannerStderr.Text(), IsStderr: true}
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		for command.scanner.Scan() {
			command.ch <- Data{Value: command.scanner.Text()}
		}
		wg.Done()
	}()
	wg.Wait()
	command.ch <- Data{EOF: true}
	close(command.ch)

	if err := command.cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), nil
		}
		return 0, err
	}
	return 0, nil
}

func (command *Cmd) Output() <-chan Data {
	return command.ch
}

func (command *Cmd) Write(data string) {
	_, _ = command.stdin.Write([]byte(data))
}
