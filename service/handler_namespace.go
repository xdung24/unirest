package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xdung24/universal-rest/database"
)

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	namespaces, err := jsonWrapper(s.db.GetNamespaces())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, string(namespaces))
}

func (s *Server) namespaceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	userId := r.Header.Get(USER_HEADER)

	vars := mux.Vars(r)
	namespace := vars["namespace"]

	switch r.Method {
	case http.MethodPost:
		dbErr := s.db.CreateNameSpace(namespace)
		if dbErr != nil {
			respondWithError(w, http.StatusInternalServerError, dbErr.Error())
		}
		s.Notify(BrokerEvent{
			Event:     EVENT_NAMESPACE_CREATED,
			User:      userId,
			Namespace: namespace,
			Key:       "",
			Value:     nil,
		})
		respondWithJSON(w, http.StatusCreated, "{}")
	case http.MethodGet:
		data, dbErr := s.db.GetAll(namespace)
		if dbErr != nil {
			switch dbErr.ErrorCode {
			case database.NAMESPACE_NOT_FOUND:
				respondWithError(w, http.StatusBadRequest, dbErr.Error())
			default:
				respondWithError(w, http.StatusInternalServerError, dbErr.Error())
			}
		}
		namespaceData, err := jsonWrapper(data)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, string(namespaceData))
	case http.MethodDelete:
		dbErr := s.db.DropNameSpace(namespace)
		if dbErr != nil {
			switch dbErr.ErrorCode {
			case database.NAMESPACE_NOT_FOUND:
				respondWithError(w, http.StatusBadRequest, dbErr.Error())
			default:
				respondWithError(w, http.StatusInternalServerError, dbErr.Error())
			}
		}
		s.Notify(BrokerEvent{
			Event:     EVENT_NAMESPACE_DELETED,
			User:      userId,
			Namespace: namespace,
			Key:       "",
			Value:     nil,
		})
		respondWithJSON(w, http.StatusAccepted, "{}")
	}
}
