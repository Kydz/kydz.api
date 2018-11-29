package Kouter

import (
	"errors"
	"fmt"
	"github.com/Kydz/kydz.api/Konfigurator"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const AnalyzePattern = "{:\\w+}"
const WholeMatchPattern = "{?:?\\w+}?"

type Kandler func(http.ResponseWriter, *http.Request)

type Kiddleware func(Kandler) Kandler

var currentRoute *Route
var kon *Konfigurator.Kon

func init() {
	kon = Konfigurator.GetKon()
}

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
				w.Header().Set("Access-Control-Allow-Origin", kon.Site)
				route.FillParamsWithValue(path)
				setCurrentRoute(&route)
				handler = route.HandleKiddleware(handler)
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

func (k *Kouter) Get(pattern string, handler Kandler) *Route {
	r, _ := k.addRoute(pattern, http.MethodGet, handler)
	return r
}
func (k *Kouter) Post(pattern string, handler Kandler) *Route {
	r, _ := k.addRoute(pattern, http.MethodPost, handler)
	return r
}
func (k *Kouter) Put(pattern string, handler Kandler) *Route {
	r, _ := k.addRoute(pattern, http.MethodPut, handler)
	return r
}
func (k *Kouter) Delete(pattern string, handler Kandler) *Route {
	r, _ := k.addRoute(pattern, http.MethodDelete, handler)
	return r
}

func (k *Kouter) addRoute(pattern string, method string, handler Kandler) (*Route, error) {
	if len(*k.routes) == 0 {
		*k.routes = append(*k.routes, *newR(pattern, method, handler))
		return &(*k.routes)[len(*k.routes) -1], nil
	} else {
		for i, route := range *k.routes {
			if route.IsWholeMatch(pattern) {
				if _, ok := route.handlers[method]; ok {
					panic("trying to add duplicate routers, pattern: " + pattern + ", method: " + method)
				} else {
					route.handlers[method] = handler
					return &route, nil
				}
			} else {
				if i == len(*k.routes) - 1 {
					*k.routes = append(*k.routes, *newR(pattern, method, handler))
					return &(*k.routes)[len(*k.routes) -1], nil
				}
			}
		}
	}
	return nil, errors.New("no router added")
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
	kiddlewares  []Kiddleware
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

func (r *Route) Kware(kid Kiddleware) *Route {
	r.kiddlewares = append(r.kiddlewares, kid)
	return r
}

func (r *Route) HandleKiddleware(k Kandler) Kandler{
	for _, kid := range r.kiddlewares {
		k = kid(k)
	}
	return k
}

func unexpectedPanic(w http.ResponseWriter) {
	if recoveredErr := recover(); recoveredErr != nil {
		log.Print(recoveredErr)
		http.Error(w, "Unexpected Error", 500)
	}
}

func corsResponse(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", kon.Site)
	w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	_, err := fmt.Fprint(w)
	if err != nil {
		log.Printf("CORS Error: %+v", err)
	}
}
