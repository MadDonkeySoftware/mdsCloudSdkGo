package sdk

// Sdk Object to interact with various MDS Cloud resources
type Sdk struct {
	identityURL         string
	qsURL               string
	smURL               string
	fsURL               string
	nsURL               string
	sfURL               string
	defaultAccount      string
	defaultAuthManager  *AuthManager
	allowSelfCert       bool
	enableAuthSemaphore bool
}

// NewSdk Creates a new SDK object
//
// account       - The account that clients will act against
// userId        - The user id used during authentication
// password      - The password used during authentication
// allowSelfCert - Allow HTTPS authentication when self-signed certificate used
//
// urls          - Key value map for various clients to act against.
//   identityUrl - identity service url
//   qsUrl       - queue service url
//   smUrl       - state machine service url
//   fsUrl       - file service url
//   nsUrl       - notification service url
//   sfUrl       - serverless function service url
func NewSdk(account string, userID string, password string, allowSelfCert bool, enableAuthSemaphore bool, urls map[string]string) *Sdk {
	// TODO: document parameters
	authManager := NewAuthManager(
		urls["identityUrl"],
		userID,
		password,
		account,
		allowSelfCert,
		enableAuthSemaphore,
	)
	sdk := Sdk{
		identityURL: urls["identityUrl"],
		qsURL:       urls["qsUrl"],
		smURL:       urls["smUrl"],
		fsURL:       urls["fsUrl"],
		nsURL:       urls["nsUrl"],
		sfURL:       urls["sfUrl"],
	}
	sdk.defaultAccount = account
	sdk.defaultAuthManager = authManager
	sdk.allowSelfCert = allowSelfCert
	sdk.enableAuthSemaphore = enableAuthSemaphore
	return &sdk
}

// GetServerlessFunctionsClient Gets a new serverless function client
func (s *Sdk) GetServerlessFunctionsClient() *ServerlessFunctionsClient {
	return &ServerlessFunctionsClient{
		serviceURL:  s.sfURL,
		authManager: s.defaultAuthManager,
	}
}

// GetIdentityClient Gets a new identity client
func (s *Sdk) GetIdentityClient() *IdentityClient {
	return &IdentityClient{
		allowSelfSignCert: s.allowSelfCert,
		authManager:       s.defaultAuthManager,
		identityURL:       s.identityURL,
	}
}

// GetQueueServiceClient Gets a new queue service client
func (s *Sdk) GetQueueServiceClient() *QueueServiceClient {
	return &QueueServiceClient{
		authManager:     s.defaultAuthManager,
		queueServiceURL: s.qsURL,
	}
}

// GetFileServiceClient Gets a new file service client
func (s *Sdk) GetFileServiceClient() *FileServiceClient {
	return &FileServiceClient{
		authManager:    s.defaultAuthManager,
		fileServiceURL: s.fsURL,
	}
}

// GetStateMachineServiceClient Gets a new state machine service client
func (s *Sdk) GetStateMachineServiceClient() *StateMachineServiceClient {
	return &StateMachineServiceClient{
		authManager:            s.defaultAuthManager,
		stateMachineServiceURL: s.smURL,
	}
}
