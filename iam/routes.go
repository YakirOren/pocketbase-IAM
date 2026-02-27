package iam

import (
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func registerRoutes(app core.App, cache *PolicyCache) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.POST("/api/iam/check", func(e *core.RequestEvent) error {
			var body struct {
				Action   string `json:"action"`
				Resource string `json:"resource"`
			}
			if err := e.BindBody(&body); err != nil {
				return e.BadRequestError("invalid request body", err)
			}

			if body.Action == "" {
				return e.BadRequestError("action is required", nil)
			}
			if body.Resource == "" {
				body.Resource = "*"
			}

			userID := e.Auth.Id

			allowed, reason, err := Evaluate(app, cache, userID, body.Action, body.Resource)
			if err != nil {
				app.Logger().Error("IAM evaluation error",
					"user", userID,
					"action", body.Action,
					"resource", body.Resource,
					"error", err,
				)
				return e.InternalServerError("internal error", nil)
			}

			if !allowed {
				app.Logger().Warn("IAM access denied",
					"user", userID,
					"action", body.Action,
					"resource", body.Resource,
					"reason", reason,
				)
			}

			return e.JSON(200, map[string]any{"allowed": allowed})
		}).Bind(apis.RequireAuth())

		return se.Next()
	})
}
