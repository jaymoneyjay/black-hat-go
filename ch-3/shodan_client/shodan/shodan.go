package shodan

const BaseURL = "https://api.shodan.io"

// Client stores information to access the api
type Client struct {
	apiKey string
}

// New creates an instance of Client
func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}
