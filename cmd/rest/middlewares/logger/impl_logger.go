package logger

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// ---------------------------------------------------------------------------------------------
// Logger is a post-middleware that logs the request as it goes in and the response as it goes out.
type Log struct {
	IP	    string
	Method  string
	Path    string
	Status  int
	Errors  []error
}
func (l *Log) Display() string {
	return fmt.Sprintf("[%s] %s %d | ip: %s | errors: %v", l.Method, l.Path, l.Status, l.IP, l.Errors)
}

func LoggerDefault(h http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		// call the handler and middleware chain
		h.ServeHTTP(w, r)

		// post middleware
		// -> fetch the log from the request's context and response writer
		var lg Log

		// -> request info
		lg.IP = r.RemoteAddr
		lg.Method = r.Method
		lg.Path = r.URL.Path

		// -> response info
		// status := w.(interface{Status() int}).Status()
		// lg.Status = status
		lg.Errors, _ = r.Context().Value(CtxKeyLogger).([]error)

		// -> log: [method] path status | ip : ... | errors: ...
		log.Println(lg.Display())
	})
}




// ---------------------------------------------------------------------------------------------
// Context: allows to connect handler layer with the middleware layer for the purpose of logging.
type ContextKey int
const (
	CtxKeyLogger ContextKey = iota
)

func Errors(r *http.Request, errs ...error) {
	// create new context and update the request's context with it
	ctx := context.WithValue(r.Context(), CtxKeyLogger, errs)
	*r = *r.WithContext(ctx)
}