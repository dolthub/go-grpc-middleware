// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_opentracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

var (
	defaultOptions = &options{
		filterOutFunc: nil,
		tracerFactory: nil,
	}
)

// FilterFunc allows users to provide a function that filters out certain methods from being traced.
//
// If it returns false, the given request will not be traced.
type FilterFunc func(ctx context.Context, fullMethodName string) bool

// UnaryRequestHandlerFunc is a custom request handler
type UnaryRequestHandlerFunc func(span opentracing.Span, req interface{})

// OpNameFunc is a func that allows custom operation names instead of the gRPC method.
type OpNameFunc func(method string) string

type TracerFactory func(ctx context.Context) opentracing.Tracer

type options struct {
	filterOutFunc           FilterFunc
	tracerFactory           TracerFactory
	traceHeaderName         string
	unaryRequestHandlerFunc UnaryRequestHandlerFunc
	opNameFunc              OpNameFunc
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	if optCopy.tracerFactory == nil {
		optCopy.tracerFactory = defaultTracerFactory
	}
	if optCopy.traceHeaderName == "" {
		optCopy.traceHeaderName = "uber-trace-id"
	}
	return optCopy
}

type Option func(*options)

// WithFilterFunc customizes the function used for deciding whether a given call is traced or not.
func WithFilterFunc(f FilterFunc) Option {
	return func(o *options) {
		o.filterOutFunc = f
	}
}

// WithTraceHeaderName customizes the trace header name where trace metadata passed with requests.
// Default one is `uber-trace-id`
func WithTraceHeaderName(name string) Option {
	return func(o *options) {
		o.traceHeaderName = name
	}
}

// WithTracer sets a custom tracer to be used for this middleware, otherwise the opentracing.GlobalTracer is used.
func WithTracer(tracer opentracing.Tracer) Option {
	return func(o *options) {
		o.tracerFactory = constantFactory(tracer)
	}
}

// WithTracerFactory sets a factory to get the tracer to be used for this middleware.
func WithTracerFactory(factory TracerFactory) Option {
	return func(o *options) {
		o.tracerFactory = factory
	}
}

// WithUnaryRequestHandlerFunc sets a custom handler for the request
func WithUnaryRequestHandlerFunc(f UnaryRequestHandlerFunc) Option {
	return func(o *options) {
		o.unaryRequestHandlerFunc = f
	}
}

// WithOpName customizes the trace Operation name
func WithOpName(f OpNameFunc) Option {
	return func(o *options) {
		o.opNameFunc = f
	}
}

func constantFactory(tracer opentracing.Tracer) TracerFactory {
	return func(ctx context.Context) opentracing.Tracer {
		return tracer
	}
}

func defaultTracerFactory(ctx context.Context) opentracing.Tracer {
	return opentracing.GlobalTracer()
}
