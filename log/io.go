package log

import log "github.com/sirupsen/logrus"

type Writer struct{}

func (w *Writer) Write(p []byte) (n int, err error) {
	log.Trace(string(p))
	return
}
