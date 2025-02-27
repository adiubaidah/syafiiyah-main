package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// function string to lower and snake case
func ToSnakeCase(s string) string {
	var result string
	for i, c := range s {
		if i > 0 && 'A' <= c && c <= 'Z' {
			result += "_"
		}
		result += string(c)
	}
	return result
}
func GetDeviceName(topic string) string {
	return strings.Split(topic, "/")[0]
}

func GetDeviceMode(topic string) string {
	return strings.Split(topic, "/")[2]
}

func Generate32ByteKey() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Printf("Failed to generate object key: %v\n", err)
		return ""
	}
	key := hex.EncodeToString(bytes)
	return key
}

func debugQuery(query string, args ...interface{}) string {
	var placeholderRegexp = regexp.MustCompile(`\$(\d+)`) // e.g. "$1"
	return placeholderRegexp.ReplaceAllStringFunc(query, func(ph string) string {
		// Get the index, e.g. for "$2" index becomes 2.
		idx, err := strconv.Atoi(ph[1:])
		if err != nil || idx < 1 || idx > len(args) {
			return ph
		}
		// Wrap value in quotes for readability.
		return fmt.Sprintf("'%v'", args[idx-1])
	})
}
