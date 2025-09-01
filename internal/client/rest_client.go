package client

import (
	"api-ptf-core-business-orchestrator-go-ms/internal/pkg/logger"
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type RestClient struct {
	HttpClient *http.Client
}

func NewRestClient(timeout time.Duration) *RestClient {
	return &RestClient{
		HttpClient: &http.Client{Timeout: timeout},
	}
}

func (rc *RestClient) doRequest(method string, reqData *RequestData) ([]byte, int, error) {
	url, err := reqData.BuildURL()
	if err != nil {
		LogError(err)
		return nil, 0, err
	}

	// Crear request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqData.Body))
	if err != nil {
		LogError(err)
		return nil, 0, err
	}

	// Headers
	for k, v := range reqData.Headers {
		req.Header.Set(k, v)
	}

	// Logging entrada
	LogRequest(method, url, reqData.Headers, reqData.Body)

	start := time.Now()

	// Hacer la llamada
	resp, err := rc.HttpClient.Do(req)
	if err != nil {
		LogError(err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		LogError(err)
		return nil, resp.StatusCode, err
	}

	// Logging salida
	LogResponse(resp.StatusCode, respBody, start)

	return respBody, resp.StatusCode, nil
}

// Métodos expuestos
func (rc *RestClient) Get(req *RequestData) ([]byte, int, error) {
	return rc.doRequest(http.MethodGet, req)
}

func (rc *RestClient) Post(req *RequestData) ([]byte, int, error) {
	return rc.doRequest(http.MethodPost, req)
}

func (rc *RestClient) Put(req *RequestData) ([]byte, int, error) {
	return rc.doRequest(http.MethodPut, req)
}

func (rc *RestClient) Patch(req *RequestData) ([]byte, int, error) {
	return rc.doRequest(http.MethodPatch, req)
}

func (rc *RestClient) Delete(req *RequestData) ([]byte, int, error) {
	return rc.doRequest(http.MethodDelete, req)
}

func LogRequest(method, url string, headers map[string]string, body []byte) {
	logger.Logger().Info("➡️  REQUEST: [%s] %s", zap.String("method", method), zap.String("url", url))
	if headers != nil {
		logger.Logger().Info("   Headers:", zap.Reflect("headers", headers))
	}
	if len(body) > 0 {
		logger.Logger().Info("   Body: %s", zap.String("body", string(body)))
	}
}

func LogResponse(status int, body []byte, start time.Time) {
	duration := time.Since(start)
	logger.Logger().Info("⬅️  RESPONSE: Status=%d, Time=%s", zap.Int("status", status), zap.Duration("duration", duration))
	if len(body) > 0 {
		logger.Logger().Info("   Body: %s", zap.String("body", string(body)))
	}
}

func LogError(err error) {
	logger.Logger().Error("❌ ERROR: %v", zap.Error(err))
}
