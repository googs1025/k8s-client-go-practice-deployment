package initClient

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

// 使用反向代理，会有token认证问题
//var K8sClient *kubernetes.Clientset
//
//func init() {
//	config := &rest.Config{
//		Host: "http://1.14.120.233:9090",
//		//BearerToken: "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb2dJQkFBS0NBUUVBcGlWVmdSc0hJZ3lrZHlUMWVRc2pIbzB1U29hRGNDNFg0cE5yMS8xdzUvK3dTK2gwClZzU0hhNXdxUjBRZG5JczBmZklMUEYxU3B3Yi9oNXEvRUx0R291dmRTV2xHeERYYU1nam4yNmRnTnhTNVU4alcKQy9YL3QvWUpYaHFxRm1iZ1FFakRmZHZZMXNZOVBvcnJJOE1razJFMWp1cVB2WW4zYmFwVUpnQkd2VUFoUXdPRQpCOHgvcjhVaVhGaHAxS1dmUmZtZGRlMENTNytxaXhzY1hBOUM5R2tYYVBMbDZaUmhibzBFa3ppYm5QNGlHM1BnCjR5eWtyZGNXbTJnMEMyTW9jbWFzd2ZZekREUjFQQkNscGtIem92cU9xTytnTWZhNWZyQXF5Z3BFWWdBNE9panEKdGNDa1hIZjI5UWcxRmxQclI0TW9WS1M0cDBSVGhVb0lsWDBoYVFJREFRQUJBb0lCQUFpdGhJVEV3NStjcDI1dApxTUNVdTFYYUsrUEttTXpnSzNFekgvdmRDZXVrS0RJZXh3ek5JUUdXMjRKelpWU0sxTWdMUDFqOHl0ZGNmelkyCjkrbkl6a3l1SXhXMWdQTzRtRmZxclNtRTJYcW5BM01EMTJJeWpCT3dyeGFTTC9ZUmszN29EZ1hoMkxhSERpWFoKSGFUMWlWQ0ZVRVhScklaSzBYaVIzK2xJTkdtbi9VUUFVSWVnTHc0TW9pNnNSbXdLajFmSlNwMmdZam5rbkswZwphVjlsQ0NmeGYvWDh5bTdDcHdrVm9Xa1F1dG1VYnpLTGEzcHBWZWRKVVBjczFWUWFQZjVKUytJaUhwQzU1NStNCmZ4Z2hJdERsZUw1MGtXeGo1czYvK1o2RWpNSnk2OXA4WGNjamwxR1luU0JtaDJtOXVxbCtjWWU3cTBpMzEyQysKVnJOYjJLMENnWUVBeFRBNFp2M0xtMWgrRGFGZUxPTUoyTVZRZzZ0UkFHZlJHTEJoWVgyVGRTRWpvbm5XakhSKwpnbXhGR002RkI2VmN5NkVkZmtBYWZUMWcrS3owd1ZhM2FERlB1Z0VRY1c5M1lwTENPSEovOE5zSzV2QXpwSjZCCmljSjNYcGhtN1pBZ0xUcFppSUFMVk5ZQ3JsWFp1TzZhb2hiMExyZ2lmTUhRbkw0ZEsrSEhrNWNDZ1lFQTE3THkKUTBLanV0RG9GU2Y4QnArM2ZqczFyQ2gzK0hKRUZibTh3a21mOWRhOHdFaTV2MGNWbGxXZXRvcTY2ZUF6UkNnRwpOZTBnQ29aSWVpN3U5ZU9ZWDZuRmZjQ3NSRWRldlNUL1FLZ0pTY0NwQWRqTGpEZ3hIVW5qRUJDUnF4c3hCQmFhCm8yTW00VExabUd0ME01Y2x4ZjZYa3VuNWNmQnZ2amFPV2t0RGt2OENnWUJxbXVvelRBeVNuS0h4Ym9kQ2p6QVkKb0h3cDR5bTBwV3ZYQkN4eGozbHovb21NWW9CS2lRU0lNRTZlM2EvdjlVZVkweTdsdlhSVXR5VkE3QWlhcWU0WQpCMmpKNzU5YkpGOFB6TFh4M0gwczBzOHZFVGRxVFVOTkhmUjVFTDI4dTRtWnlnenpqZjRTVEcxQW9TdEhIc0E4CjExb0dGQWlaR0JOWFdqVGRMNEE5V3dLQmdIVEc2OFFnV3ZZMFRjSE9jUExCRzUyYXZyY2kvYmlqWEZzS3dMZkwKRm5BSlB3MDNFbUVOUWhHdTd3dFMxbGp5U2E4WG9DMG40Tlh4MTJGVzhZWnNIcjJEODJqZW5DVW5JcEp5YWtMOQo5bkZZZmVlREVNZ3NUK0xVY3JycXpZSithUzRXY3NnTVVTdFExVjlncFh1YzFCVjZmV05MaXdIMXN2bWZIYmlpCjBNWFBBb0dBVmpGYVZZS0lEeDRKUlNmRWVLVTMzREU3ektwODl1VXp5MjRXaXdicVBQTWdrc29hSkVOTWs1Qk8KVXJiRHJIU3lUZ29xZUVGcVJuQVZoODFlUTgrOWUrNnIwN1Izb1JFdTFJcXFVVzlyYWZ6TnVYRnJMOEZXZW9XQgpKMHZGTVhRMHo4UjQvQ1dlU0swdzI5WVVMZ3JOVWo0L0dVdXpWTzh1ZlZIdUMydGF4WmM9Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==",
//
//	}
//	clientSet, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	K8sClient = clientSet
//
//}


var K8sClient *kubernetes.Clientset

func init() {

	// 两个选一个用
	//config := kubeConfig()
	config := kubeProxyConfig()


	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
	K8sClient = clientSet
}

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	return os.Getenv("USERPROFILE")

}

func kubeConfig() *rest.Config {
	// 法一：直接在k8s上运行的代码
	var kubeConfig *string

	if home := HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		log.Panic(err.Error())
	}
	return config
}

func kubeProxyConfig() *rest.Config {
	// 法二：需要用端口转换 kubectl proxy --address="0.0.0.0" --accept-hosts='^*$' --port=8009
	config := &rest.Config{
		Host: "http://1.14.120.233:8009",
	}
	return config
}