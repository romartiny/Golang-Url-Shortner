package url

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang-url-shortner/internal/lib/api/response"
	"golang-url-shortner/internal/lib/logger/sl"
	"net/http"

	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Status string `json:"status"`          //error, ok
	Error  string `json:"error,omitempty"` //omitempty for ignore empty errors
	Alias  string `json:"alias,omitempty"` //1 12 40
}

type URLSaver interface {
	Save(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context()))),
		)
	    var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err!= nil {
            log.Error("failed to decode request body", sl.Err(err))

            render.JSON(w, r, response.Error("failed to decode request"))

            return
        }
	}
}
