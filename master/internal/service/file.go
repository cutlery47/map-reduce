package service

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func splitFile(file, chunkDir, resultDir string, parts int) ([]io.Reader, error) {
	// slice of file desctiptors abstracted as io.Readers
	readers := []io.Reader{}

	// creating chunk dir
	if err := execMkdirWithParent(chunkDir); err != nil {
		return nil, fmt.Errorf("execMkdirWithParent: %v", err)
	}

	// creating results dir
	if err := execMkdirWithParent(resultDir); err != nil {
		return nil, fmt.Errorf("execMkdirWithParent: %v", err)
	}

	// splitting file into chunks in chunk dir
	if err := execSplitFile(file, chunkDir, parts); err != nil {
		return nil, fmt.Errorf("execSplitFile: %v", err)
	}

	// collecting chunks
	chunks, err := os.ReadDir(chunkDir)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %v", err)
	}

	// getting each chunk's fd
	for _, chunk := range chunks {
		path := fmt.Sprintf("%v/%v", chunkDir, chunk.Name())
		fd, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("os.Open: %v", err)
		}

		readers = append(readers, fd)
	}

	return readers, nil
}

func execMkdirWithParent(dir string) error {
	bash := "mkdir"
	arg0, arg1 := "-p", dir

	return execCmd(bash, arg0, arg1)
}

func execSplitFile(file, dir string, parts int) error {
	bash := "split"
	arg0, arg1 := "-n", strconv.Itoa(parts)
	arg2 := "-d"
	arg3, arg4 := file, dir+"/"

	return execCmd(bash, arg0, arg1, arg2, arg3, arg4)
}

func execCreateAndWriteFile(name, dir string, data []byte) error {
	return os.WriteFile(fmt.Sprintf("%v/%v", dir, name), data, 0777)
}

func execCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd.Run: %v", err)
	}

	return nil
}

func createNestedDirString(b strings.Builder, objs ...string) string {
	b.Reset()

	for _, obj := range objs {
		b.WriteString(obj)
		b.WriteRune('/')
	}

	return b.String()[:b.Len()-1]
}
