package config

type SentryConfiguration struct {
	Dsn              string
	Debug            bool
	Environment      string
	TracesSampleRate int
}
