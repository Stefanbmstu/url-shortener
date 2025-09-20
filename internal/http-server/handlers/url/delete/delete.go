package delete

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

// URLDeleter is an interface for deleting url by alias.
//
//go:generate mockery --name=URLDeleter --output=../../mocks --outpkg=mocks
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request: alias is required"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found for deletion", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("url deleted", slog.String("alias", alias))

		render.JSON(w, r, resp.OK())
	}
}
