package exec

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSleep1s(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./test-sleep1/sleep1s.go")
	if err != nil {
		t.Fatal(err)
	}

	msgs := run(t, cmd, 0)

	assert.Equal(t, 1, len(msgs))
	assert.True(t, msgs[0].EOF)
}

func TestPrintln(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./test-println/println.go")
	if err != nil {
		t.Fatal(err)
	}

	msgs := run(t, cmd, 0)

	assert.Equal(t, len(msgs), 3)
	assert.Equal(t, msgs[0].Value, "Hello World!")
	assert.False(t, msgs[0].IsStderr)
	assert.Equal(t, msgs[1].Value, "Hello Moon!")
	assert.False(t, msgs[1].IsStderr)
	assert.True(t, msgs[2].EOF)
}

func TestExit1(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./test-exit1/exit1.go")
	if err != nil {
		t.Fatal(err)
	}

	msgs := run(t, cmd, 1)

	assert.Equal(t, len(msgs), 2)
	assert.Equal(t, msgs[0].Value, "exit status 1")
	assert.True(t, msgs[0].IsStderr)
	assert.True(t, msgs[1].EOF)
}

func TestStderr(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./test-stderr/stderr.go")
	if err != nil {
		t.Fatal(err)
	}

	msgs := run(t, cmd, 0)

	assert.Equal(t, len(msgs), 3)
	assert.Equal(t, msgs[0].Value, "Hello Stderr!")
	assert.True(t, msgs[0].IsStderr)
	assert.Equal(t, msgs[1].Value, "Hello errors!")
	assert.True(t, msgs[1].IsStderr)
	assert.True(t, msgs[2].EOF)
}

func TestStdin(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./test-stdin/stdin.go")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		cmd.Write("example\n")
	}()

	msgs := run(t, cmd, 0)

	assert.Equal(t, 1, len(msgs))
	assert.True(t, msgs[0].EOF)
}

func TestStdinFail(t *testing.T) {
	cmd, err := NewCmd("go", "run", "./test-stdin/stdin.go")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		cmd.Write("bogus\n")
	}()

	msgs := run(t, cmd, 1)

	assert.Equal(t, len(msgs), 2)
	assert.Equal(t, msgs[0].Value, "exit status 1")
	assert.True(t, msgs[0].IsStderr)
	assert.True(t, msgs[1].EOF)
}

func TestCalculator(t *testing.T) {
	cmd, err := NewCmd("bc")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		cmd.Write("3 * 4\n")
		cmd.Write("2 + 2\n")
		cmd.Write("quit\n")
	}()

	msgs := run(t, cmd, 0)

	assert.Equal(t, len(msgs), 3)
	assert.Equal(t, msgs[0].Value, "12")
	assert.False(t, msgs[0].IsStderr)
	assert.Equal(t, msgs[1].Value, "4")
	assert.False(t, msgs[1].IsStderr)
	assert.True(t, msgs[2].EOF)
}

func run(t *testing.T, cmd *Cmd, _exitCode int) []Data {
	msgs := []Data{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for msg := range cmd.Output() {
			msgs = append(msgs, msg)
		}
		wg.Done()
	}()
	exitCode, err := cmd.Run()
	if err != nil {
		t.Fatalf("%s: cmd.Run failed with error %q", t.Name(), err)
	}
	if exitCode != _exitCode {
		t.Fatalf("%s: invalid exit code %d (expected %d)", t.Name(), exitCode, _exitCode)
	}
	wg.Wait()
	return msgs
}
