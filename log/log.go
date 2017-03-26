package log

import (
	"fmt"
	"github.com/op/go-logging"
	"os"
	"path"
	"runtime"
	"strings"
)

var log = logging.MustGetLogger("kijiji-rentals")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} ▶ %{level} %{id:03x}%{color:reset} %{message}",
)

//"%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
func init() {

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	syslog_backend, err := logging.NewSyslogBackend("kijiji-rentals")
	if err != nil {
		fmt.Println("Failed to initialize syslog logger ", err.Error())
		os.Exit(1)
	}

	backend_formatted := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backend_formatted, syslog_backend)

}

func trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Sprintf("[%s:%d][%s] ", path.Base(file), line, path.Base(f.Name()))
}

func joiner(data []interface{}) string {

	args_str := make([]string, len(data))
	for i := range data {
		args_str[i], _ = data[i].(string)
	}

	return strings.Join(args_str, " ")

}

func Criticalf(data ...interface{}) {
	log_message := joiner(data)
	trace_info := trace()
	log.Critical(strings.Join([]string{trace_info, log_message}, ""))
	os.Exit(1)

}

func Error(data ...interface{}) {
	log_message := joiner(data)
	trace_info := trace()
	log.Error(strings.Join([]string{trace_info, log_message}, ""))
}

func Warning(data ...interface{}) {
	log_message := joiner(data)
	trace_info := trace()
	log.Warning(strings.Join([]string{trace_info, log_message}, ""))

}

func Info(data ...interface{}) {
	log_message := joiner(data)
	trace_info := trace()
	log.Info(strings.Join([]string{trace_info, log_message}, ""))

}

func Debug(data ...interface{}) {
	log_message := joiner(data)
	trace_info := trace()
	log.Debug(strings.Join([]string{trace_info, log_message}, ""))

}
