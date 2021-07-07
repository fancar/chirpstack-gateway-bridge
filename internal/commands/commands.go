package commands

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/brocaar/chirpstack-gateway-bridge/internal/config"
	"github.com/brocaar/chirpstack-gateway-bridge/internal/integration"
	"github.com/brocaar/lorawan"
)

type command struct {
	Name                 string // a template for name the answer
	Command              string
	CompressOutput       bool
	Internal             bool
	MaxExecutionDuration time.Duration
}

var (
	mux sync.RWMutex

	commands map[string]command
)

// Setup configures the gateway commands.
func Setup(conf config.Config) error {
	mux.Lock()
	defer mux.Unlock()

	commands = preconfiguredCommands() // make(map[string]command)

	for k, v := range conf.Commands.Commands {
		commands[k] = command{
			Command:              v.Command,
			MaxExecutionDuration: v.MaxExecutionDuration,
			CompressOutput:       v.CompressOutput,
		}

		log.WithFields(log.Fields{
			"command":                k,
			"command_exec":           v.Command,
			"max_execution_duration": v.MaxExecutionDuration,
		}).Info("commands: configuring command")
	}

	i := integration.GetIntegration()
	if i == nil {
		return errors.New("integration is not set")
	}

	i.SetGatewayCommandExecRequestFunc(gatewayCommandExecRequestFunc)

	return nil
}

// commands hardcoded by default. Mod
func preconfiguredCommands() map[string]command {
	commands := make(map[string]command)

	// connectivity check
	commands["ping"] = command{
		Command:              "echo pong",
		MaxExecutionDuration: 10 * time.Millisecond,
	}

	// reboot the device
	commands["reboot"] = command{
		Command:              "/sbin/reboot",
		MaxExecutionDuration: 1 * time.Second,
	}

	// stop packet forwarder
	commands["radio_stop"] = command{
		Command:              "/etc/init.d/erth_pf stop",
		MaxExecutionDuration: 10 * time.Second,
	}

	// restart packet forwarder
	commands["radio_restart"] = command{
		Command:              "/etc/init.d/erth_pf restart",
		MaxExecutionDuration: 10 * time.Second,
	}

	// start packet forwarder
	commands["radio_start"] = command{
		Command:              "/etc/init.d/erth_pf start",
		MaxExecutionDuration: 10 * time.Second,
	}

	// init ssh
	commands["reverse_ssh"] = command{
		Command:              "/home/erth/reverse_ssh.sh",
		MaxExecutionDuration: 5 * time.Second,
	}

	// spectrum analyser
	commands["spectral_scan"] = command{
		Name:                 time.Now().Format("spectrum_02012006_1504.txt"),
		Command:              "/home/erth/spectral_scan",
		MaxExecutionDuration: 10 * time.Minute,
		CompressOutput:       false,
	}

	// internal complex instructions. processing by methods in internal.go

	// zip and return all logs from /var/log
	commands["get_logs"] = command{
		Name:                 time.Now().Format("logs_02012006_1504.zip"),
		Internal:             true, // internal command with complex logic
		MaxExecutionDuration: 1 * time.Minute,
	}

	return commands
}

func gatewayCommandExecRequestFunc(pl gw.GatewayCommandExecRequest) {
	go executeCommand(pl)
}

func compress(cmd string, input []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	f, err := zw.Create(cmd)
	if err != nil {
		return nil, fmt.Errorf("unable to create zip for cmd: %s, error: %w", cmd, err)
	}

	_, err = f.Write(input)
	if err != nil {
		return nil, fmt.Errorf("unable to zip output of cmd: %s, error: %w", cmd, err)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("unable to close zip while compressing output of cmd: %s error: %w", cmd, err)
	}

	return buf.Bytes(), nil
}

func executeCommand(cmd gw.GatewayCommandExecRequest) {
	var gatewayID lorawan.EUI64
	copy(gatewayID[:], cmd.GatewayId)

	n, stdout, stderr, err := execute(cmd.Command, cmd.Stdin, cmd.Environment)
	resp := gw.GatewayCommandExecResponse{
		Name:      n,
		GatewayId: cmd.GatewayId,
		ExecId:    cmd.ExecId,
		Stdout:    stdout,
		Stderr:    stderr,
	}
	if err != nil {
		resp.Error = err.Error()
	}

	var id uuid.UUID
	copy(id[:], cmd.ExecId)

	if err := integration.GetIntegration().PublishEvent(gatewayID, "exec", id, &resp); err != nil {
		log.WithError(err).Error("commands: publish command execution event error")
	}
}

func execute(command string, stdin []byte, environment map[string]string) (string, []byte, []byte, error) {
	mux.RLock()
	defer mux.RUnlock()

	cmd, ok := commands[command]
	if !ok {
		return "", nil, nil, errors.New("command does not exist")
	}

	if cmd.Internal {
		return cmd.ExecInternal(command)
	}

	cmdArgs, err := ParseCommandLine(cmd.Command)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "parse command error")
	}
	if len(cmdArgs) == 0 {
		return "", nil, nil, errors.New("no command is given")
	}

	log.WithFields(log.Fields{
		"command":                command,
		"exec":                   cmdArgs[0],
		"args":                   cmdArgs[1:],
		"max_execution_duration": cmd.MaxExecutionDuration,
	}).Info("commands: executing command")

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(cmd.MaxExecutionDuration))
	defer cancel()

	cmdCtx := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)

	// The default is that when cmdCtx.Env is nil, os.Environ() are being used
	// automatically. As we want to add additional env. variables, we want to
	// extend this list, thus first need to set them to os.Environ()
	cmdCtx.Env = os.Environ()
	for k, v := range environment {
		cmdCtx.Env = append(cmdCtx.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdinPipe, err := cmdCtx.StdinPipe()
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "get stdin pipe error")
	}

	stdoutPipe, err := cmdCtx.StdoutPipe()
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "get stdout pipe error")
	}

	stderrPipe, err := cmdCtx.StderrPipe()
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "get stderr pipe error")
	}

	go func() {
		defer stdinPipe.Close()
		if _, err := stdinPipe.Write(stdin); err != nil {
			log.WithError(err).Error("commands: write to stdin error")
		}
	}()

	if err := cmdCtx.Start(); err != nil {
		return "", nil, nil, errors.Wrap(err, "starting command error")
	}

	stdoutB, _ := ioutil.ReadAll(stdoutPipe)
	stderrB, _ := ioutil.ReadAll(stderrPipe)

	if err := cmdCtx.Wait(); err != nil {
		return "", nil, nil, errors.Wrap(err, "waiting for command to finish error")
	}

	if len(stdoutB) > 0 && cmd.CompressOutput {
		stdoutB, err = compress(command, stdoutB)
		if err != nil {
			return "", nil, nil, err
		}
	}

	return cmd.Name, stdoutB, stderrB, nil
}

// ParseCommandLine parses the given command to commands and arguments.
// source: https://stackoverflow.com/questions/34118732/parse-a-command-line-string-into-flags-and-arguments-in-golang
func ParseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New(fmt.Sprintf("Unclosed quote in command line: %s", command))
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}
