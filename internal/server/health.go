package server

import "net/http"

func (api *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	data := map[string]string{
		"status":  "ok",
		"env":     api.config.Env,
		"version": api.config.Version,
	}

	if err := api.jsonResponse(w, http.StatusOK, data); err != nil {
		api.internalServerError(w, r, err)
		return
	}
}
