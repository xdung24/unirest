package service

import (
	"errors"

	"github.com/gorilla/mux"
	"github.com/r3labs/sse/v2"
)

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

type BrokerEvent struct {
	Event     string      `json:"event"`
	User      string      `json:"user_id,omitempty"`
	Namespace string      `json:"namespace"`
	Key       string      `json:"key,omitempty"`
	Value     interface{} `json:"value,omitempty"`
}

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
