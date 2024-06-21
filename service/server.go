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

	"github.com/xdung24/universal-rest/database"
)

type Database interface {
	Init()
	Disconnect()
	CreateNameSpace(namespace string) *database.DbError
	GetNamespaces() []string
	DropNameSpace(namespace string) *database.DbError
	Upsert(namespace string, key string, value []byte, allowOverWrite bool) *database.DbError
	Get(namespace string, key string) ([]byte, *database.DbError)
	GetAll(namespace string) (map[string][]byte, *database.DbError)
	Delete(namespace string, key string) *database.DbError
	DeleteAll(namespace string) *database.DbError
}

const (
	NamespaceHomePattern   = "/namespace"
	NamespacePattern       = "/namespace/{namespace:[a-zA-Z0-9\\-]+}"
	DataSetPattern         = "/dataset/{namespace:[a-zA-Z0-9\\-]+}"
	DataSetKeyValuePattern = "/dataset/{namespace:[a-zA-Z0-9\\-]+}/{key:[a-zA-Z0-9\\-]+}"
	SearchPattern          = "/search/{namespace:[a-zA-Z0-9\\-]+}"
	SchemaPattern          = "/schema/{namespace:[a-zA-Z0-9\\-]+}"
	OpenAPIPattern         = "/{openapi|swagger}.json"
	BrokerPattern          = "/broker"
	SwaggerUIPattern       = "/swaggerui/"

	SchemaId = "_schema"

	EVENT_ITEM_CREATED = "ITEM_CREATED"
	EVENT_ITEM_UPDATED = "ITEM_UPDATED"
	EVENT_ITEM_DELETED = "ITEM_DELETED"

	EVENT_NAMESPACE_CREATED = "NAMESPACE_CREATED"
	EVENT_NAMESPACE_DELETED = "NAMESPACE_DELETED"

	EVENT_SCHEMA_CREATED = "SCHEMA_CREATED"
	EVENT_SCHEMA_DELETED = "SCHEMA_DELETED"

	certsPublicKey = "./certs/public-cert.pem"
)

var (
	ErrInvalidArguments = errors.New("invalid arguments")
)

type Server struct {
	Address        string
	BrokerAddress  string
	SwaggerEnabled bool
	BrokerEnabled  bool
	AuthEnabled    bool
	RawSqlEnabled  bool

	router *mux.Router
	broker *sse.Server
	db     Database
}

type BrokerEvent struct {
	Event     string      `json:"event"`
	User      string      `json:"user_id,omitempty"`
	Namespace string      `json:"namespace"`
	Key       string      `json:"key,omitempty"`
	Value     interface{} `json:"value,omitempty"`
}

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
