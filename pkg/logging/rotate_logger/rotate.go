package rotate_logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

func (rl *RotateLogger) rotate() {
	rl.mu.Lock()

	newFileName := fmt.Sprintf("%s/%s.log", rl.dir, time.Now().Format(time.DateOnly))

	newFile, err := os.OpenFile(newFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		rl.mu.Unlock()
		log.Printf("произошла ошибка при создании нового файла для логов: %s", err.Error())
		return
	}

	rl.file.Close()

	rl.file = newFile

	rl.mu.Unlock()
}
