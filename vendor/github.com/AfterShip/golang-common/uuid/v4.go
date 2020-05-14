package uuid

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func GenerateUUIDV4() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}
