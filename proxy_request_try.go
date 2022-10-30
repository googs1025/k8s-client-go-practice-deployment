package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func proxy() {
	// 请求反向代理的端口
	req, err := http.NewRequest("GET", "http://1.14.120.233:6443", nil)
	if err != nil {
		log.Fatal(err)
	}


	req.Header.Add("Authorization","Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImxwSDVFOWRrbGtfcHVCN1k5T01VWjNIQ1U3NEZFWjNxSEQ5SmJpTGxTREEifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2hY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJuYW1lc3BhY2UtY29udHJvbGxlci10b2tlbi1qNjJtNyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJuYW1lc3BhY2UtY29udHJvbGxlciIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjVlYjkzYmZhLTYxMzQtNGE1NC1hOWVmLWI3MWYwYjJiOTAwNCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTpuYW1lc3BhY2UtY29udHJvbGxlciJ9.NpMVb_FjZ3417cmiPUasyVHsUPKyR6DmIwSPtSCqO9FsfupSR8C13P9aDnXa_6FwI1P6YgBef-RjX5cr2AbkEcZhgv5PovEo2U3-VHI1anvVK7-Wus01j98FXlB8UiBdy7gpBHcIgvVHUMjVLUEq1fK5mxtvr2XvWIow5QFjFTDHept8PYF2zctFPPy0hQ2CYFfMoUCvfyI-tESVtKnLKwo2Tgc4JBsQAExsc4nqyP8bZSd3F8TpeRh-27c3wFUzz6av1pyR-fhB2MldH0VmTiyGS1yS_JgTzzH8igCTMwHVLAClo9r4KTty24JFWs9hdGJMU6q3f1h0r292rCRldw")
	rsp,err := http.DefaultClient.Do(req)
	if err != nil{
		log.Fatal(err)
	}
	defer rsp.Body.Close()
	b,_:=ioutil.ReadAll(rsp.Body)
	fmt.Println(string(b))

}

func main() {
	proxy()
}
