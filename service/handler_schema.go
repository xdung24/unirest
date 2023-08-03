package service

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) schemaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"] + SchemaId

	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		dbErr := s.db.Upsert(namespace, SchemaId, data, true)
		if dbErr != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		log.Printf("added schema for namespace '%s'\n", vars["namespace"])
		respondWithJSON(w, http.StatusCreated, string(data))
	case http.MethodGet:
		data, dbErr := s.db.Get(namespace, SchemaId)
		if dbErr != nil {
			respondWithError(w, http.StatusNotFound, dbErr.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, string(data))
	case http.MethodDelete:
		dbErr := s.db.Delete(namespace, SchemaId)
		if dbErr != nil {
			respondWithError(w, http.StatusNotFound, dbErr.Error())
			return
		}
		respondWithJSON(w, http.StatusAccepted, "{}")
	}
}
