package prettyjson

import "errors"

// ErrCreationFailed is returned when the handler cannot be constructed.
var ErrCreationFailed = errors.New("failed to create prettyjson handler")
