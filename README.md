# go-gzip
A go helper to serve gzipped static contents in website server. This require the use of [Julien Scmidt's router](https://github.com/julienschmidt/httprouter)

Gzip static contents (mainly css & js) can speed up page loading time by a lot. Read the detail [here](https://betterexplained.com/articles/how-to-optimize-your-site-with-gzip-compression/)

This go-gzip is, in essence, a Julien Schmidt's router handler which serve static contents

At this moment, only .js and .css files will be gzipped

## Installation
`go get -u github.com/albertsundjaja/go-gzip`

## How to use
#### This is how to serve files normally in httprouter
```go
router := httprouter.New()
// with the httprouter, this is how you serve static contents normally
router.ServeFiles("/rsc/static/*filepath", http.Dir("resource/"))
```

#### With go-gzip
```go
import (
  goGzip "github.com/albertsundjaja/go-gzip"
)

func main() {
  router := httprouter.New()

  // create new handler
  staticHandler := goGzip.CreateNew()
  // this is the location of the parent folder of the resources
  staticHandler.ResourceFolder = "resource"
  // this will create and rewrite a .gzip version of all .css and .js
  // CAUTION: if you skip this, then the original version will be served
  staticHandler.ProcessResourceFolder()

  router.GET("/rsc/static/*filepath", staticHandler.StaticFilesHandler)
}
```


