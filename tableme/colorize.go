package tableme

import (
	"bytes"
	"errors"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"regexp"
)

var (
	colorFuncsCache = make(map[string]func())
)

type ColorRule struct {
	Pattern         string
	CompiledPattern *regexp.Regexp
	Color           string
}

func compileRegexps(colorRules []*ColorRule) {
	for _, rule := range colorRules {
		re := rule.Pattern
		rule.CompiledPattern = regexp.MustCompile(re)
	}
}

func getColorFunc(colorName string) (func(...interface{}) string, error) {
	var fn func(...interface{}) string

	switch colorName {
	case "red":
		fn = color.New(color.FgRed).SprintFunc()
	case "red_bold":
		fn = color.New(color.FgRed, color.Bold).SprintFunc()
	case "blue":
		fn = color.New(color.FgBlue).SprintFunc()
	case "blue_bold":
		fn = color.New(color.FgBlue, color.Bold).SprintFunc()
	case "cyan":
		fn = color.New(color.FgCyan).SprintFunc()
	case "cyan_bold":
		fn = color.New(color.FgCyan, color.Bold).SprintFunc()
	case "yellow":
		fn = color.New(color.FgYellow).SprintFunc()
	case "yellow_bold":
		fn = color.New(color.FgYellow, color.Bold).SprintFunc()
	case "green":
		fn = color.New(color.FgGreen).SprintFunc()
	case "green_bold":
		fn = color.New(color.FgGreen, color.Bold).SprintFunc()
	case "magenta":
		fn = color.New(color.FgMagenta).SprintFunc()
	case "magenta_bold":
		fn = color.New(color.FgMagenta, color.Bold).SprintFunc()
	case "white":
		fn = color.New(color.FgWhite).SprintFunc()
	case "white_bold":
		fn = color.New(color.FgWhite, color.Bold).SprintFunc()
	default:
		return nil, errors.New(fmt.Sprintf("unknown color: ", colorName))
	}

	return fn, nil
}

func Colorize(raw []byte, colorRules []*ColorRule) *bytes.Buffer {
	compileRegexps(colorRules)

	var output []byte = raw

	for _, rule := range colorRules {
		//spew.Dump(rule)
		output = rule.CompiledPattern.ReplaceAllFunc(output, func(match []byte) []byte {
			//spew.Dump(match)
			fn, _ := getColorFunc(rule.Color)
			return []byte(fn(string(match)))
		})
	}

	return bytes.NewBuffer(output)
}
