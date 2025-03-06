package common

import "errors"

var ErrUserFetchDetailsError = errors.New("error fetching user details")
var ErrEnvVarNotSet = errors.New("env var not set")
