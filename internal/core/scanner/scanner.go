package scanner

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	reportsDir = "./reports/"
)

var (
	template = "./scanner --url=%s --blockRegex true --reportFormat json --reportName %s --skipWAFIdentification true --skipWAFBlockCheck true --noEmailReport true"
)

// Scan a host by using apt-scanner
func Scan(host, name string) ([]string, error) {
	r := new(report)

	// create command
	command := fmt.Sprintf(template, host, name)

	// execute command
	cmd := exec.Command(command)
	if err := cmd.Start(); err != nil {
		return r.vulnerabilities, err
	}

	// read report file
	context, err := os.ReadFile(reportsDir + name + ".json")
	if err != nil {
		return r.vulnerabilities, err
	}

	// convert type to our report
	if er := convert(context, r); er != nil {
		return r.vulnerabilities, er
	}

	return r.vulnerabilities, nil
}
