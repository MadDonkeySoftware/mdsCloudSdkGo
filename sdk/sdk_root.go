package sdk

// Sdk Object to interact with various MDS Cloud resources
type Sdk struct {
	qsURL              string
	smURL              string
	fsURL              string
	nsURL              string
	sfURL              string
	defaultAccount     string
	defaultAuthManager *AuthManager
}

// NewSdk Creates a new SDK object
//
// account       - The account that clients will act against
// userId        - The user id used during authentication
// password      - The password used during authentication
// allowSelfCert - Allow HTTPS authentication when self-signed certificate used
//
// urls          - Key value map for various clients to act against.
//   qsUrl       - queue service url
//   smUrl       - state machine service url
//   fsUrl       - file service url
//   nsUrl       - notification service url
//   sfUrl       - serverless function service url
func NewSdk(account string, userID string, password string, allowSelfCert bool, urls map[string]string) *Sdk {
	// TODO: document parameters
	authManager := NewAuthManager(
		urls["identityUrl"],
		userID,
		password,
		account,
		allowSelfCert,
	)
	sdk := Sdk{
		qsURL: urls["qsUrl"],
		smURL: urls["smUrl"],
		fsURL: urls["fsUrl"],
		nsURL: urls["nsUrl"],
		sfURL: urls["sfUrl"],
	}
	sdk.defaultAccount = account
	sdk.defaultAuthManager = authManager
	return &sdk
}

// GetServerlessFunctionsClient Gets a new serverless function client
func (s *Sdk) GetServerlessFunctionsClient() *ServerlessFunctionsClient {
	client := ServerlessFunctionsClient{
		serviceURL:     s.sfURL,
		defaultAccount: s.defaultAccount,
		authManager:    s.defaultAuthManager,
	}
	return &client
}
