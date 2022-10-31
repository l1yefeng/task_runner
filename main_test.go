package main

import (
	"fmt"
	"os"
	"testing"
)

func Test(t *testing.T) {
	in := `
- name: Create root directory
  type: create_dir
  abortOnFail: true
  args:
    path: /tmp/project
- name: Create VERSION file
  type: create_file
  args:
    path: /tmp/project/VERSION
- name: Set VERSION
  type: put_content
  args:
    path: /tmp/project/VERSION
    content: 1.0.0
    append: false # overwrite the file
# Here we could do other operations, but we don't have Type
# for them, so we do nothing.
- name: Clean up
  type: rm_dir
  abortOnFail: true
  args:
    path: /tmp/project
    recursive: true
`
	tasks, err := parseTasks([]byte(in))
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
