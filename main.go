package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"time"

	"reflect"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {

	//http.Handle("/", http.HandlerFunc(hw))
	//http.Handle("/", messagehander("Hello World"))
	//http.Handle("/", middlware(messagehander("Hello World")))
	//http.Handle("/", middlware(http.HandlerFunc(hw)))
	//http.Handle("/", timing(middlware(http.HandlerFunc(hw))))
	http.Handle("/", decorate(http.HandlerFunc(hw), middlware, timing))

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func decorate(f http.Handler, d ...func(http.Handler) http.Handler) http.Handler {
	decorated := f
	for _, decorateFn := range d {
		fmt.Printf("Decorating %v", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		fmt.Printf(" with %v\n", runtime.FuncForPC(reflect.ValueOf(decorateFn).Pointer()).Name())
		decorated = decorateFn(decorated)
	}
	return decorated
}

func hw(w http.ResponseWriter, req *http.Request) {
	log.Println("001")
	// time.Sleep(time.Second * 1)
	w.Write([]byte("Hello World"))
	log.Println("002")
}

func messagehander(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(message))
	})
}

func timing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("timer start")
		defer func(start time.Time) {
			fmt.Println("timer end", time.Since(start).Nanoseconds())

		}(time.Now())
		next.ServeHTTP(w, req)
	})
}

func middlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// fmt.Printf("%+v\n", req)

		// requestDump, err := httputil.DumpRequest(req, true)
		// if err != nil {
		// fmt.Println(err)
		// }
		// fmt.Println(string(requestDump))

		headers := make(map[string]string)

		for k, v := range req.Header {
			headers[k] = strings.Join(v, ",")
		}

		log.Println("Start")
		next.ServeHTTP(w, req)
		log.Println("Stop")

		log.WithFields(log.Fields{
			"host":       req.Host,
			"requestUri": req.RequestURI,
			"remoteAddr": req.RemoteAddr,
			"method":     req.Method,
			"headers":    headers,
		}).Info()

	})
}
