package exponea

import (
  "github.com/stretchr/testify/suite"
  "testing"
  "github.com/stretchr/testify/assert"
  "net/http/httptest"
  "net/http"
  "encoding/json"
)

const (
  mockClientProjectID = "id"
  mockClientProjectSecret = "secret"
)

type APISuite struct {
  suite.Suite

  server *httptest.Server
  api *API
}

type MockBackend struct {}

func (backend *MockBackend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  successData, _ := json.Marshal(&Response{})
  w.Write(successData)
}

func (s *APISuite) SetupSuite() {
  // create mock server
  s.server = httptest.NewServer(&MockBackend{})

  // create new API
  s.api = NewAPI(mockClientProjectID, mockClientProjectSecret)
  assert.NotEqual(s.T(), nil, s.api)
  assert.Equal(s.T(), mockClientProjectID, s.api.projectID)
  assert.Equal(s.T(), mockClientProjectSecret, s.api.projectSecret)
}

func (s *APISuite) TearDownSuite() {
  // close the server listener on cleanup
  s.server.Close()
}

func (s *APISuite) TestTrack() {
  resp, err := s.api.Track(&Event{
    Customers: map[string]string{
      "registered": "peter.malina@flowup.eu",
    },
    Type: "registration",
    Properties: map[string]string{
      "property": "and it's valueee",
    },
  })

  assert.Equal(s.T(), nil, err)
  assert.NotEqual(s.T(), nil, resp)
}

func (s *APISuite) TestUpdate() {
  resp, err := s.api.Update(&Customer{
    Customers: map[string]string{
      "registered": "peter.malina@flowup.eu",
    },
    Properties: map[string]string{
      "property": "and it's valueee",
    },
  })

  assert.Equal(s.T(), nil, err)
  assert.NotEqual(s.T(), nil, resp)
}

func TestAPISuite(t *testing.T) {
  suite.Run(t, &APISuite{})
}