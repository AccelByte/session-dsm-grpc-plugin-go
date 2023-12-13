package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// MakeTraceID make traceID.
func MakeTraceID(identifiers ...string) string {
	strInt := strconv.Itoa(generateRandomInt())
	var tID string
	for _, i := range identifiers {
		tID = fmt.Sprintf(tID + i + "_")
	}

	return fmt.Sprintf(tID + strInt)
}

//nolint:gosec
func generateRandomInt() int {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	return random.Intn(10000)
}
