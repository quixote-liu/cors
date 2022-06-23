# cors
解决跨域插件

可以对单个API解决跨域问题：
```go
package main

import (
  "http"
  "github.com/quixote-liu/cors"
)

func main() {
  mux := http.NewServeMux()
  mux.HandlerFunc("/hello", cors.WarpH(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello, %s", "world")
  }))
  log.Fatal(http.ListenAndServe(":yourPort", mux))
}
```

也可以对全局的请求设置cors跨域
```go
package main

import "github.com/quixote-liu/cors"

func main() {
  mux := http.NewServeMux()
  mux.HandlerFunc("/hello", cors.WarpH(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello, %s", "world")
  }))
  log.Fatal(http.ListenAndServe(":yourPort", mux))
}


type mux struct {
	*http.ServeMux
}

func NewServerMux() *mux {
	return &mux{ServeMux: http.NewServeMux()}
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ok := cors.Handler(w, r); ok {
    w.WriteHeader(http.StatusNotFound)
    return
  }
	m.ServeMux.ServeHTTP(w, r)
}

```

如果使用中间件，也可以类似的在中间件实现