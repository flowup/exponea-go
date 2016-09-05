package exponea

import (
  "net/http"
  "encoding/json"
  "bytes"
  "io/ioutil"
)

const (
  APIEndpoint = "https://api.exponea.com"
  EventsEndpoint = APIEndpoint + "/crm/events"
  CustomersEndpoint = APIEndpoint + "/crm/customers"
)

// Event is an encapsulation for event sending request
type Event struct {
  Customers  map[string]string `json:"customer_ids"`
  ProjectID  string `json:"project_id"`
  Type       string `json:"type"`
  Properties map[string]string `json:"properties"`
  Timestamp  int `json:"timestamp"`
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
type Client struct {
  // id of the project
  project string
}

// NewClient creates new Client configuration based on the
// given project ID
func NewClient(projectID string) *Client {
  return &Client{project: projectID}
}

// SendRequest sends request to the given endpoint, marshalling
// the model to JSON and awaiting response of type Response
func (c *Client) SendRequest(url string, model interface{}) (*Response, error) {
  requestData, err := json.Marshal(model)
  if err != nil {
    return nil, err
  }

  resp, err := http.Post(url, "application/json", bytes.NewReader(requestData))
  if err != nil {
    return nil, err
  }

  responseData, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  apiResponse := &Response{}
  if err = json.Unmarshal(responseData, apiResponse); err != nil {
    return nil, err
  }

  return apiResponse, nil
}

// SendEvent sends given event data to the Events endpoint
// and returns the server response. In case any errors were
// caused by the network or data serialization, an error will
// be returned
func (c *Client) SendEvent(event *Event) (*Response, error) {
  if event.ProjectID == "" {
    event.ProjectID = c.project
  }

  return c.SendRequest(EventsEndpoint, event)
}

// SendCustomer sends given event data to the Customers endpoint
// @see SendEvent
func (c *Client) SendCustomer(customer *Customer) (*Response, error) {
  if customer.ProjectID == "" {
    customer.ProjectID = c.project
  }

  return c.SendRequest(CustomersEndpoint, customer)
}