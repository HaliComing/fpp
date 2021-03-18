package request

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	sj "github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/HaliComing/fpp/pkg/serializer"
	"github.com/HaliComing/fpp/pkg/util"
)

// GeneralClient 通用 HTTP Client
var GeneralClient Client = HTTPClient{}

// Response 请求的响应或错误信息
type Response struct {
	Err      error
	Response *http.Response
}

// Client 请求客户端
type Client interface {
	Request(method, target string, body io.Reader, opts ...Option) *Response
}

// HTTPClient 实现 Client 接口
type HTTPClient struct {
}

// Option 发送请求的额外设置
type Option interface {
	apply(*options)
}

type options struct {
	timeout       time.Duration
	header        http.Header
	signTTL       int64
	ctx           context.Context
	contentLength int64
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func newDefaultOption() *options {
	return &options{
		header:        http.Header{},
		timeout:       time.Duration(30) * time.Second,
		contentLength: -1,
	}
}

// WithTimeout 设置请求超时
func WithTimeout(t time.Duration) Option {
	return optionFunc(func(o *options) {
		o.timeout = t
	})
}

// WithContext 设置请求上下文
func WithContext(c context.Context) Option {
	return optionFunc(func(o *options) {
		o.ctx = c
	})
}

// WithHeader 设置请求Header
func WithHeader(header http.Header) Option {
	return optionFunc(func(o *options) {
		for k, v := range header {
			o.header[k] = v
		}
	})
}

// WithoutHeader 设置清除请求Header
func WithoutHeader(header []string) Option {
	return optionFunc(func(o *options) {
		for _, v := range header {
			delete(o.header, v)
		}

	})
}

// WithContentLength 设置请求大小
func WithContentLength(s int64) Option {
	return optionFunc(func(o *options) {
		o.contentLength = s
	})
}

// Request 发送HTTP请求
func (c HTTPClient) Request(method, target string, body io.Reader, opts ...Option) *Response {
	// 应用额外设置
	options := newDefaultOption()
	for _, o := range opts {
		o.apply(options)
	}

	// 创建请求客户端
	client := &http.Client{Timeout: options.timeout}

	// size为0时将body设为nil
	if options.contentLength == 0 {
		body = nil
	}

	// 创建请求
	var (
		req *http.Request
		err error
	)
	if options.ctx != nil {
		req, err = http.NewRequestWithContext(options.ctx, method, target, body)
	} else {
		req, err = http.NewRequest(method, target, body)
	}
	if err != nil {
		return &Response{Err: err}
	}

	// 添加请求相关设置
	req.Header = options.header
	if options.contentLength != -1 {
		req.ContentLength = options.contentLength
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return &Response{Err: err}
	}

	return &Response{Err: nil, Response: resp}
}

// GetResponse 检查响应并获取响应正文
func (resp *Response) GetResponse() (string, error) {
	if resp.Err != nil {
		return "", resp.Err
	}
	respBody, err := ioutil.ReadAll(resp.Response.Body)
	_ = resp.Response.Body.Close()

	return string(respBody), err
}

// CheckHTTPResponse 检查请求响应HTTP状态码
func (resp *Response) CheckHTTPResponse(status int) *Response {
	if resp.Err != nil {
		return resp
	}

	// 检查HTTP状态码
	if resp.Response.StatusCode != status {
		resp.Err = fmt.Errorf("服务器返回非正常HTTP状态%d", resp.Response.StatusCode)
	}
	return resp
}

// DecodeResponse 尝试解析为serializer.Response，并对状态码进行检查
func (resp *Response) DecodeResponse() (*serializer.Response, error) {
	if resp.Err != nil {
		return nil, resp.Err
	}

	respString, err := resp.GetResponse()
	if err != nil {
		return nil, err
	}

	var res serializer.Response
	err = json.Unmarshal([]byte(respString), &res)
	if err != nil {
		util.Log().Debug("[Conf] Unable to parse callback server response, Error = %s", string(respString))
		return nil, err
	}
	return &res, nil
}

// NopRSCloser 实现不完整seeker
type NopRSCloser struct {
	body   io.ReadCloser
	status *rscStatus
}

type rscStatus struct {
	// http.ServeContent 会读取一小块以决定内容类型，
	// 但是响应body无法实现seek，所以此项为真时第一个read会返回假数据
	IgnoreFirst bool

	Size int64
}

// GetRSCloser 返回带有空seeker的RSCloser，供http.ServeContent使用
func (resp *Response) GetRSCloser() (*NopRSCloser, error) {
	if resp.Err != nil {
		return nil, resp.Err
	}

	return &NopRSCloser{
		body: resp.Response.Body,
		status: &rscStatus{
			Size: resp.Response.ContentLength,
		},
	}, resp.Err
}

// SetFirstFakeChunk 开启第一次read返回空数据
func (instance NopRSCloser) SetFirstFakeChunk() {
	instance.status.IgnoreFirst = true
}

// SetContentLength 设置数据流大小
func (instance NopRSCloser) SetContentLength(size int64) {
	instance.status.Size = size
}

// Read 实现 NopRSCloser reader
func (instance NopRSCloser) Read(p []byte) (n int, err error) {
	if instance.status.IgnoreFirst && len(p) == 512 {
		return 0, io.EOF
	}
	return instance.body.Read(p)
}

// Close 实现 NopRSCloser closer
func (instance NopRSCloser) Close() error {
	return instance.body.Close()
}

// Seek 实现 NopRSCloser seeker, 只实现seek开头/结尾以便http.ServeContent用于确定正文大小
func (instance NopRSCloser) Seek(offset int64, whence int) (int64, error) {
	// 进行第一次Seek操作后，取消忽略选项
	if instance.status.IgnoreFirst {
		instance.status.IgnoreFirst = false
	}
	if offset == 0 {
		switch whence {
		case io.SeekStart:
			return 0, nil
		case io.SeekEnd:
			return instance.status.Size, nil
		}
	}
	return 0, errors.New("未实现")

}

func RequestProxy(testUrl, proxy string) (bool, string, error) {
	proxyURL, _ := url.Parse(proxy)

	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	netTransport := &http.Transport{
		Proxy:               http.ProxyURL(proxyURL),
		TLSClientConfig:     tlsConfig,
		MaxIdleConnsPerHost: 50,
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 20,
		Transport: netTransport,
	}
	request, _ := http.NewRequest("GET", testUrl, nil)
	//设置一个header
	request.Header.Add("accept", "text/plain")
	resp, err := httpClient.Do(request)

	if err != nil {
		return false, "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		json, err := sj.NewFromReader(resp.Body)
		if err != nil {
			return false, "", err
		}
		origin, err := json.Get("origin").String()
		if err != nil {
			return false, "", err
		}
		return true, origin, nil
	}
	return false, "", errors.New("error")
}

func RequestIP(testUrl string) (string, error) {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	netTransport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		MaxIdleConnsPerHost: 50,
	}
	httpClient := &http.Client{
		Timeout:   time.Second * 20,
		Transport: netTransport,
	}
	request, _ := http.NewRequest("GET", testUrl, nil)
	//设置一个header
	request.Header.Add("accept", "text/plain")
	resp, err := httpClient.Do(request)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		json, err := sj.NewFromReader(resp.Body)
		if err != nil {
			return "", err
		}
		origin, err := json.Get("origin").String()
		if err != nil {
			return "", err
		}
		return origin, nil
	}
	return "", errors.New("error")
}
