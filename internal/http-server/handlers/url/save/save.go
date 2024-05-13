package save

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortner/internal/lib/response"
)

// todo move this to configuration
const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type UrlSaver interface {
	SaveURL(alias, urlToSave string) (int64, error)
}

func New(log *slog.Logger, urlSaver UrlSaver) http.Handlerfunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("op", "handlers.url.save.New"),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to parse request body", err)
			render.JSON(w, r, response.Error("failed to parse request body"))
			return
		}
		log.Info("request body parsed", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("request validation failed", err)

			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
	}

}
