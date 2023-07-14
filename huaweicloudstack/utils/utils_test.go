package utils

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	greenCode  = "\033[1;32m"
	yellowCode = "\033[1;33m"
	resetCode  = "\033[0m"
)

func green(str interface{}) string {
	return fmt.Sprintf("%s%#v%s", greenCode, str, resetCode)
}

func yellow(str interface{}) string {
	return fmt.Sprintf("%s%#v%s", yellowCode, str, resetCode)
}

func TestAccFunction_RemoveNil(t *testing.T) {
	var (
		testInput = map[string]interface{}{
			"level_one_index_zero": nil,
			"level_one_index_one": []map[string]interface{}{
				{
					"level_two_index_zero": nil,
				},
				{
					"level_two_index_one": "192.168.0.1",
				},
			},
			"level_one_index_two": []map[string]interface{}{
				{
					"level_two_index_zero": nil,
				},
			},
			"level_one_index_three": "172.16.0.237",
		}

		expected = map[string]interface{}{
			"level_one_index_one": []map[string]interface{}{
				{
					"level_two_index_one": "192.168.0.1",
				},
			},
			"level_one_index_three": "172.16.0.237",
		}
	)

	if !reflect.DeepEqual(RemoveNil(testInput), expected) {
		t.Fatalf("The processing result of RemoveNil method is not as expected, want %s, but %s", green(expected), yellow(testInput))
	}
	t.Logf("The processing result of RemoveNil method meets expectation: %s", green(expected))
}

func TestAccFunction_reverse(t *testing.T) {
	var (
		testInput = "abcdefg"
		expected  = "gfedcba"
	)

	if !reflect.DeepEqual(Reverse(testInput), expected) {
		t.Fatalf("The processing result of the function 'Reverse' is not as expected, want '%s', but got '%s'",
			green(expected), yellow(testInput))
	}
	t.Logf("The processing result of function 'Reverse' meets expectation: %s", green(expected))
}

func TestAccFunction_jsonStringsEqual(t *testing.T) {
	var (
		jsonStr1 = "{\n\"key1\":\"value1\",\n\"key2\":\"value2\"\n}"
		jsonStr2 = "{\"key2\":\"value2\",\"key1\":\"value1\"}"
	)

	if !JSONStringsEqual(jsonStr1, jsonStr2) {
		t.Fatalf("The processing result of the function 'JSONStringsEqual' is not as expected, want '%v', "+
			"but got '%v'", green(true), yellow(false))
	}
	t.Logf("The processing result of function 'JSONStringsEqual' meets expectation: %s", green(true))
}
