# go-exec

Simple wrapper package around Go standard package `os.exec`.

Example:

```go
func Example() {
    // use Linux calculator for this example:
	cmd, err := NewCmd("bc")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		exitCode, err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run failed with error %q", err)
		}
		if exitCode != 0 {
			log.Fatalf("invalid exit code %d ", exitCode)
		}
	}()

	cmd.Write("2 + 2\n")

	result := <-cmd.Output()
	fmt.Println(result.Value)

	cmd.Write("quit\n")
	// Output: 4
}
```