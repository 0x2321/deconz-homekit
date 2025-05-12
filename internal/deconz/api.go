package deconz

type ApiClient struct {
	baseUrl string
	apiKey  string
}

func NewApiClient(baseUrl string, apiKey string) *ApiClient {
	return &ApiClient{
		baseUrl: baseUrl,
		apiKey:  apiKey,
	}
}

func (ac *ApiClient) buildUrl(path string) string {
	return ac.baseUrl + "/api/" + ac.apiKey + path
}
