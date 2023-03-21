package client

import "errors"

var ErrCredentialsRequired = errors.New("both client_id and client_secret credentials must be specified")
