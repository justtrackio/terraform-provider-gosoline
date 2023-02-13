package builder

import (
	"fmt"
	"strings"
)

func Augment(input string, appId AppId, additionalReplacements ...map[string]string) string {
	values := map[string]string{
		"project": appId.Project,
		"env":     appId.Environment,
		"family":  appId.Family,
		"group":   appId.Group,
		"app":     appId.Application,
	}

	for _, replacements := range additionalReplacements {
		for key, val := range replacements {
			values[key] = val
		}
	}

	for key, val := range values {
		templ := fmt.Sprintf("{%s}", key)
		input = strings.ReplaceAll(input, templ, val)
	}

	return input
}
