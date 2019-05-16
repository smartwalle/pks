module greeter

go 1.12

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.0-20181115231424-8e868ca12c0f

require (
	github.com/golang/protobuf v1.3.1
	github.com/micro/examples v0.1.0
	github.com/micro/go-micro v1.1.0
	github.com/micro/go-plugins v1.1.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/smartwalle/jaeger4go v1.0.0
	github.com/smartwalle/pks v1.0.0
	github.com/uber/jaeger-client-go v2.16.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3
	google.golang.org/grpc v1.20.1
)
