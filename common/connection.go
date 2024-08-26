package common

import (
	"io"
	"net/http"
)

func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	uidCookie, err := request.Cookie("uid")
	if err != nil {
		return
	}
	signCookie, err := request.Cookie("sign")
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}

	req.AddCookie(&http.Cookie{
		Name:  "uid",
		Value: uidCookie.Value,
		Path:  "/",
	})

	req.AddCookie(&http.Cookie{
		Name:  "sign",
		Value: signCookie.Value,
		Path:  "/",
	})

	response, err = client.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
	return
}
