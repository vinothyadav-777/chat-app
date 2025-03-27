package utils

import (
	"bytes"
	"github.com/vinothyadav-777/chat-app/constants"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

func GetNuid() string {
	return strings.Replace(uuid.New().String(), constants.Hyphen, "", -1)
}

func JsonUnmarshal(jsonString string, output interface{}) error {
	d := json.NewDecoder(bytes.NewReader([]byte(jsonString)))
	// for this issue https://github.com/golang/go/issues/5562#issuecomment-66080065
	d.UseNumber()

	err := d.Decode(output)
	if err != nil {
		return err
	}
	return nil
}

func SplitIntoSizedChunks(chunkSize int, slice []string) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
