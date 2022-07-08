package simple

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/r3labs/sse/v2"
	"github.com/simpleflags/evaluation"
	"github.com/simpleflags/golang-server-sdk/connector"
	"github.com/simpleflags/golang-server-sdk/log"
	"go.uber.org/atomic"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Option func(c *simpleFlagsConfig)

type simpleFlagsConfig struct {
	baseURL      string
	eventsURL    string
	retryWaitMax time.Duration
}

func WithBaseURL(baseURL string) Option {
	return func(c *simpleFlagsConfig) {
		c.baseURL = baseURL
	}
}

func WithEventsURL(eventsURL string) Option {
	return func(c *simpleFlagsConfig) {
		c.eventsURL = eventsURL
	}
}

func WithRetryWaitMax(max time.Duration) Option {
	return func(c *simpleFlagsConfig) {
		c.retryWaitMax = max
	}
}

type HttpConnector struct {
	apiKey          string
	config          simpleFlagsConfig
	baseApiClient   *http.Client
	eventsApiClient *http.Client
	stream          *sse.Client
	cancelStream    context.CancelFunc
	streamConnected *atomic.Bool
}

func NewHttpConnector(apiKey string, options ...Option) *HttpConnector {

	config := simpleFlagsConfig{
		baseURL:      "http://localhost:1324/api",
		eventsURL:    "http://localhost:1324/api",
		retryWaitMax: time.Second * 60,
	}

	for _, option := range options {
		option(&config)
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = time.Second
	retryClient.RetryMax = int(time.Second * config.retryWaitMax)

	baseApiClient := retryClient.StandardClient()
	eventsApiClient := retryClient.StandardClient()

	return &HttpConnector{
		apiKey:          apiKey,
		config:          config,
		baseApiClient:   baseApiClient,
		eventsApiClient: eventsApiClient,
		streamConnected: atomic.NewBool(false),
	}
}

func (f *HttpConnector) Configurations(ctx context.Context, identifiers ...string) (evaluation.Configurations, error) {
	address, err := url.Parse(f.config.baseURL + "/configs")
	if err != nil {
		return evaluation.Configurations{}, err
	}
	q := address.Query()
	if len(identifiers) > 0 {
		strIdentifiers := strings.Join(identifiers, ",")
		q.Set("identifiers", strIdentifiers)
	}
	address.RawQuery = q.Encode()
	response, err := get(ctx, f.apiKey, f.baseApiClient, address.String())
	if err != nil {
		return []evaluation.Configuration{}, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []evaluation.Configuration{}, err
	}

	if response.StatusCode != 200 {
		return []evaluation.Configuration{}, err
	}

	var configurations evaluation.Configurations
	err = json.Unmarshal(bytes, &configurations)
	if err != nil {
		return []evaluation.Configuration{}, err
	}
	return configurations, nil
}

func (f *HttpConnector) Variables(ctx context.Context, identifiers ...string) ([]evaluation.Variable, error) {
	address, err := url.Parse(f.config.baseURL + "/vars")
	if err != nil {
		return []evaluation.Variable{}, err
	}
	q := address.Query()
	if len(identifiers) > 0 {
		strIdentifiers := strings.Join(identifiers, ",")
		q.Set("identifiers", strIdentifiers)
	}
	address.RawQuery = q.Encode()
	response, err := get(ctx, f.apiKey, f.baseApiClient, address.String())
	if err != nil {
		return []evaluation.Variable{}, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []evaluation.Variable{}, err
	}

	if response.StatusCode != 200 {
		return []evaluation.Variable{}, err
	}

	var variables []evaluation.Variable
	err = json.Unmarshal(bytes, &variables)
	if err != nil {
		return []evaluation.Variable{}, err
	}
	return variables, nil
}

func (f *HttpConnector) Stream(ctx context.Context, updater connector.Updater) error {
	if f.streamConnected.Load() {
		log.Info("stream already started")
		return nil
	}
	sseCtx, cancel := context.WithCancel(ctx)
	f.cancelStream = cancel
	f.stream = sse.NewClient(f.config.baseURL + "/stream")
	f.stream.Headers["API-Key"] = f.apiKey
	f.stream.Connection.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	f.stream.OnDisconnect(func(c *sse.Client) {
		updater.OnDisconnect()
	})

	var connErr chan error
	go func(connected *atomic.Bool) {
		connErr <- f.stream.SubscribeWithContext(sseCtx, "", func(msg *sse.Event) {
			// Got some data!
			updater.OnEvent(&connector.Msg{
				Event: msg.Event,
				Data:  msg.Data,
			})
		})
	}(f.streamConnected)

	select {
	case err := <-connErr:
		return fmt.Errorf("error subscribing to SSE %v", err)
	case <-ctx.Done():
		return nil
	default:
		updater.OnConnect()
		return nil
	}
}

func (f *HttpConnector) Close() error {
	if f.stream != nil {
		f.cancelStream()
	}
	return nil
}

func get(ctx context.Context, apiKey string, client *http.Client, requestUrl string) (*http.Response, error) {
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("API-Key", apiKey)
	return client.Do(req.WithContext(ctx))
}
