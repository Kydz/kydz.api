package Kouter

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const AnalyzePattern = "{:\\w+}"
const WholeMatchPattern = "{?:?\\w+}?"

type Kandler func(w http.ResponseWriter, r *http.Request)

var currentRoute *Route

func setCurrentRoute(r *Route) {
	currentRoute = r
}

func GetCurrentRoute() *Route {
	return currentRoute
}

func NewK() (k Kouter) {
	routes := make([]Route, 0)
	k.routes = &routes
	return k
}

type Kouter struct {
	routes *[]Route
}

func (k Kouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer unexpectedPanic(w)
	path := r.URL.Path[1:]
	method := r.Method

	if method == http.MethodOptions {
		corsResponse(w)
	} else {
		for i, route := range *k.routes {
			if handler, ok := route.handlers[method]; ok && route.IsWholeMatch(path) {
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
				route.FillParamsWithValue(path)
				setCurrentRoute(&route)
				http.HandlerFunc(handler).ServeHTTP(w, r)
				break
			} else {
				if i == len(*k.routes) - 1 {
					http.NotFound(w, r)
				}
			}
		}
	}
}

func (k *Kouter) Get(pattern string, handler Kandler) {
	addRoute(k, pattern, http.MethodGet, handler)
}
func (k *Kouter) Post(pattern string, handler Kandler) {
	addRoute(k, pattern, http.MethodPost, handler)
}
func (k *Kouter) Put(pattern string, handler Kandler) {
	addRoute(k, pattern, http.MethodPut, handler)
}
func (k *Kouter) Delete(pattern string, handler Kandler) {
	addRoute(k, pattern, http.MethodDelete, handler)
}

func addRoute(k *Kouter, pattern string, method string, handler Kandler) {
	if len(*k.routes) == 0 {
		*k.routes = append(*k.routes, *newR(pattern, method, handler))
	} else {
		for i, route := range *k.routes {
			if route.IsWholeMatch(pattern) {
				if _, ok := route.handlers[method]; ok {
					panic("trying to add duplicate routers, pattern: " + pattern + ", method: " + method)
				} else {
					route.handlers[method] = handler
				}
			} else {
				if i == len(*k.routes) - 1 {
					*k.routes = append(*k.routes, *newR(pattern, method, handler))
				}
			}
		}
	}

}

func newR(pattern string, method string, handler Kandler) *Route {
	r := new(Route)
	r.handlers = make(map[string]Kandler)
	r.paramsHolder = make(map[int]string)
	r.Params = make(map[string]string)
	r.pattern = pattern
	r.handlers[method] = handler
	r.analyzePattern()
	return r
}

type Route struct {
	pattern      string
	handlers     map[string]Kandler
	paramsHolder map[int]string
	wholeMatcher string
	Params       map[string]string
}

func (r *Route) analyzePattern() {
	re, _ := regexp.Compile(AnalyzePattern)
	frags := strings.Split(r.pattern, "/")
	fields := make(map[int]string)
	pm := make([]string, len(frags))
	for i, frag := range frags {
		res := re.FindAllString(frag, -1)
		if len(res) == 1 {
			fields[i] = res[0][2 : len(res[0])-1]
			pm[i] = WholeMatchPattern
		} else {
			pm[i] = frag
		}
	}
	r.wholeMatcher = "^" + strings.Join(pm, "/") + "$"
	r.paramsHolder = fields
}

func (r *Route) IsWholeMatch(path string) bool {
	match, _ := regexp.MatchString(r.wholeMatcher, path)
	return match
}

func (r *Route) FillParamsWithValue(path string) {
	values := strings.Split(path, "/")
	p := make(map[string]string)
	for i, field := range r.paramsHolder {
		p[field] = values[i]
	}
	r.Params = p
}

func unexpectedPanic(w http.ResponseWriter) {
	if recoveredErr := recover(); recoveredErr != nil {
		log.Print(recoveredErr)
		http.Error(w, "Unexpected Error", 500)
	}
}

func corsResponse(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprint(w)
}
