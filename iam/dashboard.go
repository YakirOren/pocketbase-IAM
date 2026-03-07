package iam

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

//go:embed dashboard/*
var dashboardFS embed.FS

// registerDashboardRoutes serves the embedded IAM admin dashboard at /_/iam/.
// Access is restricted to superusers only.
func registerDashboardRoutes(app core.App) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		subFS, err := fs.Sub(dashboardFS, "dashboard")
		if err != nil {
			return err
		}

		indexHTML, err := fs.ReadFile(subFS, "index.html")
		if err != nil {
			return err
		}

		fileServer := http.StripPrefix("/_/iam/", http.FileServerFS(subFS))

		se.Router.GET("/_/iam/{path...}", func(e *core.RequestEvent) error {
			path := e.Request.PathValue("path")

			// Serve existing files via the standard file server.
			if path != "" {
				if _, err := fs.Stat(subFS, path); err == nil {
					if strings.HasPrefix(path, "assets/") {
						e.Response.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
					}
					fileServer.ServeHTTP(e.Response, e.Request)
					return nil
				}
			}

			// SPA fallback: serve index.html for unknown routes.
			e.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
			e.Response.WriteHeader(http.StatusOK)
			_, writeErr := e.Response.Write(indexHTML)
			return writeErr
		})

		return se.Next()
	})
}
