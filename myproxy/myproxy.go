package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func Myproxy() {
	var tlsConfig = &tls.Config{
		InsecureSkipVerify: true,  // 忽略证书验证
	}

	var transport http.RoundTripper = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 30 * time.Second,
			KeepAlive: 30 * time.Second,

		}).DialContext,
		MaxIdleConns: 100,
		IdleConnTimeout: 90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: tlsConfig,
		DisableCompression: true,
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		server, _ := url.Parse("https://10.0.0.12:6443")
		log.Println(request.URL.Path)	// 显示请求路径
		p := httputil.NewSingleHostReverseProxy(server)
		p.Transport = transport
		p.ServeHTTP(writer, request)
	})

	log.Println("开始进行反向代理 k8s端口")
	err := http.ListenAndServe("0.0.0.0:9090", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func main()  {
	Myproxy()
}
