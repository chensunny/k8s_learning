## Package types
```import "github.com/AfterShip/golang-common/lang/types"```

Package `types` implements many extended types, solved several issues below:

1. Make our API to adapt to newest [API Guideline](https://docs.google.com/spreadsheets/d/16oduZidE9ofdoT6m3I9oCNa9hPjgSCYmSeg3jxjofdo/edit#gid=722104700). Details see [here](#Directory structure) 
2. Make dynamic optional parameter into possible.
3. Make type extension more easier by exposing extended type interface.

### Directory structure
```
|-- types
|-- -- wrappedtype.go   # WrappedType Interface definition
|-- -- string.go        # extended string type
|-- -- int64.go         # extended int64 type
|-- -- foat64.go        # extended float64 type
|-- -- bool.go          # extended bool type
|-- -- bytes.go         # extended []byte type
|-- -- date.go          # extended civil.Date(used by google) type
|-- -- time.go          # extended time.Time type
|-- -- encoder.go       # jsoniter encoder
```

### Adapt to newest API Guideline

Due to our API Guideline, there are several situations below need to be processed correctly:
#### 1. JSON `null` can not be treated as Golang zero value
For example

A payload like this:
```json
{
  "user": {
    "name": null
  }
}
```
Assume we have a struct:
```go
package example

type Data struct {
    User User `json:"user"`
}

type User struct {
    Name string `json:"name"`
}
```
When we try to unmarshal this payload with struct `User`, we will get:
```
data.User.Name == ""
```
`null` was abandoned.

If we want to keep the information of `null` from payload, `User` struct should be changed like this:
```go
package example

type Data struct {
    User User `json:"user"`
}

type User struct {
    Name *string `json:"name"`
}
```
This time, the result of unmarshal will be:
```
data.User.Name == nil
```
Seems, except kinds of pointers exist at anywhere, everything is going on the right way? 

No, too young too simple, sometimes naive.

If a payload like this came to your API:
```json
{
  "user": {
    "age": 16
  }
}
```
You will get:
```
data.User.Name == nil
```
But this request probably means i only want to update my `age` in a HTTP PATCH.

And then, rely on your kinds of pointers, you will set his `name` with `null`, which is not expected.

So, we can have a conclusion:

For json attributes, actually they have three situation:
1. Have a value except `null`
2. `null`
3. Not Assigned

We can not find a way to satisfy those situation at the same time based on Golang native types(such as string, int, bool, float64 ...)

Using extended types under this package will solve this problem.

Assume our new `User` struct like this:
```go
package example
import "github.com/AfterShip/golang-common/lang/types"

type Data struct {
    User User `json:"user"`
}

type User struct {
    Name types.String `json:"name"`
    Age  types.Int64  `json:"age"`
}
``` 
When received:
```json
{
  "user": {
    "name": null
  }
}
```
We will get:
```
data.User.Name.Null() == true
data.User.Age.Assigned == false
```

#### 2. When all children attributes are nullable and definitely `null` at present, parent attribute must be `null`, vice versa
This rule can be explained in two ways. First, we focus on unmarshaling request payload.

Here is payload and struct definition:
```json
{
  "user": null
}
```
```go
package example
import "github.com/AfterShip/golang-common/lang/types"

type Data struct {
    User User `json:"user"`
}

type User struct {
    Name types.String `json:"name"`
    Age  types.Int64  `json:"age"`
}
```
Then, you will get:
```
data.User.Name == types.String{}
data.User.Age == types.Int64{}
```
No, it's not correct, what we want is `null`

Another change must be applied on your code, if want the expected result:

Use `gins.ShouldBindJSON`(from package `golang-common/http/server/gins`) instead of `c.ShouldBindJSON`.

Inside of this function, it use `autoSetNullExtension` of `jsoniter` to make sure all children attributes could be set to `null` automatically.

How about json marshaling?

Call `RegisterEncoder()` in `encoder.go` when your application start to run. It will register extended types to `jsoniter` globally.

Then, use `jsoniter.Marshal` instead of `json.Marshal`, everything will be ok now

### How to use extended types?
```go
package example
import (
    "github.com/AfterShip/golang-common/lang/types"
    "github.com/gin-gonic/gin/internal/json"
    "github.com/stretchr/testify/assert"
    "testing"
)

func main() {
    var t = new(testing.T)

    types.NewString("hello") // assigned string
    types.NewNullString()    // assigned null string
    _ = types.String{}       // not assigned string

    // you can access value of type, status of assigned or null
    _ = types.NewString("world").String()
    _ = types.NewNullString().Null()
    _ = types.NewString("hello").Assigned()

    var str = new(types.String)
    // you can copy String from another
    str.CopyFrom(types.NewString("hello"))
    assert.Equal(t, "hello", str.String())
    
    // you can define struct with extended types
    type User struct {
        Name types.String `json:"name"`
        Age  types.Int64  `json:"age"`
    }
    var user = new(User)
    err := json.Unmarshal([]byte(`{"name": null, "age": 14}`), &user)
    assert.Nil(t, err)
    assert.Equal(t, true, user.Name.Null())
    assert.Equal(t, int64(14), user.Age.Int64())
}
```