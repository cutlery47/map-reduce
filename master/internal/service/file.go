package service

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

func splitFile(file, chunkDir, resultDir string, parts int) ([]io.Reader, error) {
	// slice of file desctiptors abstracted as io.Readers
	readers := []io.Reader{}

	// creating chunk dir
	if err := execMkdirWithParent(chunkDir); err != nil {
		return nil, fmt.Errorf("execMkdirWithParent: %v", err)
	}

	if err := execMkdirWithParent(resultDir); err != nil {
		return nil, fmt.Errorf("execMkdirWithParent: %v", err)
	}

	if err := execSplitFile(file, chunkDir, parts); err != nil {
		return nil, fmt.Errorf("execSplitFile: %v", err)
	}

	// collecting chunks
	chunks, err := os.ReadDir(os.ExpandEnv(chunkDir))
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %v", err)
	}

	for _, chunk := range chunks {
		path := fmt.Sprintf("%v/%v", chunkDir, chunk.Name())
		fd, err := os.Open(os.ExpandEnv(path))
		if err != nil {
			return nil, fmt.Errorf("os.Open: %v", err)
		}

		readers = append(readers, fd)
	}

	return readers, nil
}

func execMkdirWithParent(dir string) error {
	bash := "mkdir"
	arg0, arg1 := "-p", os.ExpandEnv(dir)

	return execCmd(bash, arg0, arg1)
}

func execSplitFile(file, dir string, parts int) error {
	bash := "split"
	arg0, arg1 := "-n", strconv.Itoa(parts)
	arg2 := "-d"
	arg3, arg4 := os.ExpandEnv(file), os.ExpandEnv(dir+"/")

	return execCmd(bash, arg0, arg1, arg2, arg3, arg4)
}

func execCreateAndWriteFile(name, dir string, data []byte) error {
	expDir := os.ExpandEnv(dir)
	return os.WriteFile(fmt.Sprintf("%v/%v", expDir, name), data, 0777)
}

func execCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd.Run: %v", err)
	}

	return nil
}
