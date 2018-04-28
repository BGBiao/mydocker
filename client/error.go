package client
import (
    "errors"
)

var (
    ErrConn  = errors.New("user not exist")
    ErrInvalidPasswd = errors.New("Passwd or username not right")
    ErrInvalidParams = errors.New("Invalid params")
    ErrUserExist     = errors.New("user exist")
)
