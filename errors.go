package tomeit

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func badRequestError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "リクエストに何らかの間違いがあります。",
		ErrorText:      err.Error(),
	}
}

func AuthenticationError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 401,
		StatusText:     "ログインする必要があります。",
		ErrorText:      err.Error(),
	}
}

func AuthorizationError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 403,
		StatusText:     "指定されたリソースに対してアクセス権限がありません。",
		ErrorText:      err.Error(),
	}
}

func NotFoundError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 404,
		StatusText:     "指定されたリソースが存在しません。",
		ErrorText:      err.Error(),
	}
}

func internalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "サーバ側で何らかのエラーが発生しました。",
		ErrorText:      err.Error(),
	}
}
