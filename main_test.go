package main

import (
	_ "embed"
	"fmt"
	"os"
	"testing"
)

//go:embed sample.yml
var sample string

func Test(t *testing.T) {

	tasks, err := parseTasks([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}

	os.RemoveAll("/tmp/project")

	var buf string
	err = runTasks(tasks, func(i int, err error) {
		buf += fmt.Sprintf("ERR: task %d, %v\n", i, err)
	})

	if buf != "" {
		t.Fatalf("buf is not empty: %s", buf)
	}

	if err != nil {
		t.Fatal(err)
	}
}
