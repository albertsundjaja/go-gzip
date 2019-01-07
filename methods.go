package go_gzip

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"strings"
	"os"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	filepath2 "path/filepath"
	"fmt"
)

func (obj goGzip) StaticFilesHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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
			serveFiles(obj.ResourceFolder, req.URL.Path, splitExtension, w, req)
			return

			//gzip is not supported serve normal version
		} else {
			fileserver := http.FileServer(http.Dir(obj.ResourceFolder))
			fileserver.ServeHTTP(w, req)
			return
		}

		// not js/css requested serve normally
	} else {
		fileserver := http.FileServer(http.Dir(obj.ResourceFolder))
		fileserver.ServeHTTP(w, req)
		return
	}
}

// preGzip everything in the resource folder
func (obj goGzip) ProcessResourceFolder() {
	fmt.Println("gzipping resource folder")

	checkShouldZip(obj.ResourceFolder)

	fmt.Println("gzip finish")
}

func checkShouldZip(parentFolder string) {
	files, err := ioutil.ReadDir(parentFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _,file := range files {
		if !file.IsDir() && (filepath2.Ext(file.Name()) == ".css" || filepath2.Ext(file.Name()) == ".js") {
			if err := createNewGzipFile(parentFolder, file.Name()); err != nil {
				log.Fatal(err)
			}
		} else if (file.IsDir()) {
			checkShouldZip(filepath2.Join(parentFolder, file.Name()))
		}
	}
}

func serveFiles(resourceFolder string, path string, extension string, w http.ResponseWriter, req *http.Request) {
	filePath := resourceFolder + path + ".gz"

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// do nothing .gz doesnt exist
	} else {
		// add the .gz at the end of the file
		req.URL.Path = req.URL.Path + ".gz"
		// set header to correct encoding
		w.Header().Set("Content-Encoding", "gzip")
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

	fileContent, err := ioutil.ReadFile(filepath2.Join(resourceFolder, path))
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
	err = ioutil.WriteFile(filepath2.Join(resourceFolder, path + ".gz"), b.Bytes(), 0755)
	if err != nil {
		return err
	}

	return nil
}