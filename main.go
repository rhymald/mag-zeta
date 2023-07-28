package main

import (
	"rhymald/mag-zeta/api"
	"rhymald/mag-zeta/connect"
	"fmt"
	"go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/propagation"
) 

func init() {
	fmt.Println()
	fmt.Println("Hello world!..")
	fmt.Println()
}

func main() {
	tp, tpErr := connect.JaegerTraceProvider()
  if tpErr != nil { fmt.Println(tpErr) }
  otel.SetTracerProvider(tp)
  otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	api.RunAPI()
}