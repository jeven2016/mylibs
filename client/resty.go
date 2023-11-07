package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jeven2016/mylibs/utils"
	"net/http"
	"sync"
	"time"
)

const DefaultRetries = 3

var restyInstanceMap = make(map[string]*resty.Client)
var restyLock sync.Mutex

// GetRestyClient 一个域名对应一个resty.Client
// https://github.com/go-resty/resty/issues/612
func GetRestyClient(url string, retry bool) (*resty.Client, error) {
	base := utils.ParseBaseUri(url)
	if base == "" {
		return nil, fmt.Errorf("invalid Url: %s", url)
	}

	if client, ok := restyInstanceMap[base]; !ok {
		restyLock.Lock()
		defer restyLock.Unlock()

		if client, ok = restyInstanceMap[base]; ok {
			return client, nil
		}

		newClient := resty.New()
		// Allow GET request with Payload. This is disabled by default.
		newClient.SetAllowGetMethodPayload(true)
		newClient.SetDebug(false)
		//newClient.SetBasicAuth("myuser", "mypass")

		// 设置链接超时
		//newClient.SetTimeout(1 * time.Minute)

		// 设置自定义证书 Refer: http://golang.org/pkg/crypto/tls/#example_Dial
		//newClient.SetTLSClientConfig(&tls.Config{RootCAs: roots})

		// 禁用安全检查
		newClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		if retry {
			newClient.
				SetRetryCount(DefaultRetries).
				SetRetryWaitTime(5 * time.Second).
				SetRetryMaxWaitTime(20 * time.Second).
				// SetRetryAfter sets callback to calculate wait time between retries.
				// Default (nil) implies exponential backoff with jitter
				SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
					return 0, errors.New("quota exceeded")
				})
			// 自定义重试策略
			newClient.AddRetryCondition(
				// RetryConditionFunc type is for retry condition function
				// input: non-nil Response OR request execution error
				func(r *resty.Response, err error) bool {
					return r.StatusCode() == http.StatusTooManyRequests
				},
			)
		}

		restyInstanceMap[base] = newClient
		return newClient, nil
	} else {
		return client, nil
	}

}
