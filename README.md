# PROD Framework

Эта кароч либа для моей команды на проде ну чтобы не на джаве ну кароч

ну типа пишешь ручки а она типа генерит сваггер и работает с ашипками ну кароч

```go
package main

import (
	"github.com/TaeKwonZeus/pf"
	"net/http"
)

func Ping(w pf.ResponseWriter[string], r *http.Request, body struct{}) error {
    return w.JSON("Pong!")
}

func main() {
    r := pf.NewRouter()
    pf.Get(r, "/ping", Ping)

    err := http.ListenAndServe(":8080", r)
    if err != nil {
        slog.Error("Error starting server", "err", err)
    }
}
```
