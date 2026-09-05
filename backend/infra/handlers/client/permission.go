package client

import (
	"net/http"

	"bitbucket.org/lyndus/backend/domain/client"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/utils"
)

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	if id != 0 {
		var sender = map[string]interface{}{
			"permissions": client.Permissions,
			"levels":      utils.AccessLevels,
		}

		config.JSONResponse(sender, http.StatusOK, w)
		return
	}

	config.ResponsePerErr(w, nil, config.INTERNALSERVERERROR)
	return
}

func GetByUserPermissions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	if id != 0 {
		var sender = map[string]interface{}{
			"permissions": client.Permissions,
			"levels":      utils.AccessLevels,
		}

		config.JSONResponse(sender, http.StatusOK, w)
		return
	}

	config.ResponsePerErr(w, nil, config.INTERNALSERVERERROR)
	return
}

