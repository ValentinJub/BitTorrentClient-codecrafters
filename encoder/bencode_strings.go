package encoder

import (
	"fmt"
	"slices"
)

const ErrUnsupportedBencodeType = "unsupported bencode type"

func EncodeBencode(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return encodeString(v), nil
	case int:
		return encodeInteger(v), nil
	case []interface{}:
		return encodeList(v)
	case map[string]interface{}:
		return encodeDictionary(v)
	default:
		return "", fmt.Errorf(ErrUnsupportedBencodeType)
	}
}

func encodeString(value string) string {
	return fmt.Sprintf("%d:%s", len(value), value)
}

func encodeInteger(value int) string {
	return fmt.Sprintf("i%de", value)
}

func encodeList(value []interface{}) (string, error) {
	var encodedList string
	for _, v := range value {
		encoded, err := EncodeBencode(v)
		if err != nil {
			return "", err
		}
		encodedList += encoded
	}
	return "l" + encodedList + "e", nil
}

func encodeDictionary(value map[string]interface{}) (string, error) {
	var encodedDictionary string
	keysOrder := make([]string, 0, len(value))
	for k := range value {
		keysOrder = append(keysOrder, k)
	}
	slices.Sort(keysOrder)
	for _, k := range keysOrder {
		encodedKey := encodeString(k)
		encodedValue, err := EncodeBencode(value[k])
		if err != nil {
			return "", err
		}
		encodedDictionary += encodedKey + encodedValue
	}
	return "d" + encodedDictionary + "e", nil
}
