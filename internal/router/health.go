package router

import "net/http"

func (s *service) HealthRoute(w http.ResponseWriter, r *http.Request) error {
	if err := s.db.HealthCheck(); err != nil {
		return err
	}

	return JSONResponse(w, http.StatusOK, map[string]any{"status": "ok"})
}
