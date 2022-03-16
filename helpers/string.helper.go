package helpers

import (
	"encoding/base64"
	"strings"
)

func ByteArraySeriesToBase64StringJoined(data [][]byte, joiner string) string {
	base64Strings := make([]string, 0)

	for _, value := range data {
		base64string := base64.StdEncoding.EncodeToString(value)
		base64Strings = append(base64Strings, base64string)
	}
	return strings.Join(base64Strings, joiner)
}

func Base64StringJoinedToByteArraySeries(base64Strings string, sep string) ([][]byte, error) {
	split := strings.Split(base64Strings, sep)
	byteArraySeries := make([][]byte, 0)

	for _, base64String := range split {
		byteArray, err := base64.StdEncoding.DecodeString(base64String)
		if err != nil {
			return nil, err
		}

		byteArraySeries = append(byteArraySeries, byteArray)
	}

	return byteArraySeries, nil
}

func ByteArraySeriesToString(byteArraySeries [][]byte) string {
	fulltext := ""

	for _, byteArray := range byteArraySeries {
		fulltext = fulltext + string(byteArray)
	}

	return fulltext
}
