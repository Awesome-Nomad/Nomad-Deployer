package beauter

import (
	"bytes"
	"encoding/json"
)

func formatJSON(jsonStr string) (formattedContent string, err error) {
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(jsonStr), "", "\t")
	if err != nil {
		return
	}
	formattedContent = prettyJSON.String()
	return
}
