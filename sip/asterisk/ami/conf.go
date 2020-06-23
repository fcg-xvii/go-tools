package ami

import "time"

const (
	max_client_uuid = ^uint64(0)
)

var (
	RequestTimeoutDefault = time.Second * 20
)
