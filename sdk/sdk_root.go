package sdk

// Sdk Object to interact with various MDS Cloud resources
type Sdk struct {
	qsURL          string
	smURL          string
	fsURL          string
	nsURL          string
	sfURL          string
	defaultAccount string
}

// NewSdk Creates a new SDK object
//
// account - The account that clients will act against
//
// urls    - Key value map for various clients to act against.
//   qsUrl - queue service url
//   smUrl - state machine service url
//   fsUrl - file service url
//   nsUrl - notification service url
//   sfUrl - serverless function service url
func NewSdk(account string, urls map[string]string) *Sdk {
	// TODO: document parameters
	sdk := Sdk{
		qsURL: urls["qsUrl"],
		smURL: urls["smUrl"],
		fsURL: urls["fsUrl"],
		nsURL: urls["nsUrl"],
		sfURL: urls["sfUrl"],
	}
	sdk.defaultAccount = account
	return &sdk
}

// GetServerlessFunctionsClient Gets a new serverless function client
func (s *Sdk) GetServerlessFunctionsClient(sfURL string, defaultAccount string) *ServerlessFunctionsClient {
	var url, account string
	if sfURL != "" {
		url = sfURL
	} else {
		url = s.sfURL
	}

	if defaultAccount != "" {
		account = defaultAccount
	} else {
		account = s.defaultAccount
	}

	client := ServerlessFunctionsClient{
		serviceURL:     url,
		defaultAccount: account,
	}
	return &client
}
