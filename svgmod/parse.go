package svgmod

import (
	"fmt"
	"regexp"
)

var (
	subpat = regexp.MustCompile("^s/([^/]+)/([^/]*)/?$")
	label  = regexp.MustCompile("^([xy]label) (.+)$")
)

func Parse(txt string) (*Command, error) {
	switch {
	case subpat.MatchString(txt):
		fs := subpat.FindStringSubmatch(txt)
		return CommandSubstitute(fs[1], tex2svg(fs[2]))
	case label.MatchString(txt):
		fs := label.FindStringSubmatch(txt)
		return CommandSubstitute(fmt.Sprintf("$%s$", fs[1]), tex2svg(fs[2]))
	default:
		return nil, nil
	}
}
