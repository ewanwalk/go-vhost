# go-vhost
A lightweight implementation of adding hostname routing to your project, inspired by gorilla mux. The goal of this package
is simply to provide an extremely simple method of routing based on hosts without the added overhead of other features.


#### Install

```
go get -u github.com/ewanwalk/go-vhost
```

#### Usage

net/http server:
```go
func main() {
    router := vhost.New()

    router.Handler(myHandler, "domain.com")
    
    http.Handle("/", router)
}
```

#### Configuration
Configuration is presently limited due to the simplistic nature of this package.


**Strict**
Determines if you want to ignore `www.` in hostnames or not

Default: `false`

Usage:
```go
func main() {
    router := vhost.New()
    router.Strict = true

    router.Handler(myHandler, "domain.com")
    
    http.Handle("/", router)
}
```