// Package logging provides structured logging with logrus.
package logging

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/natefinch/lumberjack"
	"io"
	"os"
	"strings"
	"time"
)



var (
	// Logger is a configured logrus.Logger.
	Logger *logrus.Logger
)

//Log formatting struct
type LogFormat struct {
	TimestampFormat string
}

//Customt formatting
func (f *LogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}


	appname:="doubtnut"
	b.WriteString("producerType==\"service\" ")
	b.WriteString("producerName==\"" + appname + "\" ")
	b.WriteString("time==\"")
	//2018-11-14T14:00:05Z
	//stdISO8601TZ, stdISO8601ColonTZ, stdISO8601SecondsTZ, stdISO8601ShortTZ, stdISO8601ColonSecondsTZ, stdNumTZ, stdNumColonTZ, stdNumSecondsTz, stdNumShortTZ, stdNumColonSecondsTZ
	b.WriteString(entry.Time.UTC().Format(time.RFC3339) + "\"")

	b.WriteString(" ")
	b.WriteString("level==")
	fmt.Fprint(b, "\""+strings.ToUpper(entry.Level.String())+"\"")

	b.WriteString(" ")
	b.WriteString("message==")
	if entry.Message != "" {
		b.WriteString(fmt.Sprintf("\"%s\"", entry.Message))
	}
	b.WriteString(" ")

	//mdc
	for key, value := range entry.Data {
		b.WriteString(" " + key)
		b.WriteString("==\"")
		fmt.Fprint(b, value)
		b.WriteString("\"")
	}

	//end mdc


	b.WriteByte('\n')
	return b.Bytes(), nil
}

// NewLogger creates and configures a new logrus Logger.
func NewLogger() *logrus.Entry {

	dir := "opt/logs"
	fileName := "info"
	fileName+= ".log"
	path := dir + "/" + fileName

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err = os.Create(fileName)
		if err != nil {
			panic(err)
		}
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    1,
		MaxBackups: 20,
		MaxAge:     28,
		LocalTime:  true,
	}
	logrus.SetOutput(lumberjackLogger)
	Logger = logrus.New()
	Logger.Formatter = &logrus.JSONFormatter{}
	Logger.Level = logrus.DebugLevel
	Logger.Out = lumberjackLogger

	mw := io.MultiWriter(os.Stdout, lumberjackLogger)
	Logger.Out = mw

	return Logger.WithFields(logrus.Fields{"producerType": "service", "producerName": "LAS"})


}

