package main

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type TaskType string

const (
	CreateDir  TaskType = "create_dir"
	CreateFile          = "create_file"
	RmFile              = "rm_file"
	RmDir               = "rm_dir"
	PutContent          = "put_content"
)

type Task struct {
	Name        string            `yaml:"name"`
	Type        TaskType          `yaml:"type"`        // required
	AbortOnFail bool              `yaml:"abortOnFail"` // default to false
	Args        map[string]string `yaml:"args"`
}

func parseTasks(in []byte) (tasks []Task, err error) {
	if err = yaml.Unmarshal(in, &tasks); err != nil {
		return
	}

	for _, task := range tasks {
		switch task.Type {
		case CreateDir, CreateFile, RmDir, RmFile, PutContent: // ok
		default:
			err = errors.New("invalid type")
			return
		}
	}

	return
}

func (task *Task) run() (err error) {
	path, hasPath := task.Args["path"]

	switch task.Type {
	case CreateDir:
		// args: path
		if !hasPath {
			return errors.New("create_dir without path")
		}
		err = os.Mkdir(path, 0750)

	case CreateFile:
		// args: path
		if !hasPath {
			return errors.New("create_file without path")
		}
		var file *os.File
		file, err = os.Create(path)
		if err != nil {
			return
		}
		file.Close()

	case RmDir:
		// args: path, recursive
		if !hasPath {
			return errors.New("rm_dir without path")
		}
		recursive, exists := task.Args["recursive"]
		if !exists {
			recursive = "false"
		}
		if recursive == "true" {
			err = os.RemoveAll(path)
		} else if recursive == "false" {
			err = os.Remove(path)
		} else {
			return errors.New(`rm_dir.recursive is neither "true" or "false"`)
		}

	case RmFile:
		// args: path
		if !hasPath {
			return errors.New("rm_file without path")
		}
		err = os.Remove(path)

	case PutContent:
		// args: path, content, append
		if !hasPath {
			return errors.New("put_content without path")
		}

		appendContent, exists := task.Args["append"]
		if !exists {
			appendContent = "false"
		}
		var flag int
		if appendContent == "true" {
			flag = os.O_APPEND
		} else if appendContent == "false" {
			flag = os.O_WRONLY
		} else {
			return errors.New(`put_content.append is neither "true" or "false"`)
		}
		var file *os.File
		file, err = os.OpenFile(path, flag, 0755)
		if err != nil {
			return
		}
		defer file.Close()

		_, err = file.WriteString(task.Args["content"])

	default:
		// skip
	}

	return
}

func runTasks(tasks []Task, handle func(i int, err error)) (err error) {
	for i, task := range tasks {
		if err = task.run(); err != nil {
			if task.AbortOnFail {
				handle(i, err)
				return
			} else {
				handle(i, err)
			}
		}
	}
	return
}
