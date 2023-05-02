package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func prettyJSON(r io.Reader) {
	bytesFromReader, err := io.ReadAll(r)
	if err != nil {
		fmt.Printf("failed to read data: %v", err)
		return
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, bytesFromReader, "", "\t")
	if err != nil {
		fmt.Printf("unable to created pretty JSON: %v", err)
		return
	}

	fmt.Println(prettyJSON.String())
}
