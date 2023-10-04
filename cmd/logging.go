/*
Copyright © 2023 Jean-Marc Meessen jean-marc@meessen-web.org

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type LOGGER struct {
	debug *log.Logger
	prod  *log.Logger
}

var lock = &sync.Mutex{}
var loggers *LOGGER

func GetLoggerInstance() *LOGGER {
	lock.Lock()
	defer lock.Unlock()

	if loggers == nil {
		loggers = &LOGGER{}
	}

	return loggers
}
func initLoggers() {
	loggers := GetLoggerInstance()
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("debug log file not created", err.Error())
	}
	loggers.debug = log.New(f, "[DEBUG]", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
	loggers.prod = log.New(os.Stderr, "[log]", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
}
