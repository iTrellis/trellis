package cmd

type Option func(*Options)

type Options struct {
	ConfigFile string
}

func ConfigFile(filepath string) Option {
	return func(o *Options) {
		o.ConfigFile = filepath
	}
}
