package exec

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func Run(exe string, parms []string) ([]byte, []byte, error) {

	// log.Printf("#### Run: %v %v", exe, strings.Join(parms, " "))

	cmd := exec.Command(exe, parms...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("Run, StdoutPipe err: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("Run, StderrPipe err: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("Run, Start err: %v", err)
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, nil, fmt.Errorf("Run, ReadAll stderr err: %v", err)
	}

	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, nil, fmt.Errorf("Run, ReadAll output err: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return output, bytesErr, err
	}

	return output, bytesErr, nil
}
