package llog

import "errors"

const (
	FORMAT_CONSOLE = "console"
	OUTPUT_ERR     = "stderr"
	OUTPUT_OUT     = "stdout"
)

type Options struct {
	OutputPaths      []string `json:"outputPaths" mapstructure:"outputPaths"`
	ErrorOutputPaths []string `json:"errorOutputPaths" mapstructure:"errorOutputPaths"`
	LogLevel         Level    `json:"logLevel" mapstructure:"logLevel"`
	Format           string   `json:"format" mapstructure:"format"`
	Name             string   `json:"name" mapstructure:"name"`
}

type Option interface {
	apply(*Options) error
}

type optionFunc func(*Options) error

func (f optionFunc) apply(o *Options) error {
	return f(o)
}

func New(opts ...Option) *Options {
	options := &Options{
		OutputPaths:      []string{OUTPUT_OUT},
		ErrorOutputPaths: []string{OUTPUT_ERR},
		LogLevel:         Info,
		Format:           FORMAT_CONSOLE,
		Name:             "",
	}
	for _, opt := range opts {
		if err := opt.apply(options); err != nil {
			// TODO: do not panic
			panic("Log new error")
		}
	}
	return options
}

func WithLogLevel(level Level) Option {
	return optionFunc(func(o *Options) error {
		o.LogLevel = level
		return nil
	})
}

func WithOutputPaths(paths []string) Option {
	return optionFunc(func(o *Options) error {
		for _, v := range paths {
			o.OutputPaths = append(o.OutputPaths, v)
		}
		return nil
	})
}

func WithErrorOutputPaths(paths []string) Option {
	return optionFunc(func(o *Options) error {
		for _, v := range paths {
			o.ErrorOutputPaths = append(o.ErrorOutputPaths, v)
		}
		return nil
	})
}

func WithFormat(f string) Option {
	return optionFunc(func(o *Options) error {
		if f == "" {
			return errors.New("format should not be empty")
		}
		o.Format = f
		return nil
	})
}

func WithName(n string) Option {
	return optionFunc(func(o *Options) error {
		if n == "" {
			return errors.New("name should not be empty")
		}
		o.Name = n
		return nil
	})
}
