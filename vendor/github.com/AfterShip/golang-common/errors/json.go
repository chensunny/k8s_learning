package errors

import (
	"encoding/json"
)

//struct to json format
type apierrorForJson struct {
	MainCode *Code  `json:"main_code,omitempty"`
	SubCode  *Code  `json:"sub_code,omitempty"`
	Scene    *Scene `json:"scene,omitempty"`
}

//never used
func (e *APIError) UnmarshalJSON(data []byte) error {
	errorForJson := new(apierrorForJson)
	err := json.Unmarshal(data, errorForJson)
	if err != nil {
		return err
	}
	e.mainCode = errorForJson.MainCode
	e.subCode = errorForJson.SubCode
	e.scene = errorForJson.Scene
	return nil
}

func (e *APIError) MarshalJSON() (data []byte, err error) {
	errorForJson := &apierrorForJson{
		MainCode: e.mainCode,
		SubCode:  e.subCode,
		Scene:    e.scene,
	}
	return json.Marshal(errorForJson)
}

type codeForJson struct {
	Code       int    `json:"code,omitempty"`
	MessageKey string `json:"message_key,omitempty"`
}

func (c *Code) UnmarshalJSON(data []byte) error {
	code := new(codeForJson)
	err := json.Unmarshal(data, code)
	if err != nil {
		return err
	}
	c.code = code.Code
	c.messageKey = code.MessageKey
	return nil
}

func (c *Code) MarshalJSON() (data []byte, err error) {
	code := &codeForJson{
		Code:       c.code,
		MessageKey: c.messageKey,
	}
	return json.Marshal(code)
}

type errorWithSceneForJson struct {
	Cause error  `json:"cause,omitempty"`
	Scene *Scene `json:"scene,omitempty"`
}

//never used
func (e *ErrorWithScene) UnmarshalJSON(data []byte) error {
	errorWithSceneForJson := new(errorWithSceneForJson)
	err := json.Unmarshal(data, errorWithSceneForJson)
	if err != nil {
		return err
	}
	e.cause = errorWithSceneForJson.Cause
	e.scene = errorWithSceneForJson.Scene
	return nil
}

func (e *ErrorWithScene) MarshalJSON() (data []byte, err error) {
	var causeErr error
	if e.cause != nil {
		if _, ok := e.cause.(interface{ MarshalJSON() ([]byte, error) }); ok {
			causeErr = e.cause
		} else {
			switch e.cause.(type) {
			case *APIError, *ErrorWithScene:
				causeErr = e.cause
			default:
				causeErr = &errorMessage{Message: e.cause.Error()}
			}
		}
	}
	errorWithSceneForJson := &errorWithSceneForJson{
		Cause: causeErr,
		Scene: e.scene,
	}
	return json.Marshal(errorWithSceneForJson)
}

type sceneForJson struct {
	//items
	Items []interface{} `json:"items,omitempty"`
	//cause by
	Cause error `json:"cause,omitempty"`
	//debug.Stack()
	Stack string `json:"stack,omitempty"`
	// for custom fields
	Fields map[string]interface{} `json:"fields,omitempty"`
}

//never used
func (s *Scene) UnmarshalJSON(data []byte) error {
	scene := new(sceneForJson)
	err := json.Unmarshal(data, scene)
	if err != nil {
		return err
	}
	s.items = scene.Items
	s.cause = scene.Cause
	s.stack = scene.Stack
	s.fields = scene.Fields
	return nil
}

func (s *Scene) MarshalJSON() (data []byte, err error) {
	var errorCause error
	if s.cause != nil {
		if _, ok := s.cause.(interface{ MarshalJSON() ([]byte, error) }); ok {
			errorCause = s.cause
		} else {
			switch s.cause.(type) {
			case *APIError, *ErrorWithScene:
				errorCause = s.cause
			default:
				errorCause = &errorMessage{Message: s.cause.Error()}
			}
		}
	}
	scene := &sceneForJson{
		Items:  s.items,
		Cause:  errorCause,
		Stack:  s.stack,
		Fields: s.fields,
	}
	return json.Marshal(scene)
}

type errorMessage struct {
	Message string `json:"message"`
}

func (err *errorMessage) Error() string {
	return err.Message
}
