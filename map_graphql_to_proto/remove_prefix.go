package map_graphql_to_proto

import "strings"

func removePrefix(input string) string {
	config := GetConfig()

	for _, prefix := range config.RemovePrefixes {
		if strings.HasPrefix(input, prefix) {
			return input[len(prefix):]
		}
	}

	return input
}
