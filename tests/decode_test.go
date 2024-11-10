package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/codecrafters-io/bittorrent-starter-go/decoder"
	"github.com/codecrafters-io/bittorrent-starter-go/utils"
)

var DecodeBencodeTests = []struct {
	description    string
	input          string
	expectedValue  interface{}
	expectedLength int
	expectedError  error
}{
	{
		description:    "Decode string",
		input:          "8:Valentin",
		expectedValue:  "Valentin",
		expectedLength: 10,
		expectedError:  nil,
	},
	{
		description:    "Decode integer",
		input:          "i123e",
		expectedValue:  123,
		expectedLength: 5,
		expectedError:  nil,
	},
	{
		description:    "Decode list",
		input:          "l3:one3:twoi3ee",
		expectedValue:  []interface{}{"one", "two", 3},
		expectedLength: 15,
		expectedError:  nil,
	},
	{
		description:    "Decode dictionary",
		input:          "d3:onei1e3:twoi2ee",
		expectedValue:  map[string]interface{}{"one": 1, "two": 2},
		expectedLength: 18,
		expectedError:  nil,
	},
	{
		description: "Decode nested list",
		input:       "l3:one3:twoi3el3:one3:twoi3eee",
		expectedValue: []interface{}{
			"one",
			"two",
			3,
			[]interface{}{
				"one",
				"two",
				3,
			},
		},
		expectedLength: 30,
	},
	{
		description:    "Decode unsupported type",
		input:          "k3:one3:twoi3ez",
		expectedValue:  "",
		expectedLength: 0,
		expectedError:  fmt.Errorf(decoder.ErrUnsupportedBencodeType),
	},
}

func TestDecodeBencode(t *testing.T) {
	for _, test := range DecodeBencodeTests {
		t.Run(test.description, func(t *testing.T) {
			value, length, err := decoder.DecodeBencode(test.input)
			if !reflect.DeepEqual(value, test.expectedValue) {
				t.Errorf("Expected value to be %v, got %v", test.expectedValue, value)
			}
			if length != test.expectedLength {
				t.Errorf("Expected length to be %v, got %v", test.expectedLength, length)
			}
			if err != nil {
				if err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error to be %v, got %v", test.expectedError, err)
				}
			}
		})
	}
}

func TestDecodeTorrentFile(t *testing.T) {
	t.Run("Decode torrent file", func(t *testing.T) {
		// The file is at the root of the project
		fileContent, err := utils.ReadFile("../sample.torrent")
		if err != nil {
			t.Errorf("unable to open the sample torrent file: %v", err)
		}
		torrent, _, err := decoder.DecodeTorrentFile(fileContent.String())
		if err != nil {
			t.Errorf("Expected error to be nil, got %v", err)
		}
		expected := map[string]interface{}{
			"announce": "http://bittorrent-test-tracker.codecrafters.io/announce",
			"info": map[string]interface{}{
				"length": 92063,
			},
		}
		// Check that decoded contains the expected values
		if expected["announce"] != torrent.Announce {
			t.Errorf("Expected announce to be %v, got %v", expected["announce"], torrent.Announce)
		}
		if expected["info"].(map[string]interface{})["length"] != torrent.Length {
			t.Errorf("Expected length to be %v, got %v", expected["info"].(map[string]interface{})["length"], torrent.Length)
		}
	})
}
