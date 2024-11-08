package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientRun(t *testing.T) {
	cmd := exec.Command("go", "run", "cmd/client/main.go", "--output-dir=thumbnails", "--async", "--server-address=localhost:50051", "test_video")

	err := cmd.Run()

	assert.NoError(t, err, "Ошибка при запуске клиентского приложения")
}
