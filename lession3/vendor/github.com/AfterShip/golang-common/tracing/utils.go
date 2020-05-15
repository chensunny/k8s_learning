package tracing

import "github.com/AfterShip/golang-common/uuid"

func GenerateTracingID() string {
	return uuid.GenerateUUIDV4() + "/0"
}
