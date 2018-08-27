package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

func main() {
	http.Handle("/", errorHandler(handleJSRequest))

	hostport := fmt.Sprintf("%s:%s", getHost(), getPort())
	logrus.WithField("hostport", hostport).Info("starting server")
	if err := http.ListenAndServe(hostport, nil); err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("server failed")
	}
}

func handleJSRequest(resp http.ResponseWriter, req *http.Request) error {
	filePath := makePathAbsolute(req.URL.Path)

	if !strings.HasSuffix(filePath, ".js") {
		http.Error(resp, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}

	file, err := http.Dir("./").Open(path.Clean(filePath))
	if err != nil {
		return errors.Wrap(err, "failed to open js file")
	}

	js, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "failed to read js file")
	}

	vm := otto.New()
	vm.Set("__request_method", func(call otto.FunctionCall) otto.Value {
		result, _ := otto.ToValue(req.Method)
		return result
	})
	vm.Set("__request_path", func(call otto.FunctionCall) otto.Value {
		result, _ := otto.ToValue(req.RequestURI)
		return result
	})
	vm.Set("__response_println", func(call otto.FunctionCall) otto.Value {
		text, _ := call.Argument(0).ToString()
		fmt.Fprint(resp, text)
		return otto.NullValue()
	})
	if _, err := vm.Run(bootstrapJS); err != nil {
		return errors.Wrap(err, "failed to run bootstrap js code")
	}
	if _, err := vm.Run(string(js)); err != nil {
		return errors.Wrap(err, "failed to run file js code")
	}
	return nil
}

func errorHandler(delegate func(resp http.ResponseWriter, req *http.Request) error) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if err := delegate(resp, req); err != nil {
			logrus.WithError(err).Error("failed to process request")
			http.Error(resp, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})
}

func makePathAbsolute(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}

func getHost() string {
	host, found := os.LookupEnv("HOST")
	if !found {
		return "127.0.0.1"
	}
	return host
}

func getPort() string {
	port, found := os.LookupEnv("PORT")
	if !found {
		return "8080"
	}
	return port
}

const bootstrapJS = `
function require(module) {
	switch (module) {
	case 'http/v3/request':
		return {
			getMethod: __request_method,
			getPath: __request_path
		};
		break;
	case 'http/v3/response':
		return {
			println: __response_println
		}
	}
};
`
