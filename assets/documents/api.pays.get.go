package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var (
		request, _ = http.NewRequest(http.MethodGet, "{GATEWAY}/api/pays", nil)
	)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}
