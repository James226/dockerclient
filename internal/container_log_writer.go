package internal

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/pkg/stdcopy"
)

// ContainerLogWriter is an implementation of io.Writer used to capture and
// write a container's stdout and stderr to a specified io.Writer.
type ContainerLogWriter struct {
	name string
	dst  io.Writer
}

// NewContainerLogWriter is used to return a new instance of ContainerLogWriter,
// with a container's name and an output destination.
func NewContainerLogWriter(name string, dst io.Writer) *ContainerLogWriter {
	return &ContainerLogWriter{
		name: name,
		dst:  dst,
	}
}

// Write is used to writes the given data to the underlying, destination io.Writer.
// The data will be split by each line, then outputed with the container's name as
// a prefix.
//
// Note: the returned integer value will match the length of the data provided.
// This is due to an issue with the Docker stdcopy package, which expects the output length
// to match exactly the length of the data provided; not allowing any greater.
// https://github.com/moby/moby/issues/45377
func (w *ContainerLogWriter) Write(data []byte) (int, error) {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		_, err := w.dst.Write([]byte(fmt.Sprintf("[%s]: %s\n", w.name, []byte(line))))
		if err != nil {
			return 0, err
		}
	}
	return len(data), nil
}

// PrintContainerLogs is used to print the output of a container to the stdout/err streams.
func PrintContainerLogs(name string, out io.Reader) {
	_, err := stdcopy.StdCopy(
		NewContainerLogWriter(name, os.Stdout),
		NewContainerLogWriter(name, os.Stderr),
		out,
	)
	if err != nil {
		fmt.Printf("Failed to print the logs of container '%s'\n", name)
	}
}
