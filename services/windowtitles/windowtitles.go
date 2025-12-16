package windowtitles

import (
	"strings"
)

const (
	Delimiter = "_"
	Separator = ":"

	ComponentOrientation = "O"
	ComponentName        = "N"
	ComponentId          = "ID"

	ApplicationReader = "com.lab126.booklet.reader"
)

type KindleTitle string

func (title KindleTitle) Get(key string) (string, bool) {
	for _, component := range strings.Split(string(title), Delimiter) {
		start := key + Separator
		if strings.HasPrefix(component, start) {
			return strings.Replace(component, start, "", -1), true
		}
	}
	return "", false
}

func (title KindleTitle) Set(key, value string) KindleTitle {
	components := strings.Split(string(title), Delimiter)
	prefix := key + Separator
	replaced := false

	for i, component := range components {
		if strings.HasPrefix(component, prefix) {
			components[i] = prefix + value
			replaced = true
			break
		}
	}

	if replaced {
		return KindleTitle(strings.Join(components, Delimiter))
	}

	return title
}

func (title KindleTitle) IsApplication(name string) bool {
	if strings.Contains(string(title), ComponentName+Separator+"application") {
		return strings.Contains(string(title), ComponentId+Separator+name)
	}
	return false
}
