package common

import (
	"context"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
	"sync"
)

type ChromePool struct {
	pool       *sync.Pool
	deferFuncs []func()
}

func NewChromePool() *ChromePool {
	var cleanFunc []func()
	pool := &ChromePool{
		pool: &sync.Pool{
			New: func() any {
				// 创建一个自定义的Chrome选项
				opts := append(chromedp.DefaultExecAllocatorOptions[:],
					chromedp.Flag("headless", true), // 取消headless模式
				)

				// 创建一个自定义的Chrome执行器
				allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)

				// 使用自定义的执行器创建新的上下文
				ctx, chdCancel := chromedp.NewContext(allocCtx)

				cleanFunc = append(cleanFunc, func() {
					chdCancel()
					cancelAlloc()
				})
				return ctx
			},
		},
	}
	pool.deferFuncs = cleanFunc
	return pool
}

func (c *ChromePool) Close() {
	zap.L().Info("Closing Chrome processes")
	for _, f := range c.deferFuncs {
		f()
	}
}

func (c *ChromePool) GetInstance() context.Context {
	return c.pool.Get().(context.Context)
}
func (c *ChromePool) PutInstance(instance context.Context) {
	c.pool.Put(instance)
}

func OpenChrome(cnt context.Context) (ctx context.Context, cleanFunc func()) {
	var customOpts = []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
	}

	// 创建一个自定义的Chrome选项
	opts := append(chromedp.DefaultExecAllocatorOptions[:])
	customOpts = append(customOpts, chromedp.Flag("proxy-server", "http://localhot:10809")) //todo
	//set http proxy
	//if proxy := GetConfig().Http.Proxy; proxy != "" {
	//	customOpts = append(customOpts, chromedp.Flag("proxy-server", proxy))
	//}
	customOpts = append(chromedp.DefaultExecAllocatorOptions[:], customOpts...)

	// 创建一个自定义的Chrome执行器
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(cnt, opts...)

	// 使用自定义的执行器创建新的上下文
	ctx, chdCancel := chromedp.NewContext(allocCtx)

	cleanFunc = func() {
		chdCancel()
		cancelAlloc()
	}
	return
}
