# Specification

* APIError for api layer
* Wrap error for any layer
* Use xerrors.Is(err,causedByErr) to detect source error reason

# Example

```
var (
   causeByErr = ${some error}
)

func doSomethingInAnyLayer() error{
    receivedErr:= doOther() 

    if (xerrors.Is(receivedErr,causeByErr){
        // do some error handling
    }
    //wrap err
    errors.Wrap(receivedErr,errerrors.Field("account_id","xxxxxx"))
}

func doSomethingInAPILayer() error{
    receivedErr:=doSomethingInAnyLayer()
    if receivedErr!=nil{
        if (xerrors.Is(receivedErr,causeByErr){
            // do some error handling and wrap error with code
            return errors.WithScene(errors.ErrConflict,errors.Cause(causeByErr)
        }
        // do some error handling and wrap error with code
        return errors.WithScene(errors.ErrInternalError,errors.Cause(causeByErr)
    }
    return nil
}

func handleRequest(ginCtx *gin.Context){
    err:= doSomethingInAPILayer()
    if err!=nil{
        handleResponseError(ginCtx,err)
    }
    // do response okay
}

func handleResponseError(ginCtx *gin.Context,err error){
    if (xerrors.Is(err,errors.ErrInternalError){
        // do something
    }

    if (xerrors.Is(err,${some error}){
        // do something
    }
}

```

