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
  // see explanation on ServeMode below
  staticHandler.ServeMode = goGzip.MODE_CREATE_IF_NOT_EXIST

  router.GET("/rsc/static/*filepath", staticHandler.StaticFilesHandler)
}
```

### ServeMode
ServeMode is the action to be taken if `.gz` file version does not yet exist

**MODE\_CREATE\_IF\_NOT\_EXIST** this will create the .gz file from original file on the fly

**MODE\_SERVE\_ORIGINAL\_IF\_NOT\_EXIST** it will serve the normal file if .gz file not found
  
**MODE\_ASSUME\_EXIST** this mode assumes you have created all .gz files manually, *CAUTION!* if the .gz file doesn't exist, it will return an error 404

