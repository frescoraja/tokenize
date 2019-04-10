package tokenize

import (
	"errors"
	"strings"
)

// ValueBoundaryCount is number of times given boundary rune must repeat in order to signify value start/end
const ValueBoundaryCount = 3

// Param struct to hold parameter_name and parameter_value
type Param struct {
	Name  string
	Value string
}

// func main() {
// filename := flag.String("filename", "params.txt", "file containing params input string.")
// flag.Parse()

// f, err := ioutil.ReadFile(*filename)
// s := string(f)

// ps, err := GetParams(s, []rune{'|', '\n'}, '`', '\\')
// if err != nil {
// panic(err)
// }

// for _, p := range ps {
// fmt.Printf("thing: %+v\n", p)
// }
// }

// GetParams converts slice of strings to slice of Params
func GetParams(s string, paramSeps []rune, valDel rune, escape rune) (params []Param, err error) {
	var trimmedTokens []string

	tokens, err := tokenizeParamString(s, paramSeps, valDel, escape)
	if err != nil {
		return nil, err
	}

	for _, t := range tokens {
		if strings.Contains(t, "\n") {
			trimmedTokens = append(trimmedTokens, t)
		} else {
			trimmedToken := strings.TrimSpace(t)
			if trimmedToken != "" {
				trimmedTokens = append(trimmedTokens, trimmedToken)
			}
		}
	}
	if len(trimmedTokens)%2 != 0 {
		return nil, errors.New("invalid parameters: error parsing input into parameter name/value pairs")
	}

	for x := 0; x < len(trimmedTokens); x += 2 {
		pname := trimmedTokens[x]
		pvalue := trimmedTokens[x+1]
		p := Param{Name: pname, Value: pvalue}
		params = append(params, p)
	}

	return params, nil
}

func tokenizeParamString(s string, paramSeps []rune, valueBoundary, escape rune) (tokens []string, err error) {
	var (
		runes             []rune
		boundaryCount     int
		inEscape, inValue bool
	)
	for _, r := range s {
		switch {
		case r == escape:
			if inValue {
				runes = append(runes, r)
			} else {
				inEscape = true
			}
		case r == valueBoundary:
			if inValue {
				if boundaryCount < (ValueBoundaryCount - 1) {
					boundaryCount++
				} else {
					tokens = append(tokens, string(runes))
					runes = runes[:0]
					inValue = false
					boundaryCount = 0
				}
			} else {
				if boundaryCount < (ValueBoundaryCount - 1) {
					boundaryCount++
				} else {
					inValue = true
					boundaryCount = 0
				}
			}
		case strings.ContainsRune(string(paramSeps), r):
			if inValue || inEscape {
				runes = append(runes, r)
			} else if !inEscape {
				if len(runes) > 0 {
					tokens = append(tokens, string(runes))
					runes = runes[:0]
				}
			}
			inEscape = false
		default:
			for x := 0; x < boundaryCount; x++ {
				runes = append(runes, valueBoundary)
			}
			boundaryCount = 0
			inEscape = false
			runes = append(runes, r)
		}
	}

	if inEscape {
		err = errors.New("invalid terminal escape")
	}
	if inValue {
		err = errors.New("invalid value termination")
	}

	tokens = append(tokens, string(runes))

	return tokens, err
}
