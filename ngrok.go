package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"

	. "ytd/constants"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2"
)

type NgrokService struct {
	runtime   *wails.Runtime
	pid       int
	publicUrl string
}

type NgrokProcessResult struct {
	err       error
	output    string
	errCode   string
	status    string
	publicUrl string
}

type NgrokTunnelInfo struct {
	PublicUrl string `json:"url"`
	Err       error
}

type ngrokTunnels struct {
	Tunnels []ngrokTunnel `json:"tunnels"`
}

type ngrokTunnel struct {
	Name      string            `json:"name"`
	URI       string            `json:"uri"`
	PublicURL string            `json:"public_url"`
	Proto     string            `json:"proto"`
	Config    ngrokTunnelConfig `json:"config"`
}

type ngrokTunnelConfig struct {
	Addr    string `json:"addr"`
	Inspect bool   `json:"inspect"`
}

type ngrokCmdError struct {
	err     error
	errCode string
}

func newProcessResultWithError(err error, output string, errCode string) NgrokProcessResult {
	return NgrokProcessResult{
		status:  NgrokStatusError,
		err:     err,
		output:  output,
		errCode: errCode,
	}
}

func newProcessResultWithUrl(url string) NgrokProcessResult {
	return NgrokProcessResult{
		status:    NgrokStatusRunning,
		publicUrl: url,
	}
}

func (n *NgrokService) isRunning() bool {
	return n.pid != 0
}

func (n *NgrokService) StartProcess(restart bool) NgrokProcessResult {
	// if is running kill & restart to apply new config
	if restart && n.isRunning() {
		n.KillProcess()
	}

	// set ngrok authtoken first
	err := n.SetAuthToken()
	if err != nil {
		return newProcessResultWithError(errors.Wrap(err, "ngrok StartProcess()"), "", "")
	}

	// ngrok --http http://localhost:8080 --region eu --bind-tls true --auth "ytd:ytd"
	ngrokPath, _ := n.IsNgrokInstalled()
	args := []string{
		"http", fmt.Sprintf("http://localhost:8080"),
		"--region", "eu",
		"--bind-tls", "true",
		// "--log=stdout",
	}
	if appState.Config.PublicServer.Ngrok.Auth.Enabled {
		args = append(args, "--auth")
		args = append(args, fmt.Sprintf("%s:%s", appState.Config.PublicServer.Ngrok.Auth.Username, appState.Config.PublicServer.Ngrok.Auth.Password))
	}
	cmd := exec.Command(
		ngrokPath,
		args...,
	)
	// Use the same pipe for standard error
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return newProcessResultWithError(errors.Wrap(err, "ngrok StartProcess() cmd.StderrPipe()"), "", "")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return newProcessResultWithError(errors.Wrap(err, "ngrok StartProcess() cmd.StdoutPipe()"), "", "")
	}

	if err := cmd.Start(); err != nil {
		return newProcessResultWithError(errors.Wrap(err, "ngrok StartProcess() cmd.Start()"), "", "")
	}

	// channel will send an error if ngrok quits unexpectedly.
	errorChan := make(chan ngrokCmdError)
	go errorReciever(cmd, stdout, stderr, errorChan)

	// channel will recieve the string of the connection URL.
	waitForConnectionChan := make(chan NgrokTunnelInfo, 1)
	go n.GetPublicUrl(waitForConnectionChan)

	// and finally, make a channel that will time out if all else fails.
	timeoutChan := time.After(20 * time.Second)

	// wait for something to happen...
	for {
		select {
		case info := <-waitForConnectionChan:
			fmt.Println("CONN INFO READ CHANNLEEEEE", info)
			if info.Err != nil {
				return newProcessResultWithError(errors.Wrap(err, "ngrok StartProcess() cmd.StderrPipe()"), "", "")
			}
			n.pid = cmd.Process.Pid
			return newProcessResultWithUrl(info.PublicUrl)
		case errResult := <-errorChan:
			fmt.Println("ERRROR READ CHANNLEEEEE", err)
			return newProcessResultWithError(errResult.err, "", errResult.errCode)
		case <-timeoutChan:
			fmt.Println("TIMEOUT READ CHANNLEEEEE")
			return newProcessResultWithError(errors.New("Ngrok start was timeouted"), "", "ERR_NGROK_START_TIMEOUT")
		default:
			time.Sleep(500 * time.Millisecond)
			fmt.Println("DEFAULT CASEEEE")
		}
	}
}

