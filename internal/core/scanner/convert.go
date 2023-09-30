package scanner

import (
	"encoding/json"
)

// convert scanner report to our system report
func convert(context []byte, r *report) error {
	p := make([]string, 0)

	// convert our json object into our report
	if er := json.Unmarshal(context, &p); er != nil {
		return er
	}

	// adding vulnerabilities into report list
	r.vulnerabilities = p

	return nil
}
