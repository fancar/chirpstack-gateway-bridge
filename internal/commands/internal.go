package commands

import (
	"archive/zip"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

// internal commands to process some internal logic

func (cmd *command) ExecInternal(command string) (string, []byte, []byte, error) {
	log.WithFields(log.Fields{
		"command":                command,
		"max_execution_duration": cmd.MaxExecutionDuration,
	}).Info("execCmd:internal: executing internal command ...")

	switch command {
	case "get_logs":
		return cmd.getFiles("/var/log/")
	default:
		return "", nil, nil, fmt.Errorf("execCmd:internal: unknown internal command: %s", command)
	}
}

func (cmd *command) getFiles(fpath string) (string, []byte, []byte, error) {
	var buf bytes.Buffer
	pref := fmt.Sprintf("execCmd:Internal:getFiles:%s:", fpath)
	zw := zip.NewWriter(&buf)
	fnames, err := ioutil.ReadDir(fpath)
	if err != nil {
		return "", nil, nil, fmt.Errorf("%s unable to read log dir. Error: %v", pref, err)
	}

	c := 0
	for _, name := range fnames {
		if name.IsDir() {
			continue
		}
		// for _, file := range files {
		file, err := os.Open(fpath + name.Name())
		defer file.Close()
		if err != nil {
			log.Errorf("%s unable to open the file %s: %v\n", pref, name.Name(), err)
			continue
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Errorf("%s unable to read the file %s: %v\n", pref, file.Name(), err)
			continue
		}

		f, err := zw.Create(file.Name())
		if err != nil {
			return "", nil, nil, fmt.Errorf("%s unable to add file %s to zip. Error: %v", pref, file.Name(), err)
		}

		_, err = f.Write(data)
		if err != nil {
			return "", nil, nil, fmt.Errorf("%s unable to write file %s to zip. Error: %v", pref, file.Name(), err)
		}
		c++
	}

	if err := zw.Close(); err != nil {
		return "", nil, nil, fmt.Errorf("%s unable to close zip while compressing output of cmd: %s error: %v", pref, cmd, err)
	}

	if c == 0 {
		return "", nil, nil, fmt.Errorf("%s unable to zip any file in the directory", pref)
	}
	log.Debug("%s zipped %d files with size: %d\n", pref, c, len(buf.Bytes()))

	return cmd.Name, buf.Bytes(), nil, nil
}
