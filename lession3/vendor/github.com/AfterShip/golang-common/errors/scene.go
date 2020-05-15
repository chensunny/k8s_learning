package errors

import "encoding/json"

//翻车现场，错误详情信息，会出现在logging里，便于问题排查
type Scene struct {
	//items
	items []interface{}
	//cause by
	cause error
	//debug.Stack()
	stack string
	// for custom fields
	fields map[string]interface{}
}

var emptyScene = &Scene{}

func (s *Scene) Items() []interface{} {
	return s.items
}

func (s *Scene) Cause() error {
	return s.cause
}

func (s *Scene) Stack() string {
	return s.stack
}

func (s *Scene) Fields() map[string]interface{} {
	return s.fields
}

func (s *Scene) JsonString() string {
	data, err := json.Marshal(s)
	if err == nil {
		return string(data)
	} else {
		return ""
	}
}

func (s *Scene) String() string {
	return s.JsonString()
}

type SceneField func(scene *Scene)

//引起error的源error, cause by
func Cause(cause error) SceneField {
	return func(scene *Scene) {
		scene.cause = cause
	}
}

//产生错误的源代码行stacktrace
func DefaultStack() SceneField {
	stacktrace := SmallerStacktrace(1, 1)
	return func(scene *Scene) {
		scene.stack = stacktrace
	}
}

func Stack(stack string) SceneField {
	return func(scene *Scene) {
		scene.stack = stack
	}
}

//错误明细列表，比如用于请求参数校验失败的时候，放非法参数明细
func Item(item interface{}) SceneField {
	return func(scene *Scene) {
		if scene.items == nil {
			scene.items = make([]interface{}, 0)
		}
		scene.items = append(scene.items, item)
	}
}

//错误明细列表，比如用于请求参数校验失败的时候，放非法参数明细
func Items(items ...interface{}) SceneField {
	return func(scene *Scene) {
		if scene.items == nil {
			scene.items = make([]interface{}, 0)
		}
		scene.items = append(scene.items, items...)
	}
}

//额外的错误详情字段信息
func Field(name string, value interface{}) SceneField {
	return func(scene *Scene) {
		if scene.fields == nil {
			scene.fields = make(map[string]interface{})
		}
		scene.fields[name] = value
	}
}
