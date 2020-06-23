package ami

import "strconv"

func initEvent(data ActionData) Event {
	res := Event{
		ActionData: data,
	}
	if src, check := data["Uniqueid"]; check {
		res.uuid, _ = strconv.ParseInt(src, 10, 64)
	}
	return res
}

type Event struct {
	ActionData
	uuid int64
}

func (s Event) Name() string {
	return s.ActionData["Event"]
}

func (s Event) UUID() int64 {
	return s.uuid
}
