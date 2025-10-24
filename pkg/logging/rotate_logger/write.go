package rotate_logger

func (rl *RotateLogger) Write(p []byte) (n int, err error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return rl.file.Write(p)
}