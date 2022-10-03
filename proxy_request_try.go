package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func proxy() {
	// 请求反向代理的端口
	req, err := http.NewRequest("GET", "http://1.14.120.233:9090", nil)
	if err != nil {
		log.Fatal(err)
	}


	//req.Header.Add("Authorization","Bearer kubeconfig-user-mtxnk.c-gfv2c:h86t2zzpjcq8lksd82l24l6ld7pkdwsh4264thvbfxldntkmdmf2c8")
	rsp,err := http.DefaultClient.Do(req)
	if err != nil{
		log.Fatal(err)
	}
	defer rsp.Body.Close()
	b,_:=ioutil.ReadAll(rsp.Body)
	fmt.Println(string(b))

}
