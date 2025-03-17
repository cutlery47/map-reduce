package core

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	mr "github.com/cutlery47/map-reduce/mapreduce"
)

// file stuff handler
type FileHandler struct {
	chunkPath  string // path of the chunk directory
	resultPath string // path of the results directory
	conf       mr.Config
}

func NewFileHandler(conf mr.Config) *FileHandler {
	return &FileHandler{
		chunkPath:  createNestedPaths(conf.MappedDir, "chunks"),
		resultPath: createNestedPaths(conf.MappedDir, "result"),
		conf:       conf,
	}
}

// splits provided file into parts-amount of chunks
// returns io.Reader instead tho
func (fh *FileHandler) Split(parts int) ([]io.Reader, error) {
	// splitting file into chunks in chunk dir
	if err := execSplitFile(fh.conf.File, fh.chunkPath, parts); err != nil {
		return nil, fmt.Errorf("execSplitFile: %v", err)
	}

	entries, err := os.ReadDir(fh.chunkPath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %v", err)
	}

	var (
		chunks = make([]io.Reader, 0, len(entries))
	)

	// getting fd for each chunk and storing them in result slice
	for _, entry := range entries {
		fd, err := os.Open(fmt.Sprintf("%v/%v", fh.chunkPath, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("os.Open: %v", err)
		}

		chunks = append(chunks, fd)
	}

	return chunks, nil
}

// creates chunk and result directories
func (fh *FileHandler) CreateDirs() error {
	if err := execMkdirWithParent(fh.chunkPath); err != nil {
		return err
	}

	if err := execMkdirWithParent(fh.resultPath); err != nil {
		return err
	}

	return nil
}

func (fh *FileHandler) CreateResult(name string, data []byte) error {
	var (
		path = fmt.Sprintf("%v/%v", fh.resultPath, name)
	)

	return os.WriteFile(path, data, 0777)
}

// creates a nested file path for given objects
func createNestedPaths(objs ...string) string {
	var b strings.Builder

	for _, obj := range objs {
		b.WriteString(obj)
		b.WriteRune('/')
	}

	return b.String()[:b.Len()-1]
}

func execMkdirWithParent(path string) error {
	bash := "mkdir"
	arg0, arg1 := "-p", path

	return execCmd(bash, arg0, arg1)
}

func execSplitFile(src, dst string, parts int) error {
	bash := "split"
	arg0, arg1 := "-n", strconv.Itoa(parts)
	arg2 := "-d"
	arg3, arg4 := src, dst+"/"

	return execCmd(bash, arg0, arg1, arg2, arg3, arg4)
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
