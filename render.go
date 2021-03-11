package zerror

type Resetter interface {
	Reset()
}
type Render interface {
	SetCode(code string)
	SetMessage(msg string)
	Error() string
}

type StdResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (r *StdResponse) SetCode(code string) {
	r.Code = code
}

func (r *StdResponse) SetMessage(msg string) {
	r.Msg = msg
}

func (r *StdResponse) Error() string {
	return r.Msg
}
