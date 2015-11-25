package svgmod

import (
	"fmt"
	"regexp"
)

var (
	subpat = regexp.MustCompile("^s/([^/]+)/([^/]*)/?$")
	subtexpat = regexp.MustCompile("^t/([^/]+)/([^/]*)/?$")
	label  = regexp.MustCompile("^([xy]label) (.+)$")
	title  = regexp.MustCompile("^(title[0-9]*) (.+)$")
)

func Parse(txt, font string) (*Command, error) {
	switch {
	case subpat.MatchString(txt):
		fs := subpat.FindStringSubmatch(txt)
		return CommandSubstitute(fs[1], fs[2])
	case subtexpat.MatchString(txt):
		fs := subtexpat.FindStringSubmatch(txt)
		return CommandSubstitute(fs[1], tex2svg(fs[2], font))
	case label.MatchString(txt):
		fs := label.FindStringSubmatch(txt)
		return CommandSubstitute(fmt.Sprintf("$%s$", fs[1]), tex2svg(fs[2], font))
	case title.MatchString(txt):
		fs := title.FindStringSubmatch(txt)
		return CommandSubstitute(fmt.Sprintf("$%s$", fs[1]), fs[2])
	default:
		return nil, nil
	}
}
