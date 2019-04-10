package tokenize

import (
	"errors"
	"strings"
)

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
		return nil, errors.New("invalid parameters: each name must have value")
	}

	for x := 0; x < len(trimmedTokens); x += 2 {
		pname := trimmedTokens[x]
		pvalue := trimmedTokens[x+1]
		p := Param{Name: pname, Value: pvalue}
		params = append(params, p)
	}

	return params, nil
}

func tokenizeParamString(s string, paramSeps []rune, longValueSep, escape rune) (tokens []string, err error) {
	var (
		runes                 []rune
		longSepCount          int
		inEscape, inLongValue bool
	)
	for _, r := range s {
		switch {
		default:
			for x := 0; x < longSepCount; x++ {
				runes = append(runes, longValueSep)
			}
			longSepCount = 0
			runes = append(runes, r)
			if inEscape {
				inEscape = false
			}
		case r == longValueSep:
			if inLongValue {
				if longSepCount < 2 {
					longSepCount++
				} else {
					tokens = append(tokens, string(runes))
					runes = runes[:0]
					inLongValue = false
					longSepCount = 0
				}
			} else {
				if longSepCount < 2 {
					longSepCount++
				} else {
					inLongValue = true
					longSepCount = 0
				}
			}
		case r == escape:
			if inLongValue {
				runes = append(runes, r)
			} else {
				inEscape = true
			}
		case strings.ContainsRune(string(paramSeps), r):
			if inLongValue || inEscape {
				runes = append(runes, r)
			} else if !inEscape {
				if len(runes) > 0 {
					tokens = append(tokens, string(runes))
					runes = runes[:0]
				}
			}
		}
	}
	tokens = append(tokens, string(runes))
	if inEscape {
		err = errors.New("invalid terminal escape")
	}
	return tokens, err
}
