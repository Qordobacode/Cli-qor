package file

import (
	"fmt"
	"testing"
)

func Test_Pattern(t *testing.T) {
	regexp := buildFilepathReplaceRegexp("en")
	source := "./v510-en/topics-other/ncc/file.​json"
	//result := regexp.FindStringSubmatch(source)
	rs := regexp.ReplaceAllStringFunc(source, func(m string) string {
		return regexp.ReplaceAllString(m, `"${2}whatever"`)
	})
	fmt.Printf("result = %v", rs)
}
