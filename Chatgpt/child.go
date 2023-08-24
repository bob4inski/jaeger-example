package main

import (
    "context"
    "github.com/opentracing/opentracing-go"
    "github.com/opentracing/opentracing-go/log"
    //"github.com/uber/jaeger-client-go"
    jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
    // Инициализация Jaeger Tracer
    cfg := jaegercfg.Configuration{
        ServiceName: "child-service",
        Sampler: &jaegercfg.SamplerConfig{
            Type:  "const",
            Param: 1,
        },
        Reporter: &jaegercfg.ReporterConfig{
            LogSpans:           true,
            LocalAgentHostPort: "localhost:6831", // адрес и порт Jaeger Agent
        },
    }

    tracer, closer, err := cfg.NewTracer()
    if err != nil {
        // обработка ошибки
    }
    defer closer.Close()

    // Установка глобального трейсера
    opentracing.SetGlobalTracer(tracer)

    // Получение контекста с родительским спаном
    ctx := opentracing.GlobalTracer().StartSpan("child-service").Context()

    // Вызов функции, которая будет использовать родительский спан
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    // Извлечение родительского спана из контекста
    span, _ := opentracing.StartSpanFromContext(ctx, "child-span")
    defer span.Finish()

    // Добавление логов в спан
    span.LogFields(
        log.String("event", "child-span-started"),
    )

    // Ваш код для обработки запроса
    // ...
}
