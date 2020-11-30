package beauter

import "github.com/hashicorp/hcl/v2/hclwrite"

func formatHCL(content string) (formattedContent string, err error) {
	formattedContent = string(hclwrite.Format([]byte(content)))
	return
}
