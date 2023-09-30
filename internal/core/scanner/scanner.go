package scanner

import (
	"fmt"
	"os/exec"
)

var (
	template = "python scanner.py --host %s"
)

// Scan a host by using apt-scanner
func Scan(host string) ([]string, error) {
	r := new(report)

	// create command
	command := fmt.Sprintf(template, host)

	// execute command
	cmd := exec.Command(command)
	if err := cmd.Start(); err != nil {
		return r.vulnerabilities, err
	}

	// read output
	context, err := cmd.Output()
	if err != nil {
		return r.vulnerabilities, err
	}

	// convert type to our report
	if er := convert(context, r); er != nil {
		return r.vulnerabilities, er
	}

	return r.vulnerabilities, nil
}
