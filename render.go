package zerror

type Responser interface {
	SetCode(code string)
	SetMessage(msg string)
}

type StdResponse struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
	Msg  *string     `json:"msg"`
}

func (r *StdResponse) SetCode(code string) {
	r.Code = code
}

func (r *StdResponse) SetMessage(msg string) {
	r.Msg = &msg
}
