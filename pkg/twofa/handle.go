package twofa

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Handle struct {
	twoFaService *TwoFaService
}

func NewHandle(twoFaService *TwoFaService) Handle {
	return Handle{
		twoFaService: twoFaService,
	}
}

// Initiate sending 2fa code
// (POST /2fa:init)
func (h Handle) Post2faInit(w http.ResponseWriter, r *http.Request) *Response {
	var resp SuccessResponse

	data := &Post2faInitJSONRequestBody{}
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		return &Response{
			Code: http.StatusBadRequest,
			body: "unable to parse body",
		}
	}

	loginId, err := uuid.Parse(data.LoginID)
	if err != nil {
		return &Response{
			Code: http.StatusBadRequest,
			body: "invalid login id",
		}
	}

	err = h.twoFaService.InitTwoFa(r.Context(), loginId, data.TwofaType, data.Email)
	if err != nil {
		return &Response{
			Code: http.StatusInternalServerError,
			body: "failed to init 2fa: " + err.Error(),
		}
	}

	return Post2faInitJSON200Response(resp)
}
