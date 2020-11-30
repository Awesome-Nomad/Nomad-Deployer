package beauter

type ContentType string

const (
	HCL  ContentType = "application/vnd.hcl"
	JSON ContentType = "application/json"
	YAML ContentType = "text/yaml"
)

func (c ContentType) Format(content string) (formattedContent string, err error) {
	switch string(c) {
	case string(JSON):
		return formatJSON(content)
	case string(HCL):
		return formatHCL(content)
	default:
		return content, err
	}
}

func FormatJSON(content string) (formattedContent string, err error) {
	return formatJSON(content)
}

func FormatHCL(content string) (formattedContent string, err error) {
	return formatHCL(content)
}
