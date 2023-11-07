package common

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const CollyMaxRetries = 3

func NewCollector(httpProxy string, maxRetries int) (c *colly.Collector, err error) {
	// c := colly.NewCollector(colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"))
	// c := colly.NewCollector(colly.CacheDir("./temp"))
	c = colly.NewCollector(
		//设置忽略robots协议
		colly.IgnoreRobotsTxt())
	c.SetRequestTimeout(50 * time.Second)

	httpTransport := &http.Transport{
		DisableKeepAlives: true, // Colly uses HTTP keep-alive to enhance scraping speed
		DialContext: (&net.Dialer{
			Timeout:   90 * time.Second,
			KeepAlive: 90 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   90 * time.Second,
		ExpectContinueTimeout: 90 * time.Second,
		Proxy:                 http.ProxyFromEnvironment, //从环境变量获取http proxy地址
	}

	//set http proxy
	if httpProxy != "" {
		proxyUrl, err := url.Parse(httpProxy)
		if err != nil {
			return
		}
		httpTransport.Proxy = http.ProxyURL(proxyUrl)
	}

	c.WithTransport(httpTransport)

	// 对于匹配的域名(当前配置为任何域名),将请求并发数配置为2
	// 通过测试发现,RandomDelay参数对于同步模式也生效
	if err = c.Limit(&colly.LimitRule{
		// glob模式匹配域名
		DomainGlob: "*",

		// 匹配到的域名的并发请求数
		Parallelism: 5,
		// 在发起一个新请求时的随机等待时间
		RandomDelay: time.Duration(500) * time.Millisecond,
	}); err != nil {
		return
	}

	// 是否允许重复请求相同url
	c.AllowURLRevisit = true
	c.Async = false
	c.DetectCharset = true

	// Rotate two socks5 proxies
	//"github.com/gocolly/colly/proxy"
	//rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1337", "socks5://127.0.0.1:1338")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Set error api
	c.OnError(func(r *colly.Response, err error) {
		// r.Request.Retry()
		zap.S().Warn("[Request URL]:", r.StatusCode, " ", r.Request.URL, "failed with response:", r, "\nError:", err)

		if r.StatusCode == http.StatusNotFound {
			r.Ctx.Put("inValidPage", true)
			zap.S().Warn("no retries for 404 page", zap.String("testUrl", r.Request.URL.String()))
			return
		}

		retries := 0
		if lastRetries := r.Ctx.GetAny("retries"); lastRetries != nil {
			retries = lastRetries.(int) + 1
		}
		max := CollyMaxRetries
		if maxRetries > 0 {
			max = maxRetries
		}
		if retries >= max {
			zap.L().Warn("retry aborted after multiple retires", zap.String("retries", strconv.Itoa(retries)))
			return
		}
		r.Ctx.Put("retries", retries)
		err = r.Request.Retry()
		if err != nil {
			zap.S().Warn("error occurs", zap.Error(err))
		}
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("[Response URL]:", r.StatusCode, " ", r.Request.URL)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("[Visiting]", r.URL.String())
	})

	// 随机设置
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	return
}
