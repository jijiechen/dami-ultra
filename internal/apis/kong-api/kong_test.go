package kong_api

import (
	"fmt"
	"testing"
)

const (
	kong_url = "http://ec2-54-166-250-69.compute-1.amazonaws.com:8001/routes"

	wrongMsg = `{"name": "test_route"}`
	goodMsg  = `{
    "protocols": [
        "http",
        "https"
    ],
    "regex_priority": 0,
    "updated_at": 1741244044,
    "hosts": null,
    "name": "test1",
    "request_buffering": true,
    "response_buffering": true,
    "sources": null,
    "created_at": 1741244044,
    "strip_path": true,
    "destinations": null,
    "service": null,
    "path_handling": "v0",
    "methods": null,
    "snis": null,
    "https_redirect_status_code": 426,
    "paths": [
        "/test1"
    ],
    "preserve_host": false,
    "tags": null,
    "headers": null,
    "id": "e3d2b00a-4144-46fe-9477-6a134b504009"
}`
)

func TestApplyKongConfig(t *testing.T) {
	err := ApplyKongConfig(kong_url, goodMsg)
	if nil != err {
		panic(err)
	}
}

/*
{"data":[{"created_at":1741247139,"headers":null,"strip_path":true,"regex_priority":0,"updated_at":1741247139,"hosts":null,"request_buffering":true,"response_buffering":true,"https_redirect_status_code":426,"snis":null,"id":"14a24254-4f03-4d89-a369-2688a49be91f","tags":null,"name":"test","sources":null,"service":{"id":"69192155-e8e1-4b5c-a99f-a8a1890dfe1b"},"destinations":null,"path_handling":"v0","paths":["/test"],"methods":null,"preserve_host":false,"protocols":["http","https"]}],"next":null}
*/
func TestGetKongConfig(t *testing.T) {
	kongConfig, err := GetKongConfig(kong_url)
	if nil != err {
		panic(err)
	}
	fmt.Println(kongConfig)
}

func TestDeleteKongConfig(t *testing.T) {
	err := DeleteKongConfig(kong_url, "test1")
	if nil != err {
		panic(err)
	}
}
