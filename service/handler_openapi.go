package service

import (
	"encoding/json"
	"net/http"
)

func (s *Server) openAPIHandler(w http.ResponseWriter, r *http.Request) {
	namespaces := s.db.GetNamespaces()

	rootMap, err := s.generateOpenAPIMap(namespaces)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		output, err := json.MarshalIndent(rootMap, "", "  ")
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		output = append(output, '\n')

		respondWithJSON(w, http.StatusOK, string(output))
	case http.MethodPost:
		respondWithError(w, http.StatusNotImplemented, "cannot POST to this endpoint!")
	case http.MethodDelete:
		respondWithError(w, http.StatusNotImplemented, "cannot DELETE this endpoint!")
	}
}
