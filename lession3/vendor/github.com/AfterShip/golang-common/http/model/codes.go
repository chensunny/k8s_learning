package model

import "github.com/AfterShip/golang-common/errors"

type StatusCodeDescription struct {
	TypeName              string
	Message               string
	ErrorMessageFormatter func(desc StatusCodeDescription, err *errors.APIError) string
}

func AddStatusCodeDescriptions(descriptions map[int]StatusCodeDescription) {
	if descriptions == nil {
		return
	}
	for key, value := range descriptions {
		statusCodeDescriptions[key] = value
	}
}

func GetStatusCodeDescription(code int) StatusCodeDescription {
	return statusCodeDescriptions[code]
}

var statusCodeDescriptions = map[int]StatusCodeDescription{
	20000: {
		TypeName: "OK",
		Message:  "The request was successfully processed by AfterShip.",
	},
	20100: {
		TypeName: "Created",
		Message:  "The request has been fulfilled and a new resource has been created.",
	},
	20200: {
		TypeName: "Accepted",
		Message:  "The request has been accepted for processing, but the processing has not been completed.",
	},
	40000: {
		TypeName: "BadRequest",
		Message:  "The request was not understood by the server, generally due to bad syntax or because the Content-Type header was not correctly set to application/json.",
	},
	40100: {
		TypeName: "Unauthorized",
		Message:  "The necessary authentication credentials are not present in the request or are incorrect.",
	},
	40101: {
		TypeName: "Unauthorized",
		Message:  "The necessary authentication credentials are not present in the request or are incorrect.",
	},
	40102: {
		TypeName: "Unauthorized",
		Message:  "The necessary authentication credentials are not present in the request or are incorrect.",
	},
	40103: {
		TypeName: "Unauthorized",
		Message:  "The necessary authentication credentials are not present in the request or are incorrect.",
	},
	40200: {
		TypeName: "PaymentRequired",
		Message:  "Payment Required",
	},
	40300: {
		TypeName: "Forbidden",
		Message:  "The server is refusing to respond to the request. This is generally because you have not requested the appropriate scope for this action.",
	},
	40400: {
		TypeName: "NotFound",
		Message:  "The requested resource was not found but could be available again in the future.",
	},
	40500: {
		TypeName: "MethodNotAllowed",
		Message:  "The method received in the request-line is known by the server but not supported by the target resource.",
	},
	40900: {
		TypeName: "Conflict",
		Message:  "The request conflicts with another request (perhaps due to using the same idempotent key).",
	},
	42200: {
		TypeName: "UnprocessableEntity",
		Message:  "The request body was well-formed but contains semantical errors. The response body will provide more details in the errors or error parameters.",
	},
	42900: {
		TypeName: "TooManyRequests",
		Message:  "The request was not accepted because the application has exceeded the rate limit. The default API call limit is 10 requests per second.",
	},
	50000: {
		TypeName: "InternalError",
		Message:  "Something went wrong on AfterShip's end.  Also, some error that cannot be retried happened on an external system that this call relies on.",
	},
}
