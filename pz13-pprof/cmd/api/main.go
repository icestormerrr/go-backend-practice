package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // регистрирует /debug/pprof/* на DefaultServeMux
	"runtime"

	"github.com/icestormerrr/pz13-pprof/internal/work"
)

func main() {
	enableLocks()

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		n := 35

		defer work.TimeIt("Fib(35)")()
		res := work.Fib(n)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(w, "%d\n", res)
	})

	http.HandleFunc("/work-fast", func(w http.ResponseWriter, r *http.Request) {
		n := 35

		defer work.TimeIt("FibFast(35)")()
		res := work.FibFast(n)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(w, "%d\n", res)
	})

	log.Println("Server on :8080; pprof on /debug/pprof/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enableLocks() {
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
}
