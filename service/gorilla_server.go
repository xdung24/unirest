package service

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/r3labs/sse/v2"
	"github.com/xeipuuv/gojsonschema"
)

func (s *Server) Init(db Database) {
	s.db = db
	s.db.Init()

	s.router = mux.NewRouter()

	s.router.HandleFunc(NamespaceHomePattern, s.homeHandler)
	s.router.HandleFunc(NamespacePattern, s.namespaceHandler).Methods(http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions)

	s.router.HandleFunc(DataSetPattern, s.dataSetHandler).Methods(http.MethodGet, http.MethodDelete, http.MethodOptions)
	s.router.HandleFunc(DataSetKeyValuePattern, s.dataSetKeyValueHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions)

	s.router.HandleFunc(SearchPattern, s.searchHandler).Queries("filter", "{filter}")
	s.router.HandleFunc(SchemaPattern, s.schemaHandler)

	if s.SwaggerEnabled {
		s.router.HandleFunc(OpenAPIPattern, s.openAPIHandler)
		s.router.PathPrefix(SwaggerUIPattern).Handler(http.StripPrefix(SwaggerUIPattern, http.FileServer(http.Dir("./swagger-ui/"))))
		log.Println("swagger extension enabled")

	}

	if s.BrokerEnabled {
		sseServer := sse.New()
		sseServer.EventTTL = time.Second * 15 // keep message alive for 15 seconds, so the client can reconnect

		sseServer.CreateStream("messages")
		brokerServer := http.NewServeMux()
		brokerServer.HandleFunc(BrokerPattern, func(w http.ResponseWriter, r *http.Request) {
			// Send a heartbeat every 30 seconds to keep the connection alive
			ticker := time.NewTicker(30 * time.Second)
			go func() {
				for {
					select {
					case <-ticker.C:
						// Send a comment as a heartbeat, clients will ignore this
						fmt.Fprintf(w, ":heartbeat\n\n")
						w.(http.Flusher).Flush()
					case <-r.Context().Done():
						ticker.Stop()
						return
					}
				}
			}()

			sseServer.ServeHTTP(w, r)
		})
		s.broker = sseServer

		go func() {
			log.Println("broker server started at: " + s.BrokerAddress)
			http.ListenAndServe(s.BrokerAddress, sseServer)
		}()

		log.Println("broker extension enabled")
	}

	s.router.Use(mux.CORSMethodMiddleware(s.router))

	if s.AuthEnabled {
		verifyBytes, err := os.ReadFile(certsPublicKey)
		if err != nil {
			log.Fatalf("auth required but error on reading public key for JWT: %v", err)
		}
		middleware := JWTAuthMiddleware{
			VerifyBytes: verifyBytes,
		}
		s.router.Use(middleware.GetMiddleWare(s.router))
		log.Println("authentication middleware enabled")
	}

	srv := &http.Server{
		Handler: handlers.CompressHandlerLevel(s.router, gzip.BestSpeed),
		// Handler:      s.router,
		Addr:         s.Address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (s *Server) validate(namespace string, data []byte) (interface{}, error) {
	var parsed interface{}

	// if namespace has a schema, validate against it
	schemaJson, dbErr := s.db.Get(namespace+SchemaId, SchemaId)
	if dbErr == nil {
		schemaLoader := gojsonschema.NewBytesLoader(schemaJson)
		documentLoader := gojsonschema.NewBytesLoader(data)

		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			return nil, err
		}

		if result.Valid() {
			json.Unmarshal(data, &parsed)
		} else {
			log.Printf("The document is not valid according to its schema. see errors :")
			errorLog := ""
			for _, desc := range result.Errors() {
				errorLog += desc.String()
			}
			log.Println(errorLog)
			return nil, errors.New(errorLog)
		}
	} else {
		// otherwise just validate as json
		err := json.Unmarshal(data, &parsed)
		if err != nil {
			log.Println("The document is not valid JSON")
			return nil, err
		}
	}
	return parsed, nil
}

func (s *Server) Notify(event BrokerEvent) {
	if s.broker != nil {
		jsonData, _ := json.Marshal(event)
		s.broker.Publish("messages", &sse.Event{
			Data: []byte(jsonData),
		})
	}
}
