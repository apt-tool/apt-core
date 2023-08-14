package scanner

import (
	"encoding/json"

	"github.com/amirhnajafiz/apt-scanner/pkg/proto"
)

// convert scanner report to our system report
func convert(context []byte, r *report) error {
	p := new(proto.Proto)

	// convert our json object into our report
	if er := json.Unmarshal(context, p); er != nil {
		return er
	}

	// adding vulnerabilities into report list
	for key := range p.Summary.Positive.Tests {
		r.vulnerabilities = append(r.vulnerabilities, key)
	}

	return nil
}
