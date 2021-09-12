package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	httpclient "github.com/ddliu/go-httpclient"
)

const (
	USERAGENT       = "my awsome httpclient"
	TIMEOUT         = 30 // ms
	CONNECT_TIMEOUT = 5  // ms
)

type EnforcerInput struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	UserId int    `json:"user_id"`
}

type EnforcerService struct {
	Endpoint string
}

var enforcerInstance *EnforcerService
var enforcerInstanceOnce sync.Once

func GetEnforcerService() *EnforcerService {
	enforcerInstanceOnce.Do(func() {
		httpclient.Defaults(httpclient.Map{
			httpclient.OPT_USERAGENT: "go-app",
			"Accept-Language":        "en-us",
		})

		enforcerInstance = &EnforcerService{os.Getenv("ENFORCER_ENPOINT")}
	})
	return enforcerInstance
}

func (e *EnforcerService) Enforcer(user_id int, r *http.Request) bool {
	type Data struct {
		Input EnforcerInput `json:"input"`
	}
	type Result struct {
		Allow bool `json:"result"`
	}
	data := Data{
		EnforcerInput{
			Method: r.Method,
			Path:   r.URL.Path,
			UserId: user_id,
		},
	}

	res, err := httpclient.PostJson(e.Endpoint, data)
	if err != nil || res.StatusCode != 200 {
		return false
	}
	body, err := res.ReadAll()
	if err != nil {
		return false
	}

	var result Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false
	}

	return result.Allow
}
