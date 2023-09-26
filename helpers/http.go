package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	ContentTypeJson           = "application/json"
	ContentTypeFormUrlEncoded = "application/x-www-form-urlencoded" // GetContentType 默认值
)

const (
	MethodPost = "POST"
	MethodGet  = "GET"
)

const (
	HTTPRequestEncodeJson = "_json"
	HTTPRequestEncodeRaw  = "_raw"
	HTTPRequestEncodeForm = "_form"
)

type HttpResult struct {
	HttpCode int
	Response []byte
	Header   http.Header
}

type HttpClient struct {
	Host      string
	BasicAuth struct {
		Username string
		Password string
	}

	Timeout    time.Duration // 当前客户端超时时间
	HTTPClient *http.Client  // HTTPClient 客户端对象
	clientInit sync.Once     // 对象创建单次锁
}

type HttpRequestOptions struct {
	// 超时时间
	Timeout time.Duration

	ContentType string

	// 追加url后面的参数
	UrlData map[string]string

	//指定请求头
	Headers map[string]string

	RequestBody interface{}

	Encode string
}

func (hro *HttpRequestOptions) getBody() (string, error) {
	if hro.RequestBody == nil {
		return "", nil
	}
	switch hro.Encode {
	case HTTPRequestEncodeJson:
		reqBody, err := json.Marshal(hro.RequestBody)
		return string(reqBody), err
	case HTTPRequestEncodeRaw:
		var err error
		encodeData, ok := hro.RequestBody.(string)
		if !ok {
			err = errors.New("raw data need string type")
		}
		return encodeData, err
	case HTTPRequestEncodeForm:
		fallthrough
	default:
		return hro.getFormRequestData()
	}
}

// getFormRequestData 格式化输出 form 请求data
// 不支持form a[0]=?&a[1]=?的方式，可以使用 []type{} 的方式进行json上的转换，读取key 之后需要做json_decode()
func (hro *HttpRequestOptions) getFormRequestData() (string, error) {

	v := url.Values{}

	if data, ok := hro.RequestBody.(map[string]string); ok {
		for key, value := range data {
			v.Add(key, value)
		}
		return v.Encode(), nil
	}

	if data, ok := hro.RequestBody.(map[string]interface{}); ok {
		for key, value := range data {
			var vStr string
			switch value := value.(type) {
			case string:
				vStr = value
			default:
				if tmp, err := json.Marshal(value); err != nil {
					return "", err
				} else {
					vStr = string(tmp)
				}
			}

			v.Add(key, vStr)
		}
		return v.Encode(), nil
	}

	return "", errors.New("unSupport RequestBody type [" + reflect.TypeOf(hro.RequestBody).Name() + "]")
}

// getData 追加url 后面的querystring &aa=bb&cc=dd
func (hro *HttpRequestOptions) getData() (string, error) {
	if len(hro.UrlData) > 0 {
		v := url.Values{}
		for key, value := range hro.UrlData {
			v.Add(key, value)
		}
		return v.Encode(), nil
	}

	return "", nil
}

// GetContentType  默认值 "application/x-www-form-urlencoded"
func (hro *HttpRequestOptions) GetContentType() (cType string) {

	if cType = hro.ContentType; cType != "" {
		return cType
	}
	return ContentTypeJson
}

func (client *HttpClient) makeRequest(method string, urlPath string, data io.Reader, opts HttpRequestOptions) (*http.Request, error) {
	req, err := http.NewRequest(method, urlPath, data)
	if err != nil {
		log.Printf("HTTP_ERROR Method[%s] urlPath[%s] ERROR[%+v]", method, urlPath, err)
	}
	if opts.Headers != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	if client.BasicAuth.Username != "" {
		req.SetBasicAuth(client.BasicAuth.Username, client.BasicAuth.Password)
	}

	req.Header.Set("Content-Type", opts.GetContentType())

	return req, nil
}

func (client *HttpClient) httpDo(req *http.Request, opts *HttpRequestOptions) (res HttpResult, err error) {
	client.clientInit.Do(func() {
		if client.HTTPClient == nil {
			timeout := time.Second * 3
			if client.Timeout > 0 {
				timeout = client.Timeout
			}

			client.HTTPClient = &http.Client{
				Timeout: timeout,
			}
		}
	})

	resp, doErr := client.HTTPClient.Do(req)
	if doErr != nil {
		return res, doErr
	}
	if resp != nil {
		res.HttpCode = resp.StatusCode
		res.Response, err = ioutil.ReadAll(resp.Body)
		res.Header = resp.Header
		_ = resp.Body.Close()
	}

	return res, err
}

func (client *HttpClient) HttpGet(urlPath string, opts HttpRequestOptions) (res *HttpResult, err error) {
	urlData, err := opts.getData()
	if err != nil {
		return nil, err
	}

	var reqPath string
	if urlData == "" {
		reqPath = fmt.Sprintf("%s%s", client.Host, urlPath)
	} else {
		reqPath = fmt.Sprintf("%s%s?%s", client.Host, urlPath, urlData)
	}

	req, err := client.makeRequest(MethodGet, reqPath, nil, opts)
	if err != nil {
		log.Printf("HttpGet.makeRequest.failed[%+v]", err)
	}
	log.Printf("http.HttpCient.Get._reqPath[%s], _req[%+v] ,err [%+v]", reqPath, req, err)
	resp, err := client.httpDo(req, &opts)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

func (client *HttpClient) HttpPost(urlPath string, opts HttpRequestOptions) (res *HttpResult, err error) {
	urlData, err := opts.getData()
	if err != nil {
		return nil, err
	}

	requestBody, err := opts.getBody()
	if err != nil {
		return nil, err
	}
	var reqPath string
	if urlData == "" {
		reqPath = fmt.Sprintf("%s%s", client.Host, urlPath)
	} else {
		reqPath = fmt.Sprintf("%s%s?%s", client.Host, urlPath, urlData)
	}
	log.Printf("http.HttpClient.HttpPost._reqPath[%s] ", reqPath)

	req, err := client.makeRequest(MethodPost, reqPath, strings.NewReader(requestBody), opts)
	if err != nil {
		log.Printf("HttpClient.HttpPost.makeRequest.Error._reqPath[%s], _req[%+v] err[%+v]", reqPath, req, err)
		return nil, err
	}
	resp, err := client.httpDo(req, &opts)
	if err != nil {
		return nil, err
	}
	return &resp, err

}
