package types

import "encoding/json"

type SocketActionHandler func(json.RawMessage) (bool, any, error)
