package model

import (
	"fmt"
	"strconv"

	"github.com/AfterShip/golang-common/errors"
)

// see https://docs.google.com/spreadsheets/d/16oduZidE9ofdoT6m3I9oCNa9hPjgSCYmSeg3jxjofdo/edit#gid=524712657
type ResponseBody struct {
	Meta ResponseMeta `json:"meta"`
	Data interface{}  `json:"data"`
}

// Common response meta.
//
// swagger:model
type ResponseMeta struct {
	// Response code
	//
	// required: true
	// example: 20000
	Code int `json:"code"`
	// Response status
	//
	// required: true
	// example: OK
	Type string `json:"type,omitempty"`
	// Message string
	Message string `json:"message,omitempty"`
	// Error
	Errors []interface{} `json:"errors,omitempty"`
}

func BuildMetaCode(mainCode, subCode int) int {
	code, _ := strconv.Atoi(fmt.Sprintf("%d%0*d", mainCode, 2, subCode))
	return code
}

func BuildResponseMeta(metaCode int, err *errors.APIError) ResponseMeta {
	desc := GetStatusCodeDescription(metaCode)
	message := desc.Message
	if desc.ErrorMessageFormatter != nil && err != nil {
		message = desc.ErrorMessageFormatter(desc, err)
	}
	return ResponseMeta{
		Code:    metaCode,
		Type:    desc.TypeName,
		Message: message,
	}
}
