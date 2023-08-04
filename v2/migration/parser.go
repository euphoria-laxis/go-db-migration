package migration

import "strings"

func parseTag(tag string) map[string]string {
	parsed := make(map[string]string)
	items := strings.Split(tag, ";")
	for _, item := range items {
		s := strings.Split(item, ":")
		if len(s) != 2 {
			continue
		}
		key := s[0]
		value := s[1]
		parsed[key] = value
	}

	return parsed
}
