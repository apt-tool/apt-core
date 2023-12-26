package ai

import "math/rand"

type AI struct {
	Cfg Config
}

func (a AI) GetAttacks(list, vulnerabilities []string) []string {
	records := make([]string, 0)

	if !a.Cfg.Enable {
		return list
	}

	switch a.Cfg.Method {
	case "svm":
		break
	case "nbias":
		break
	default:
		for _, item := range list {
			if rand.Intn(a.Cfg.Limit) > a.Cfg.Factor {
				records = append(records, item)
			}
		}
	}

	return records
}
