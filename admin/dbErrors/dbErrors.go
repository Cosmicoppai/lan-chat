package dbErrors

import "github.com/lib/pq"

func InternalServerError(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		codeClass := err.Code.Class()
		if codeClass == "XX" || codeClass == "08" {
			return true
		}
		return false
	}
	return false
}

func IntegrityViolation(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		codeClass := err.Code.Class()
		if codeClass == "23" {
			return true
		}
		return false
	}
	return false
}
