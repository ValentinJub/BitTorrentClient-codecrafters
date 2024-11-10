package decoder

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"
)

const ErrUnsupportedBencodeType = "unsupported bencode type"

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func DecodeBencode(bencodedString string) (interface{}, int, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		return decodeString(bencodedString)
	} else if bencodedString[0] == 'i' {
		return decodeInteger(bencodedString)
	} else if bencodedString[0] == 'l' {
		return decodeList(bencodedString)
		// } else if bencodedString[0] == 'd' {
		// 	return decodeDictionary(bencodedString)
	} else {
		return "", 0, fmt.Errorf(ErrUnsupportedBencodeType)
	}
}

// Decode a bencoded string into a string value, returning the value, the number of bytes read, and an error if any
func decodeString(bencodedString string) (value string, bytesRead int, err error) {
	var firstColonIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", 0, err
	}

	value = bencodedString[firstColonIndex+1 : firstColonIndex+1+length]
	return value, len(lengthStr) + 1 + len(value), nil
}

// Decode a bencoded string into an integer value, returning the value, the number of bytes read, and an error if any
func decodeInteger(bencodedString string) (value, length int, err error) {
	integerRegex := regexp.MustCompile(`^i(-?\d+)e`)
	matches := integerRegex.FindStringSubmatch(bencodedString)
	if len(matches) == 0 {
		return 0, 0, fmt.Errorf("invalid integer format")
	}
	n, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, err
	}
	return n, len(matches[1]) + 2, nil
}

// Decode a bencoded string into a list value, returning the value, the number of bytes read, and an error if any
func decodeList(bencodedString string) (values []interface{}, length int, err error) {
	var list []interface{}
	if bencodedString[0] != 'l' {
		return nil, 0, fmt.Errorf("invalid list format")
	}
	bencodedString = bencodedString[1:]
	for len(bencodedString) > 0 && bencodedString[0] != 'e' {
		element, elementLength, err := DecodeBencode(bencodedString)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, element)
		bencodedString = bencodedString[elementLength:]
		length += elementLength
	}
	return list, length + 2, nil
}

// func decodeDictionary(bencodedString string) (values map[string]interface{}, length int, err error) {
// 	return map[string]interface{}{}, 0, nil
// }
