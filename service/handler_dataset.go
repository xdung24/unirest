package service

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xdung24/universal-rest/database"
)

func (s *Server) dataSetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	userId := r.Header.Get(USER_HEADER)

	vars := mux.Vars(r)
	namespace := vars["namespace"]

	switch r.Method {
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
		format := r.URL.Query().Get("format")
		switch format {
		case "", "1":
			namespaceData, err := jsonWrapper(data)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondWithJSON(w, http.StatusOK, string(namespaceData))
		case "2":
			namespaceData, err := jsonWrapper2(data)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondWithJSON(w, http.StatusOK, string(namespaceData))
		case "3":
			namespaceData, err := jsonWrapper3(data)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			respondWithJSON(w, http.StatusOK, string(namespaceData))
		default:
			respondWithError(w, 400, "Invalid query")
		}

	case http.MethodDelete:
		dbErr := s.db.DeleteAll(namespace)
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

func (s *Server) dataSetKeyValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	userId := r.Header.Get(USER_HEADER)

	vars := mux.Vars(r)
	namespace := vars["namespace"]
	key := vars["key"]

	switch r.Method {
	case http.MethodPost, http.MethodPut:
		defer r.Body.Close()
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		data, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		_onUpsert(s, w, r.Method, userId, namespace, key, data)
	case http.MethodGet:
		data, dbErr := s.db.Get(namespace, key)
		if dbErr != nil {
			switch dbErr.ErrorCode {
			case database.ID_NOT_FOUND:
				respondWithError(w, http.StatusNotFound, dbErr.Error())
			case database.NAMESPACE_NOT_FOUND:
				respondWithError(w, http.StatusBadRequest, dbErr.Error())
			default:
				respondWithError(w, http.StatusInternalServerError, dbErr.Error())
			}
			return
		}
		respondWithJSON(w, http.StatusOK, string(data))
	case http.MethodDelete:
		err := s.db.Delete(namespace, key)
		if err != nil {

			switch err.ErrorCode {
			case database.ID_NOT_FOUND:
				respondWithError(w, http.StatusNotFound, err.Error())
			case database.NAMESPACE_NOT_FOUND:
				respondWithError(w, http.StatusBadRequest, err.Error())
			default:
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		s.Notify(BrokerEvent{
			Event:     EVENT_ITEM_DELETED,
			User:      userId,
			Namespace: namespace,
			Key:       key,
			Value:     nil,
		})
		respondWithJSON(w, http.StatusAccepted, "{}")
	}
}

// both POST and PUT methods will create new item
// POST will reject updating record while PUT will update record when existing
func _onUpsert(s *Server, w http.ResponseWriter, method, userId, namespace, key string, data []byte) {
	parsedData, err := s.validate(namespace, data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if s.AuthEnabled {
		// override data with a payload
		payload := Payload{
			User: userId,
			Data: parsedData,
		}
		data, err = payload.wrap()

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	allowOverWrite := false
	if method == http.MethodPut {
		allowOverWrite = true
	}

	dbErr := s.db.Upsert(namespace, key, data, allowOverWrite)
	if dbErr != nil {
		switch dbErr.ErrorCode {
		case database.NAMESPACE_NOT_FOUND:
			respondWithError(w, http.StatusBadRequest, dbErr.Error())
		case database.ITEM_CONFLICT:
			respondWithError(w, http.StatusConflict, dbErr.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, dbErr.Error())
		}
		return
	}

	event := EVENT_ITEM_CREATED
	if method == http.MethodPut {
		event = EVENT_ITEM_UPDATED
	}

	s.Notify(BrokerEvent{
		Event:     event,
		User:      userId,
		Namespace: namespace,
		Key:       key,
		Value:     parsedData,
	})
	respondWithJSON(w, http.StatusCreated, string(data))
}
