package Kouter

import (
	"fmt"
	"log"
	"net/http"
)

const MaxRouterNumber = 20

type Kandler func(w http.ResponseWriter, r *http.Request)

func GetK() (k Kouter) {
	k.routes = make(map[string]map[string]Kandler, MaxRouterNumber)
	return k
}

type Kouter struct {
	routes map[string]map[string]Kandler
}

func (k Kouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer unexpectedPanic(w)
	path := r.URL.Path
	method := r.Method
	if handler, ok := k.routes[path][method]; ok {
		http.HandlerFunc(handler).ServeHTTP(w, r)
	} else {
		if method == http.MethodOptions {
			corsResponse(w)
		} else {
			http.NotFound(w, r)
		}
	}
}

func (k *Kouter) Get(path string, handler Kandler) {
	addRoute(k, path, http.MethodGet, handler)
}
func (k *Kouter) Post(path string, handler Kandler) {
	addRoute(k, path, http.MethodPost, handler)
}
func (k *Kouter) Put(path string, handler Kandler) {
	addRoute(k, path, http.MethodPut, handler)
}
func (k *Kouter) Delete(path string, handler Kandler) {
	addRoute(k, path, http.MethodDelete, handler)
}

func addRoute(k *Kouter, path string, method string, handler Kandler) {
	if _, ok := k.routes[path]; !ok {
		p := make(map[string]Kandler)
		k.routes[path] = p
	}
	k.routes[path][method] = handler
}

func unexpectedPanic(w http.ResponseWriter) {
	if recoveredErr := recover(); recoveredErr != nil {
		log.Print(recoveredErr)
		http.Error(w, "Unexpected Error", 500)
	}
}

func corsResponse(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
	w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprint(w)
}