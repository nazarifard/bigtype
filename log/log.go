package log

import (
	"os"
	"slices"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogInterface interface {
	Debug(msg string)
	Info(msg string)
	// Warn(msg string)
	// Error(msg string)
	// Fatal(msg string)
}

var Logger LogInterface

func SetLogger(lg LogInterface) { Logger = lg }

// type nolog struct{}
// func (l nolog) Debug(msg string) {}
// func (l nolog) Info(msg string)  {}

var VerboseMode bool

type ZeroLog struct{}

func (l ZeroLog) Debug(msg string) {
	if VerboseMode {
		log.Debug().Msg(msg)
	}
}

func (l ZeroLog) Info(msg string) {
	if VerboseMode {
		log.Info().Msg(msg)
	}
}

func init() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z07:00" //RFC3339Milli
	SetLogger(ZeroLog{})

	VerboseMode = slices.Contains(os.Args, "-v") ||
		slices.Contains(os.Args, "-test.v=true") //go test option
	if VerboseMode {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	//fmt.Println("\n", VerboseMode, "  ", os.Args, "aaaaaaaaaaaaa")
}
