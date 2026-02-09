package apm

import (
	"log"

	appconfig "github.com/afandimsr/cashbook-backend/internal/config"
	"go.elastic.co/apm/v2"
)

func Init(cfg *appconfig.Config) *apm.Tracer {
	// tracer := apm.DefaultTracer

	// APM setup
	tracer, err := apm.NewTracer(
		cfg.ElasticApm.ServiceName,
		cfg.ElasticApm.ServiceVersion,
	)

	if err != nil {
		panic(err)
	}
	apm.SetDefaultTracer(tracer)

	if !tracer.Active() {
		log.Fatal("APM tracer not active — check ENV")
	}

	log.Println("APM started:", cfg.ElasticApm.ServiceName)
	return tracer
}
