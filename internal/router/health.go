package router

import "net/http"

func HealthRoute(w http.ResponseWriter, r *http.Request) error {
	return JSONResponse(w, http.StatusOK, map[string]any{"status": "ok"})
}
