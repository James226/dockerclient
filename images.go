package dockerclient

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/james226/dockerclient/options"
)

type Image struct {
	Name string
}

type ImageOperations struct {
	cli *client.Client
}

func (i ImageOperations) Pull(ctx context.Context, name string) (*Image, error) {
	reader, err := i.cli.ImagePull(ctx, name, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}

	defer reader.Close()
	_, err = io.Copy(os.Stdout, reader)
	return &Image{name}, err
}

func (i ImageOperations) Build(ctx context.Context, name string, path string, opts ...*options.BuildImageOptions) (*Image, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	err := loadBuildContext(path, "", tw)
	if err != nil {
		return nil, err
	}
	contextReader := bytes.NewReader(buf.Bytes())
	opt := options.BuildImage()
	if len(opts) > 0 {
		opt = opts[0]
	}
	buildOptions := types.ImageBuildOptions{
		Context:    bytes.NewReader(buf.Bytes()),
		Dockerfile: opt.Dockerfile(),
		Tags:       []string{name},
		Remove:     true,
	}
	platform, ok := opt.Platform()
	if ok {
		buildOptions.Platform = platform
	}
	build, err := i.cli.ImageBuild(ctx, contextReader, buildOptions)
	if err != nil {
		return nil, err
	}
	defer build.Body.Close()
	err = logBuildOutput(build.Body, name)
	if err != nil {
		return nil, err
	}
	return &Image{Name: name}, nil
}

// Used to recursively load files from the specified path a .tar file.
func loadBuildContext(path, relativePath string, tw *tar.Writer) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read dir: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			err = loadBuildContext(filepath.Join(path, entry.Name()), filepath.Join(relativePath, entry.Name()), tw)
			if err != nil {
				return err
			}
			continue
		}
		filename := filepath.Join(path, entry.Name())
		rdr, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open %s: %v", filename, err)
		}
		defer rdr.Close()
		data, err := io.ReadAll(rdr)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", filename, err)
		}
		// Use ToSlash to make the filepath generic. This solves the issue where
		// Windows uses backslashes and Docker uses forward slashs.
		name := filepath.ToSlash(filepath.Join(relativePath, entry.Name()))
		header := &tar.Header{
			Name: name,
			Size: int64(len(data)),
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return fmt.Errorf("failed to write tar header for %s: %v", name, err)
		}
		_, err = tw.Write(data)
		if err != nil {
			return fmt.Errorf("failed to write tar body for %s: %v", name, err)
		}
	}
	return nil
}

func logBuildOutput(r io.Reader, name string) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read image build output: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fmt.Printf("[%s]: %s\n", name, line)
		var errorData map[string]interface{}
		err = json.Unmarshal([]byte(line), &errorData)
		if err == nil {
			message, ok := errorData["error"]
			if ok {
				fmt.Printf("[%s]: failed to build image: %s\n", name, message)
				return fmt.Errorf("failed to build image: %s", message)
			}
		}
	}
	return nil
}
