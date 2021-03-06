package covid

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// DefaultBasepath is the basepath used to make requets, unless modified when
// initlizing the client.
const DefaultBasepath = "https://api.coronavirus.data.gov.uk"

// Client is used to make requests to the gov.uk coronavirus API.
type Client struct {
	basepath   string
	httpClient *http.Client
}

// Option[s] are used to configure the covid.Client.
type Option func(*Client)

// Filter[s] are used to control the data retrieved from the GET /data endpoint.
type Filter interface {
	MetricName() string
	Value() string
}

// Filters is a typedef for a slice of Filter[s] to allow methods on the slice.
type Filters []Filter

// AsParameter formats a query parameter from the filters, which is necessary as
// the API doesn't follow standard array query convention.
func (filters Filters) AsParameter() string {
	param := ""
	for i, filter := range filters {
		if i != 0 {
			param = param + ";"
		}
		param = param + fmt.Sprintf("%s=%s", filter.MetricName(), filter.Value())
	}
	return param
}

// Format indicates the format of data in the response body.
type Format string

const (
	FormatJSON = "json"
	FormatCSV  = "csv"
	FormatXML  = "xml"
)

// NewClient initializes a Client, applying the options passed in order.
func NewClient(opts ...Option) *Client {
	// Initialize with the default basepath and http client.
	client := &Client{
		basepath:   DefaultBasepath,
		httpClient: &http.Client{},
	}
	// Initialize as a V1 client by default.
	V1(client)
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// V1 sets up the client to use the 'v1' path.
var V1 = Option(func(client *Client) {
	client.basepath = fmt.Sprintf("%s/v1", client.basepath)
})

// WithBasepath configures the client to use the basepath provided when making
// requests.
func WithBasepath(basepath string) Option {
	return func(client *Client) {
		client.basepath = basepath
	}
}

// WithHTTPClient configures the client to use the http client provided when
// making requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

type response struct {
	Length       int             `json:"length"`
	MaxPageLimit int             `json:"maxPageLimit"`
	Data         json.RawMessage `json:"data"`
	Pagination   *pagination     `json:"pagination"`
}

type pagination struct {
	First    *string `json:"first"`
	Previous *string `json:"previous"`
	Current  *string `json:"current"`
	Next     *string `json:"next"`
	Last     *string `json:"last"`
}

// Data returns a io.Reader for the fetched given the parameters provided.
func (c *Client) GetData(ctx context.Context, structure json.Marshaler, format Format, areaType *AreaType, filters ...Filter) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/data", c.basepath), nil)
	if err != nil {
		return nil, err
	}
	// Add the required headers.
	req.Header.Add("Accepts", "application/json; application/xml; text/csv; application/vnd.PHE-COVID19.v1+json; application/vnd.PHE-COVID19.v1+xml")
	req.Header.Add("Content-Type", "application/json")
	// Setup query parameters.
	q := req.URL.Query()
	q.Add("format", string(format))
	marshaled, err := structure.MarshalJSON()
	if err != nil {
		return nil, err
	}
	q.Add("structure", string(marshaled))
	q.Add("filters", Filters(append(filters, areaType)).AsParameter())
	req.URL.RawQuery = q.Encode()
	// Make the request.
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	switch format {
	case FormatJSON:
		body := &response{}
		if err := json.NewDecoder(res.Body).Decode(body); err != nil {
			return nil, err
		}
		res.Body.Close()
		return io.NopCloser(bytes.NewBuffer(body.Data)), nil
	case FormatCSV:
		return res.Body, nil
	default:
		res.Body.Close()
		return nil, fmt.Errorf("format '%s' is unsupported", format)
	}
}
