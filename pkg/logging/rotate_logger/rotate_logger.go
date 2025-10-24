package rotate_logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type RotateLogger struct {
	dir  string
	file *os.File
	mu   sync.Mutex
}

func NewRotateLogger(dir string) (*RotateLogger, error) {
	if err := os.MkdirAll(dir, 0644); err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("%s/%s.log", dir, time.Now().Format(time.DateOnly))

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &RotateLogger{
		dir:  dir,
		file: file,
		mu:   sync.Mutex{},
	}, nil
}
