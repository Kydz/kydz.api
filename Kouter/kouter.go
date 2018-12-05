package Kouter

import (
	"errors"
	"fmt"
	"github.com/Kydz/kydz.api/Kogger"
	"github.com/Kydz/kydz.api/Konfigurator"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const AnalyzePattern = "{:\\w+}"
const WholeMatchPattern = "{?:?\\w+}?"

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

type Kiddleware func(http.HandlerFunc) http.HandlerFunc

func (k Kouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer unexpectedPanic(w)
	path := r.URL.Path[1:]
	method := r.Method


	if method == http.MethodOptions {
		corsResponse(w)
	} else {
		Kogger.Debug("Router processed: %+v", *k.routes)
		for i, route := range *k.routes {
			if kandler, ok := route.handlers[method]; ok && route.isWholeMatch(path) {
				w.Header().Set("Access-Control-Allow-Origin", kon.Site)
				route.fillParamsWithValue(path)
				setCurrentRoute(&route)
				Kogger.Debug("Kandler processed: %+v", kandler)
				handler := kandler.handleKiddleware()
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

func (k *Kouter) Get(pattern string, hf http.HandlerFunc) *Kandler {
	kandler, _ := k.addRoute(pattern, http.MethodGet, hf)
	return kandler
}

func (k *Kouter) Post(pattern string, hf http.HandlerFunc) *Kandler {
	kandler, _ := k.addRoute(pattern, http.MethodPost, hf)
	return kandler
}

func (k *Kouter) Put(pattern string, hf http.HandlerFunc) *Kandler {
	kandler, _ := k.addRoute(pattern, http.MethodPut, hf)
	return kandler
}

func (k *Kouter) Delete(pattern string, hf http.HandlerFunc) *Kandler {
	kandler, _ := k.addRoute(pattern, http.MethodDelete, hf)
	return kandler
}

func (k *Kouter) addRoute(pattern string, method string, hf http.HandlerFunc) (*Kandler, error) {
	kandler := newKandler(hf)
	if len(*k.routes) == 0 {
		*k.routes = append(*k.routes, *newR(pattern, method, kandler))
		return kandler, nil
	} else {
		for i, route := range *k.routes {
			if route.isWholeMatch(pattern) {
				if _, ok := route.handlers[method]; ok {
					panic("trying to add duplicate routers, pattern: " + pattern + ", method: " + method)
				} else {
					(&(*k.routes)[i]).handlers[method] = kandler
					return kandler, nil
				}
			} else {
				if i == len(*k.routes) - 1 {
					*k.routes = append(*k.routes, *newR(pattern, method, kandler))
					return kandler, nil
				}
			}
		}
	}
	return nil, errors.New("no router added")
}

func newR(pattern string, method string, kandler *Kandler) *Route {
	r := new(Route)
	r.handlers = make(map[string]*Kandler)
	r.paramsHolder = make(map[int]string)
	r.Params = make(map[string]string)
	r.pattern = pattern
	r.handlers[method] = kandler
	r.analyzePattern()
	return r
}

type Route struct {
	pattern      string
	handlers     map[string]*Kandler
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

func (r *Route) isWholeMatch(path string) bool {
	match, _ := regexp.MatchString(r.wholeMatcher, path)
	return match
}

func (r *Route) fillParamsWithValue(path string) {
	values := strings.Split(path, "/")
	p := make(map[string]string)
	for i, field := range r.paramsHolder {
		p[field] = values[i]
	}
	r.Params = p
}

func (k *Kandler) Kware(kids... Kiddleware) {
	for _, kid := range kids {
		k.Kiddlewares = append(k.Kiddlewares, kid)
	}
}

func (k *Kandler) handleKiddleware() http.HandlerFunc{
	hf := k.HandlerFunc
	for _, kid := range k.Kiddlewares {
		hf = kid(hf)
	}
	return hf
}

type Kandler struct {
	HandlerFunc http.HandlerFunc
	Kiddlewares []Kiddleware
}

func newKandler(hf http.HandlerFunc) *Kandler {
	k := new(Kandler)
	k.HandlerFunc = hf
	k.Kiddlewares = make([]Kiddleware, 0)
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, x-auth-token")
	_, err := fmt.Fprint(w)
	if err != nil {
		log.Printf("CORS Error: %+v", err)
	}
}
