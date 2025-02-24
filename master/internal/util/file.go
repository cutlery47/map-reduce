package util

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func SplitFile(file, chunkDir, resultDir string, parts int) ([]io.Reader, error) {
	// slice of file desctiptors abstracted as io.Readers
	readers := []io.Reader{}

	// creating chunk dir
	if err := ExecMkdirWithParent(chunkDir); err != nil {
		return nil, fmt.Errorf("execMkdirWithParent: %v", err)
	}

	// creating results dir
	if err := ExecMkdirWithParent(resultDir); err != nil {
		return nil, fmt.Errorf("execMkdirWithParent: %v", err)
	}

	// splitting file into chunks in chunk dir
	if err := ExecSplitFile(file, chunkDir, parts); err != nil {
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

func ExecMkdirWithParent(dir string) error {
	bash := "mkdir"
	arg0, arg1 := "-p", dir

	return execCmd(bash, arg0, arg1)
}

func ExecSplitFile(file, dir string, parts int) error {
	bash := "split"
	arg0, arg1 := "-n", strconv.Itoa(parts)
	arg2 := "-d"
	arg3, arg4 := file, dir+"/"

	return execCmd(bash, arg0, arg1, arg2, arg3, arg4)
}

func ExecCreateAndWriteFile(name, dir string, data []byte) error {
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

func CreateNestedDirString(b strings.Builder, objs ...string) string {
	b.Reset()

	for _, obj := range objs {
		b.WriteString(obj)
		b.WriteRune('/')
	}

	return b.String()[:b.Len()-1]
}
