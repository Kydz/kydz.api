# kydz.api

API server for http://kydz.site in Go

## About Kouter

A light weight Router that supports named parameters.

### Usage

Simply call Kouter.NewK() will get you a new instance of Kouter.
Then you are good to go.

```
k := Kouter.NewK();
```

#### Add Routes
Kouter provides a straight forward mean to add route, based on
`Http Method` you are going to coupe with, you can call .Get(), .Post(), etc...

Say we have a `Restful` resource: books, then we can construct the
following:
```
k.Get("books", GetBooksListHandler)
k.Post("books", CreateNewBookHandler)
k.Get("books/{:id}", GetBookHandler)
k.Put("books/{:id}", UpdateBookHandler)
k.Del("books/{:id}", DeleteBookHandler)
```

#### Compatibility

For all `***Handler`, you can consider each of them as a [http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc) type,
it is a standard Go type so you can easily migrate your project to
implement Kouter.

#### Named parameter

Guess you have noticed the `{:id}` notation, this is how Kouter handles named
parameters, by registering `"books/{:id}"` into Kouter, Kouter will know
there is a parameter named `id`, and assign the relative value to this
parameter.

You can easily access named parameters via `Kouter.GetCurrentRoute().Param`,
this is a type of [map](https://golang.org/ref/spec#Map_types) that contains
the named parameters of the current route, say you are currently in "https://yourServer.com/books/32",
then you will get `32` by accessing `Kouter.GetCurrentRoute().Param["id"]`

Of course you can add multiple named parameters, like so:
```
k.Get("books/{:publisher}/{:category}", GetBookHandler)
```
Access them with their names, how easy is that.

#### Middleware

You can add middleware by chaining a `Kware` method after add routes, a middleware is typically a function that
 follows the signature: `func SomeMiddleware(next Kouter.Kandler) Kouter.Kandler {}`, here is an example.
 > Kandler is nothing but a re-name of [http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc)

```go
func AuthMiddleware(next Kouter.Kandler) Kouter.Kandler {
    return func(w http.ResponseWritter, r *http.Request) {
        // get a user and some do auth check
        ...
        if user.CanLogin() {
            next(w, r)
        } else {
            http.Error(w, "access deny", http.StatusUnauthorized)
        }
    }
}

func main() {
    k := Kouter.NewK()
    k.Post("books/{:id}", CreateNewBookHandler).Kware(AuthMiddleware)
    log.Fatal(http.ListenAndServe(":8088", k))
}
```

Also you can add multiple middleware by pass them all into the Kware() method, like so: `Kware(middleware1, middleware2, middleware3...)`

### TODO
- enhance named parameters with customizable regex