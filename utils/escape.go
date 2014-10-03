package utils

import "strings"

func Escape(input string) string {
	output := strings.Replace(input, "\n", "\\n", -1)
	output = strings.Replace(output, "\r", "\\r", -1)
	output = strings.Replace(output, "\t", "\\t", -1)
	//output = strings.Replace(output, "<", `\<`, -1)
	//output = strings.Replace(output, ">", `\>`, -1)
	//output = strings.Replace(output, "&", `\&`, -1)
	output = strings.Replace(output, `'`, `\'`, -1)
	output = strings.Replace(output, `"`, `\"`, -1)

	return output
}