func (n *NgrokService) KillProcess() error {
	if n.pid == 0 {
		return nil
	}
	pgid, err := syscall.Getpgid(n.pid)
	if err != nil {
		return err
	}
	err = syscall.Kill(-pgid, 15) // note the minus sign
	if err != nil {
		return err
	}
	return nil
}

func (n *NgrokService) SetAuthToken() error {
	ngrokPath, _ := n.IsNgrokInstalled()
	args := []string{
		"authtoken", appState.Config.PublicServer.Ngrok.Authtoken,
	}
	cmd := exec.Command(
		ngrokPath,
		args...,
	)
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "ngrok SetAuthToken()")
	}
	return nil
}

func (n *NgrokService) GetPublicUrl(waitForConnectionChan chan NgrokTunnelInfo) {
	time.Sleep(3 * time.Second)
	tunnels := &ngrokTunnels{}
	res, err := http.Get("http://localhost:4040/api/tunnels")
	if err != nil {
		waitForConnectionChan <- NgrokTunnelInfo{Err: err}
		return
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		waitForConnectionChan <- NgrokTunnelInfo{Err: errors.New("Ngrok api/tunnels bad status code")}
		return
	}

	if err := json.NewDecoder(res.Body).Decode(&tunnels); err != nil {
		waitForConnectionChan <- NgrokTunnelInfo{Err: err}
		return
	}
	n.publicUrl = tunnels.Tunnels[0].PublicURL
	waitForConnectionChan <- NgrokTunnelInfo{PublicUrl: tunnels.Tunnels[0].PublicURL}
	return
}

func (n *NgrokService) IsNgrokInstalled() (string, error) {
	ngrok, err := exec.LookPath("ngrok")
	if n.runtime.System.Platform() == "darwin" && err != nil {
		// on darwin check if ngrok is maybe installed by homebrew
		// (searching for ngrok only give wrong results if installed with homebrew)
		ngrok, err := exec.LookPath("/opt/homebrew/bin/ngrok")
		return ngrok, err
	}
	return ngrok, err
}

func errorReciever(cmd *exec.Cmd, stdoutPipe io.ReadCloser, stderrPipe io.ReadCloser, errorChan chan ngrokCmdError) {
	fmt.Println("STARTED ERROR RECEIVER")
	stdout, err := ioutil.ReadAll(stdoutPipe)
	if err != nil {
		errorChan <- ngrokCmdError{err: errors.Wrap(err, "errorReciever ioutil.ReadAll(stdoutPipe)")}
		return
	}

	stderr, err := ioutil.ReadAll(stderrPipe)
	if err != nil {
		errorChan <- ngrokCmdError{err: errors.Wrap(err, "errorReciever ioutil.ReadAll(stderrPipe)")}
		return
	}
	// in the "happy case", there is no output from ngrok. So if there is ANY
	//  output, we treat it as an error.
	/* 	if len(output) > 0 {
		errorChan <- errors.Wrap(err, "errorReciever len(output) > 0")
		return
	} */
	// otherwise, we wait on the process to retrieve it's potentially non-
	//  zero exit code.
	err = cmd.Wait()
	if err != nil {
		fmt.Println("----------------------------")
		fmt.Println("errorReciever ERROR NGROK", err)
		fmt.Println("STDOUT", extractErrorCode(string(stdout)))
		fmt.Println("STDERR", extractErrorCode(string(stderr)))
		fmt.Println("----------------------------")
		errorChan <- ngrokCmdError{err: errors.Wrap(err, "errorReceiver() cmd.Wait()"), errCode: extractErrorCode(string(stdout))}
	}
	// close the channel before ending the goroutine.
	close(errorChan)
	return
}

func extractErrorCode(output string) string {
	if strings.Contains(output, "ERR_NGROK") {
		re := regexp.MustCompile("ERR_NGROK_([0-9].*)")
		errCode := re.FindString(output)
		return strings.TrimSpace(errCode)
	}
	return ""
}
