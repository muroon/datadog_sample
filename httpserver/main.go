package main

import (
	"datadog_sample/httpserver/usecases"
	"log"
	"net/http"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// Datadog
	tracer.Start(
		tracer.WithEnv("sample"),
	)
	defer tracer.Stop()

	err := usecases.Init()
	defer usecases.End()
	if err != nil {
		log.Fatalf("init error: %v", err)
		return
	}

	// http handle
	mux := httptrace.NewServeMux(
		httptrace.WithServiceName("web-service"),
	) // http handler for Datadog
	mux.HandleFunc("/grpc/", usecases.GrpcList)
	mux.HandleFunc("/grpc/post", usecases.GrpcPost)
	mux.HandleFunc("/custom/", usecases.Custom)
	mux.HandleFunc("/db/", usecases.DBList)
	mux.HandleFunc("/db/post", usecases.DBPost)
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("http listenServe error: %v", err)
	}
}
