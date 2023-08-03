package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itchyny/gojq"
)

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	result := struct {
		Results []interface{} `json:"results"`
	}{
		Results: make([]interface{}, 0),
	}

	switch r.Method {
	case http.MethodGet:
		vars := mux.Vars(r)
		query, err := gojq.Parse(vars["filter"])
		if err != nil {
			log.Println("error on parsing", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		data, dbErr := s.db.GetAll(vars["namespace"])
		if dbErr != nil {
			log.Println("error on GetAll", err)
			respondWithError(w, http.StatusBadRequest, dbErr.Error())
			return
		}
		for key, value := range data {
			var jsonContent map[string]interface{}
			err := json.Unmarshal(value, &jsonContent)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
			iter := query.Run(jsonContent)
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}
				if err, ok := v.(error); ok {
					log.Println("error on query", err)
					respondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}
				result.Results = append(result.Results, map[string]interface{}{"key": key, "value": v})
			}
		}
		jsonResponse, _ := json.Marshal(result)
		respondWithJSON(w, http.StatusOK, string(jsonResponse))
	}
}
