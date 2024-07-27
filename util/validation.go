package util

import "regexp"

var (
	ReWhiteSpace     *regexp.Regexp
	ReFieldWhiteList *regexp.Regexp
)

func init() {
	ReWhiteSpace, _ = regexp.Compile(`\s`)
	ReFieldWhiteList, _ = regexp.Compile(`^[A-Za-z0-9]+$`)
}

func FilterField(field string) bool {
	return ReFieldWhiteList.MatchString(field)
}
