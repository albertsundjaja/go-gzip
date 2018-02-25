package go_gzip

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"strings"
	"os"
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func (obj goGzip) StaticFilesHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// get the router params *filepath
	req.URL.Path = ps.ByName("filepath")

	splitPath := strings.Split(req.URL.Path, "/")
	splitFilename := strings.Split(splitPath[len(splitPath)-1 ], ".")
	splitExtension := splitFilename[len(splitFilename)-1]

	// only gzip js and css file
	if splitExtension == "js" || splitExtension == "css" {
		//check if header supports gzip compression
		gzipSupport := req.Header.Get("Accept-Encoding")
		var flagGzip bool = false
		if strings.Contains(gzipSupport, "gzip") && strings.Contains(gzipSupport, "deflate") {
			flagGzip = true
		}

		// if gzip is supported serve the gz version
		if flagGzip {
			serveFilesWithMode(obj.ServeMode, obj.ResourceFolder, req.URL.Path, splitExtension, w, req)
			return


			//gzip is not supported serve normal version
		} else {
			fileserver := http.FileServer(http.Dir("resource"))
			fileserver.ServeHTTP(w, req)
			return
		}

		// not js/css requested serve normally
	} else {
		fileserver := http.FileServer(http.Dir("resource"))
		fileserver.ServeHTTP(w, req)
		return
	}
}

// depening on mode selected we might need to create gzip on the fly
func serveFilesWithMode(mode int, resourceFolder string, path string, extension string, w http.ResponseWriter, req *http.Request) {
	filePath := resourceFolder + path + ".gz"

	switch mode {
	case MODE_SERVE_ORIGINAL_IF_NOT_EXIST:
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// gz file doesnt exist break
			break
		} else {
			// add the .gz at the end of the file
			req.URL.Path = req.URL.Path + ".gz"
			// set header to correct encoding
			w.Header().Set("Content-Encoding", "gzip")
			break
		}

	case MODE_CREATE_IF_NOT_EXIST:
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			err := createNewGzipFile(resourceFolder, path)
			// failed to create gzip
			if err != nil {
				http.Error(w, "File not found / Server error", http.StatusNotFound)
				return
			}
		}
		// add the .gz at the end of the file
		req.URL.Path = req.URL.Path + ".gz"
		// set header to correct encoding
		w.Header().Set("Content-Encoding", "gzip")
		break

	case MODE_ASSUME_EXIST:
		// add the .gz at the end of the file
		req.URL.Path = req.URL.Path + ".gz"
		// set header to correct encoding
		w.Header().Set("Content-Encoding", "gzip")
		break
	}


	// set content type, if not it will automatically set as app/gzip
	if (extension == "js") {
		w.Header().Set("Content-Type", "application/javascript")
	} else if (extension == "css"){
		w.Header().Set("Content-Type", "text/css")
	}

	fileserver := http.FileServer(http.Dir(resourceFolder))
	fileserver.ServeHTTP(w, req)
}

func createNewGzipFile(resourceFolder string, path string) error {

	fileContent, err := ioutil.ReadFile(resourceFolder + path)
	if err != nil {
		return err
	}

	// zip the file content to a buffer
	var b bytes.Buffer
	wGzip := gzip.NewWriter(&b)
	wGzip.Write(fileContent)
	wGzip.Flush()
	wGzip.Close()

	// save buffer to file
	err = ioutil.WriteFile(resourceFolder + path + ".gz", b.Bytes(), 0755)
	if err != nil {
		return err
	}

	return nil
}