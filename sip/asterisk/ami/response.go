package ami

func initResponseError(err error) Response {
	return Response{
		ActionData{
			"Action":  "Error",
			"Message": err.Error(),
		},
	}
}

type Response struct {
	ActionData
}

func (s Response) IsError() bool {
	return s.ActionData["Response"] == "Error"
}

func (s Response) ErrorMessage() string {
	return s.ActionData["Message"]
}
