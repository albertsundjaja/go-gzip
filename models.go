package go_gzip

var MODE_SERVE_ORIGINAL_IF_NOT_EXIST = 0
var MODE_CREATE_IF_NOT_EXIST = 1
var MODE_ASSUME_EXIST = 2

type goGzip struct {
	ServeMode int
	ResourceFolder string
}

//constructor
func CreateNew() goGzip {
	var tempGoGzip goGzip
	tempGoGzip.ServeMode = MODE_SERVE_ORIGINAL_IF_NOT_EXIST
	return tempGoGzip
}