package client

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"log"
	"os"
	"testing"
	"time"
)

var testUrl = "https://www.nosadfun.com/book/14683/292251.html"

func TestFirstCase(t *testing.T) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run
	var b1 []byte
	if err := chromedp.Run(ctx,
		// emulate iPhone 7 landscape
		chromedp.Emulate(device.IPadPro),
		chromedp.Navigate(testUrl),
		chromedp.CaptureScreenshot(&b1),

		// reset
		chromedp.Emulate(device.Reset),
	); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("screenshot1.png", b1, 0o644); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote screenshot1.png and screenshot2.png")
}

/*
*

	    ctx, cancel := chromedp.NewContext(
	        context.Background(),
	    )
	    defer cancel()

	    ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	    defer cancel()

	    testUrl := "http://webcode.me/click.html"

	    var ua string

	    err := chromedp.Run(ctx,

	        chromedp.Emulate(device.IPhone11),
	        chromedp.Navigate(testUrl),
	        chromedp.Click("button", chromedp.NodeVisible),
	        chromedp.Text("#output", &ua),
	    )

	    if err != nil {
	        log.Fatal(err)
	    }

	    log.Printf("User agent: %s\n", ua)
	}

In the example, we click on a button of a web page. The web page shows the client's user agent in the output div.

err := chromedp.Run(ctx,

	chromedp.Emulate(device.IPhone11),
	chromedp.Navigate(testUrl),
	chromedp.Click("button", chromedp.NodeVisible),
	chromedp.Text("#output", &ua),

)
In the task list, we navigate to the URL, click on the button, and retrieve the text output. We get our user agent. We emulate an IPhone11 device with chromedp.Emulate.

# Advertisements

Create screenshot
We can create a screenshot of an element with chromedp.Screenshot. The chromedp.FullScreenshot takes a screenshot of the entire browser viewport.

screenshot.go
package main

import (

	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chromedp/chromedp"

)

func main() {

	    ctx, cancel := chromedp.NewContext(
	        context.Background(),
	    )

	    defer cancel()

	    testUrl := "http://webcode.me"

	    var buf []byte
	    if err := chromedp.Run(ctx, ElementScreenshot(testUrl, "body", &buf)); err != nil {
	        log.Fatal(err)
	    }

	    if err := ioutil.WriteFile("body.png", buf, 0o644); err != nil {
	        log.Fatal(err)
	    }

	    if err := chromedp.Run(ctx, FullScreenshot(testUrl, 90, &buf)); err != nil {
	        log.Fatal(err)
	    }

	    if err := ioutil.WriteFile("full.png", buf, 0o644); err != nil {
	        log.Fatal(err)
	    }

	    fmt.Println("screenshots created")
	}

func ElementScreenshot(testUrl, sel string, res *[]byte) chromedp.Tasks {

	    return chromedp.Tasks{

	        chromedp.Navigate(testUrl),
	        chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	    }
	}

func FullScreenshot(testUrl string, quality int, res *[]byte) chromedp.Tasks {

	    return chromedp.Tasks{

	        chromedp.Navigate(testUrl),
	        chromedp.FullScreenshot(res, quality),
	    }
	}
*/
func TestReadContent(t *testing.T) {
	// 创建一个实例
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 创建一个自定义的Chrome选项
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // 取消headless模式
	)

	// 创建一个自定义的Chrome执行器
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	// 使用自定义的执行器创建新的上下文
	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// run
	var content string
	if err := chromedp.Run(ctx,
		// emulate iPhone 7 landscape
		chromedp.Emulate(device.IPadPro),
		chromedp.Navigate(testUrl),
		//chromedp.InnerHTML("div[class=RBGsectionThree-content]", &content, chromedp.ByQuery),
		chromedp.InnerHTML(".RBGsectionThree-content", &content, chromedp.ByQuery),
		chromedp.Tasks{},
	); err != nil {
		log.Fatal(err)
	}
	println(content)
}

func TestMultipleInstances(t *testing.T) {
	// 创建一个实例
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 创建一个自定义的Chrome选项
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // 取消headless模式
	)

	// 创建一个自定义的Chrome执行器
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	// 使用自定义的执行器创建新的上下文
	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate("http://www.baidu.com")); err != nil {
		t.Fatal(err)
	}
	if err := chromedp.Run(ctx, chromedp.Navigate("http://www.bing.com")); err != nil {
		t.Fatal(err)
	}

	select {}
}

func TestDefer(t *testing.T) {
	deferFuncs := childFUnc()
	defer deferFuncs()
	t.Log("TestDefer")
}

func childFUnc() func() {
	// 创建一个实例
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 创建一个自定义的Chrome选项
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // 取消headless模式
	)

	// 创建一个自定义的Chrome执行器
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	// 使用自定义的执行器创建新的上下文
	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()
	println("child process")
	return func() {
		println("defer funcs")
		cancelAlloc()
		cancel()
	}
}

func TestOpen(t *testing.T) {
	ctx, cleanFunc := OpenChrome(context.Background())
	defer cleanFunc()

	chromedp.Run(ctx,
		chromedp.Navigate("http://www.google.com"),
		chromedp.WaitVisible(".lnXdpd", chromedp.ByQuery),
	)
	chromedp.Stop()
	cleanFunc()
	time.Sleep(5 * time.Second)
}

func TestNotReady(t *testing.T) {
	ctx, cleanFunc := OpenChrome(context.Background())
	defer cleanFunc()

	var content string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.nosadfun.com/book/2651/2449946.html"),
		chromedp.WaitNotPresent("//p[contains(text(),'内容未加载完成')]", chromedp.BySearch),
		chromedp.InnerHTML("//div[@class='RBGsectionThree-content']", &content, chromedp.BySearch),
		//chromedp.InnerHTML(".RBGsectionThree-content", &content, chromedp.ByQuery),
	)
	println(content)
	if err != nil {
		log.Fatal(err)
	}
	cleanFunc()
	time.Sleep(5 * time.Second)

}
