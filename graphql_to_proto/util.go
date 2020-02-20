package graphql_to_proto

import (
	"strings"
)

func cleanNameFromPrefix(name string, config Config) string {

	if len(config.RemovePrefixes) < 1 {
		return name
	}

	for _, prefix := range config.RemovePrefixes {
		result := cleanNameFromSinglePrefix(name, prefix)
		if result != name {
			return result
		}
	}
	return name
}

func cleanNameFromSinglePrefix(name string, prefix string) string {
	if strings.HasPrefix(name, "repeated") {
		possiblePrefixedNames := name[9:]
		return "repeated " + cleanNameFromSinglePrefix(possiblePrefixedNames, prefix)
	}

	if !strings.HasPrefix(name, prefix) {
		return name
	}

	return name[len(prefix):]
}
