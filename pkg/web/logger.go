package web

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

func init() {
	DefaultLogger = RequestLogger(NewLogger(os.Stdout))
}

var DefaultLogger func(next http.Handler) http.Handler

func Logger(next http.Handler) http.Handler {
	return DefaultLogger(next)
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	data *responseData
}

func (w *loggingResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.data.size += size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.data.status = statusCode
}

func RequestLogger(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Printf("err: %v, trace: %s\n", err, debug.Stack())
				}
			}()
			lrw := loggingResponseWriter{
				ResponseWriter: w,
				data: &responseData{
					status: 200,
					size:   0,
				},
			}
			next.ServeHTTP(&lrw, r)
			if 400 <= lrw.data.status && lrw.data.status <= 599 {
				LogEntry(logger, "[error] ", lrw.data.status, r)
				return
			}
			LogEntry(logger, " [info] ", lrw.data.status, r)
			return
		}
		return http.HandlerFunc(fn)
	}
}

func LogEntry(logger *log.Logger, prefix string, code int, r *http.Request) {
	format, values := "%s - - [%s] \"%s %s\t%s\"\t%d %s\n", []interface{}{
		r.RemoteAddr,
		time.Now().Format(time.RFC1123Z),
		r.Method,
		r.URL.EscapedPath(),
		r.Proto,
		code,
		http.StatusText(code),
	}
	logger.SetPrefix(prefix)
	logger.Printf(format, values...)
}
