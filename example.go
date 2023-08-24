package main

import (
"fmt"
"net/http"
"os"

"github.com/opentracing/opentracing-go"
"github.com/opentracing/opentracing-go/ext"
"github.com/opentracing/opentracing-go/log"
"github.com/uber/jaeger-client-go"
"github.com/uber/jaeger-client-go/config"
)

func main() {
// Создание и настройка трассировщика Jaeger
	cfg := config.Configuration{
		ServiceName: "your-service-name",
		Sampler: &config.SamplerConfig{
		Type:  jaeger.SamplerTypeConst,
		Param: 1,
		},
		Reporter: &config.ReporterConfig{
		LogSpans:            true,
		LocalAgentHostPort:  "localhost:6831",
		BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
		if err != nil {
			fmt.Printf("Failed to create tracer: %v\n", err)
			return
		}
	defer closer.Close()

// Установка глобального трассировщика
	opentracing.SetGlobalTracer(tracer)

// Создание HTTP-обработчика
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// Извлечение родительского спана из HTTP-заголовков
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := tracer.StartSpan("your-operation-name", ext.RPCServerOption(spanCtx))
	defer span.Finish()

// Пример работы с родительским спаном
	span.LogFields(
		log.String("event", "your-event"),
		log.Int("value", 42),
	)

// Добавление информации о родительском спане в HTTP-заголовки
	tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

// Ваш код обработки запроса

// Отправка ответа
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
	})

// Запуск HTTP-сервера
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
