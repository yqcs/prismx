package server

// RegisterRequest is a request for client registration to interactsh server.
type RegisterRequest struct {
	// PublicKey is the public RSA Key of the client.
	PublicKey string `json:"public-key"`
	// SecretKey is the secret-key for correlation ID registered for the client.
	SecretKey string `json:"secret-key"`
	// CorrelationID is an ID for correlation with requests.
	CorrelationID string `json:"correlation-id"`
}

// DeregisterRequest is a request for client deregistration to interactsh server.
type DeregisterRequest struct {
	// CorrelationID is an ID for correlation with requests.
	CorrelationID string `json:"correlation-id"`
	// SecretKey is the secretKey for the interactsh client.
	SecretKey string `json:"secret-key"`
}

// PollResponse is the response for a polling request
type PollResponse struct {
	Data    []string `json:"data"`
	Extra   []string `json:"extra"`
	AESKey  string   `json:"aes_key"`
	TLDData []string `json:"tlddata,omitempty"`
}
