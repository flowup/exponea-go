package exponea

import (
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"strings"
)

const (
	DefaultAPIEndpoint = "https://api.exponea.com/"
	EventsEndpoint = "crm/events"
	CustomersEndpoint = "crm/customers"
	BulkEndpoint = "bulk"
)

// Event is an encapsulation for event sending request
type Event struct {
	Customers  map[string]string `json:"customer_ids"`
	ProjectID  string `json:"project_id"`
	Type       string `json:"type"`
	Properties map[string]string `json:"properties"`
	Timestamp  int `json:"timestamp,omitempty"`
}

// Customer is an encapsulation for customer value updating
// request
type Customer struct {
	Customers  map[string]string `json:"ids"`
	ProjectID  string `json:"project_id"`
	Properties map[string]string `json:"properties"`
}

// Response defines errors and status of the call to the API
type Response struct {
	Errors  []string `json:"errors"`
	Success bool `json:"success"`
}

// Client implements single project binding to
// exponea API
type API struct {
	// id of the project
	projectID     string

	projectSecret string

	// target API URL
	target        string

	httpClient    *http.Client
}

// NewAPI creates new Client configuration based on the
// given project ID
func NewAPI(projectID, projectSecret string) *API {
	return NewAPIWithTarget(projectID, projectSecret, DefaultAPIEndpoint)
}

// NewAPIWithTarget creates new API client and let's user configure
// target endpoint. This can be used in case of testing or using more
// Exponea api targets
func NewAPIWithTarget(projectID, projectSecret, target string) *API {
	return &API{
		projectID: projectID,
		projectSecret: projectSecret,
		target: target,
		httpClient: &http.Client{},
	}
}

func (c *API) sendRequest(url string, model, response interface{}) error {
	requestData, err := json.Marshal(model)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post(c.target + url, "application/json", bytes.NewReader(requestData))
	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(responseData, response); err != nil {
		return err
	}

	return nil
}

// SendRequest sends request to the given endpoint, marshalling
// the model to JSON and awaiting response of type Response
func (c *API) SendRequest(url string, model interface{}) (*Response, error) {
	apiResponse := &Response{}

	return apiResponse, c.sendRequest(url, model, apiResponse)
}

// SendBulkRequest sends bulk request to the API
func (c *API) SendBulkRequest(bulk *Bulk) (*BulkResponse, error) {
	bulkResponse := &BulkResponse{}

	return bulkResponse, c.sendRequest(BulkEndpoint, bulk, bulkResponse)
}

// SendEvent sends given event data to the Events endpoint
// and returns the server response. In case any errors were
// caused by the network or data serialization, an error will
// be returned
func (c *API) Track(event *Event) (*Response, error) {
	if event.ProjectID == "" {
		event.ProjectID = c.projectID
	}

	return c.SendRequest(EventsEndpoint, event)
}

// SendCustomer sends given event data to the Customers endpoint
// @see SendEvent
func (c *API) Update(customer *Customer) (*Response, error) {
	if customer.ProjectID == "" {
		customer.ProjectID = c.projectID
	}

	return c.SendRequest(CustomersEndpoint, customer)
}

// Bulk represents bulk request to the API, which can
// be used to send multiple requests at once
func (c *API) Bulk(models ...interface{}) (*BulkResponse, error) {
	bulk := &Bulk{
		Commands: make([]*BulkCommand, 5),
	}

	for _, model := range models {
		command := &BulkCommand{}

		switch model.(type) {
		case Event:
			command.Name = EventsEndpoint
		case Customer:
			command.Name = CustomersEndpoint
		}

		// assign model
		command.Data = model
		// escape the slash for the json url
		command.Name = strings.Replace(command.Name, "/", "\\/", -1)

		// append to the bulk
		bulk.Commands = append(bulk.Commands, command)
	}

	return c.SendBulkRequest(bulk)
}
