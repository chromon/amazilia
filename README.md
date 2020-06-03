# Amazilia

Amazilia is a simple HTTP web framework.

## Quick start

example.go

```
package main

import (
    "github.com/chromon/amazilia"
)

func main() {
    r := ama.New()
    r.GET("/index", func(c *ama.Context) {
        c.HTML(http.StatusOK, "<h2>Hello Ama</h2>", nil)
    })
}
```

run example.go and visit `localhost:8080/index`  on browser

```
Hello Ama
```

## LICENSE
Amazilia is distributed under the terms of the GPL-3.0 License.