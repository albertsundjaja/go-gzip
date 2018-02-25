# go-gzip
A go helper to serve gzipped static contents in website server. This require the use of [Julien Scmidt's router](https://github.com/julienschmidt/httprouter)

Gzip static contents (mainly css & js) can speed up page loading time by a lot. Read the detail [here](https://betterexplained.com/articles/how-to-optimize-your-site-with-gzip-compression/)

This go-gzip is, in essence, a Julien Schmidt's router handler which serve static contents

At this moment, only .js and .css files will be gzipped

## How to use
#### This is how to serve files normally in httprouter
```go
router := httprouter.New()
// with the httprouter, this is how you serve static contents normally
router.ServeFiles("/rsc/static/*filepath", http.Dir("resource/"))
```

#### With go-gzip
```go
router := httprouter.New()

// prepare the settings
staticHandler := goGzip.CreateNew()
// this is the location of the parent folder of the resources
staticHandler.ResourceFolder = "resource"
/*
this is the mode on what to do if a file doesnt exist
MODE_CREATE_IF_NOT_EXIST : this will create the gzip files on the fly if it doesnt exist
MODE_SERVE_ORIGINAL_IF_NOT_EXIST : like what the name suggest, 
                                   it will serve the normal file if .gz files not there
MODE_ASSUME_EXIST : this mode assumes you have created all .gz files manually,
                    CAUTION! if the .gz files dont exist, it will return an error 404
*/
staticHandler.ServeMode = goGzip.MODE_CREATE_IF_NOT_EXIST

router.GET("/rsc/static/*filepath", staticHandler.StaticFilesHandler)
```
