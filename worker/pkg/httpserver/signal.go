package httpserver

// signal for passing errors / messages down to httpserver
type AppSignal struct {
	Message string
	Error   error
}
