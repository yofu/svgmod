package svgmod

import (
	"fmt"
	"strings"
)

const (
	NORMAL = iota
	SUB
	SUP
)

var (
	greek = map[string]string{
		"\\alpha":   "α",
		"\\beta":    "β",
		"\\gamma":   "γ",
		"\\delta":   "δ",
		"\\epsilon": "ε",
		"\\zeta":    "ζ",
		"\\eta":     "η",
		"\\theta":   "θ",
		"\\iota":    "ι",
		"\\kappa":   "κ",
		"\\lambda":  "λ",
		"\\mu":      "μ",
		"\\nu":      "ν",
		"\\xi":      "ξ",
		"\\pi":      "π",
		"\\rho":     "ρ",
		"\\sigma":   "σ",
		"\\tau":     "τ",
		"\\upsilon": "υ",
		"\\phi":     "φ",
		"\\chi":     "χ",
		"\\psi":     "ψ",
		"\\omega":   "ω",
		"\\Alpha":   "Α",
		"\\Beta":    "Β",
		"\\Gamma":   "Γ",
		"\\Delta":   "Δ",
		"\\Epsilon": "Ε",
		"\\Zeta":    "Ζ",
		"\\Eta":     "Η",
		"\\Theta":   "Θ",
		"\\Iota":    "Ι",
		"\\Kappa":   "Κ",
		"\\Lambda":  "Λ",
		"\\Mu":      "Μ",
		"\\Nu":      "Ν",
		"\\Xi":      "Ξ",
		"\\Pi":      "Π",
		"\\Rho":     "Ρ",
		"\\Sigma":   "Σ",
		"\\Tau":     "Τ",
		"\\Upsilon": "Υ",
		"\\Phi":     "Φ",
		"\\Chi":     "Χ",
		"\\Psi":     "Ψ",
		"\\Omega":   "Ω",
	}
)

func tex2svg(txt, font string) string {
	rtn := []string{}
	tmp := []rune{}
	mode := NORMAL
	inbrace := false
	for _, s := range txt + " " {
		switch s {
		case ' ':
			if inbrace {
				tmp = append(tmp, s)
			} else {
				if len(tmp) != 0 {
					rtn = append(rtn, interpret(string(tmp), mode))
					tmp = []rune{' '}
					mode = NORMAL
				} else {
					tmp = []rune{' '}
				}
			}
		case '{':
			rtn = append(rtn, interpret(string(tmp), mode))
			tmp = []rune{}
			inbrace = true
		case '}':
			inbrace = false
			if len(tmp) != 0 {
				rtn = append(rtn, interpret(string(tmp), mode))
				tmp = []rune{}
			}
			mode = NORMAL
		case '_':
			rtn = append(rtn, interpret(string(tmp), mode))
			tmp = []rune{}
			mode = SUB
		case '^':
			rtn = append(rtn, interpret(string(tmp), mode))
			tmp = []rune{}
			mode = SUP
		default:
			tmp = append(tmp, s)
			if !inbrace && mode != NORMAL {
				rtn = append(rtn, interpret(string(tmp), mode))
				tmp = []rune{}
				mode = NORMAL
			}
		}
	}
	return fmt.Sprintf(`<tspan style="font-style:normal;font-variant:normal;font-weight:normal;font-stretch:normal;font-family:%s;-inkscape-font-specification:%s">%s</tspan>`, font, font, strings.Join(rtn, ""))
}


func interpret(txt string, mode int) string {
	if txt == " " {
		return " "
	}
	var style string
	if strings.HasPrefix(txt, "\\rm") {
		txt = strings.TrimPrefix(txt, "\\rm")
		txt = strings.TrimLeft(txt, " ")
	} else {
		style = `<tspan font-style="italic">`
	}
	rtn := txt
	if g, ok := greek[txt]; ok {
		rtn = g
	}
	switch mode {
	case SUB:
		rtn = fmt.Sprintf(`<tspan style="font-size:65%%;baseline-shift:sub">%s</tspan>`, rtn)
	case SUP:
		rtn = fmt.Sprintf(`<tspan style="font-size:65%%;baseline-shift:super">%s</tspan>`, rtn)
	}
	if style == "" {
		return rtn
	} else {
		return fmt.Sprintf("%s%s</tspan>", style, rtn)
	}
}
