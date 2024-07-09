package main

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var projectNameToProjectUrl map[string]string = map[string]string{
	"receiving-api": "https://receivingApiUrl",
}

func proxyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	apiName := r.Header.Get("api-name")
	if apiName == "" {
		http.Error(w, "api-name header cannot be empty!", http.StatusBadRequest)
		return
	}
	apiUrl, ok := projectNameToProjectUrl[apiName]
	if !ok {
		http.Error(w, apiName+"not found!", http.StatusBadRequest)
		return
	}

	path := ps.ByName("path")
	fullUrl := apiUrl + path

	req, err := http.NewRequest(r.Method, fullUrl, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func registerHandlers(router *httprouter.Router, path string, handler httprouter.Handle) {
	router.GET(path, handler)
	router.POST(path, handler)
	router.PUT(path, handler)
	router.DELETE(path, handler)
	router.PATCH(path, handler)
	router.OPTIONS(path, handler)
	router.HEAD(path, handler)
}

func main() {
	router := httprouter.New()
	registerHandlers(router, "/*path", proxyHandler)

	http.ListenAndServe(":8080", router)
}
