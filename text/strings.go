package text

import (
	"fmt"
	"github.com/mvdan/xurls"
	"regexp"
	"strings"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}

func Quantify(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}

	return fmt.Sprintf("%d %s", n, plural)
}

func IsURLStrict(s string) bool {
	urls := xurls.Strict.FindAllString(s, -1)
	if len(urls) != 1 {
		return false
	}

	return urls[0] == s
}

func RegexpGroups(pattern *regexp.Regexp, src string) map[string]string {
	match := pattern.FindStringSubmatch(src)
	if match == nil {
		return nil
	}

	groups := map[string]string{}
	for i, name := range pattern.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}

		groups[name] = match[i]
	}

	return groups
}
