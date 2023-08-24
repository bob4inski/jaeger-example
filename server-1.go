package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"io/ioutil"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	// "go.opentelemetry.io/otel/trace"
)

func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			attribute.String("service.name", "server-1"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "localhost:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func mainHandler(w http.ResponseWriter, r *http.Request) {

    //инициализация подключения
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    shutdown, err := initProvider()
    if err != nil {
      log.Fatal(err)
    }
    defer func() {
      if err := shutdown(ctx); err != nil {
        log.Fatal("failed to shutdown TracerProvider: %w", err)
      }
    }()
    
    tracer := otel.Tracer("server-1")

    // Получаем значение переменной из запроса
    myvar := r.URL.Query().Get("myvar")
      
    // Создаем URL для второго сервера
    url := "http://localhost:6767/second?myvar=" + myvar

    ctx, goto_span := tracer.Start(
      ctx,
      "goto-second-server")
    defer goto_span.End()
    goto_span.SetAttributes(attribute.String("myvar", myvar))

    // Отправляем GET-запрос на второй сервер
    resp, err := http.Get(url)

    if err != nil {
      fmt.Fprintf(w, "Ошибка при отправке запроса на второй сервер: %v", err)
      return
    }
    _, iSpan_goto := tracer.Start(ctx, "go-to" ) //put work code here
		iSpan_goto.End()
    defer resp.Body.Close()

    ctx, child_span := tracer.Start(
      ctx,
      "child-server")
    defer child_span.End()
    
    

    // Читаем ответ от второго сервера
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      fmt.Fprintf(w, "Ошибка при чтении ответа от второго сервера: %v", err)
      return
    }
    
    // _, iSpan_getfrom := tracer.Start(ctx, "get-from" ) //put work code here
    // iSpan_getfrom.SetAttributes(attribute.String("Body", string(body)))
		// iSpan_getfrom.End()

    fmt.Fprintf(w, "Ответ от второго сервера: %s", body)
  }

func main() {

  http.HandleFunc("/main", mainHandler)
  http.ListenAndServe(":5656", nil)
}