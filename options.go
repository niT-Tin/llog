package llog

type Options struct {
	OutputPaths      []string `json:"outputPaths" mapstructure:"outputPaths"`
	ErrorOutputPaths []string `json:"errorOutputPaths" mapstructure:"errorOutputPaths"`
	// Level
}
