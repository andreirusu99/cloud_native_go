package log

import (
	"bufio"
	"cloud_native_go/pkg/misc"
	"fmt"
	"os"
)

var outFormatString = "%d\t%d\t%s\t%s\n" // index - type - key - value
var inFormatString = "%d\t%d\t%s\t%s" // index - type - key - value

type FileTransactionLogger struct {
	events 		chan<- misc.Event
	errors 		<-chan error
	lastIndex 	uint64
	file 		*os.File
}

// constructor
func NewFileTransactionLogger(path string) (TransactionLogger, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		return nil, fmt.Errorf("transaction log file could not be opened: %w", err)
	}

	return &FileTransactionLogger{file: file}, nil
}

// put events on a channel to be processed by a goroutine
func (l *FileTransactionLogger) LogPut(key, value string) {
	l.events <- misc.Event{Type: misc.EventPut, Key: key, Value: value}
}

func (l *FileTransactionLogger) LogDelete(key string) {
	l.events <- misc.Event{Type: misc.EventDelete, Key: key}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) Run() {
	// create the event and error channels
	events := make(chan misc.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	// goroutine to consume the events channel
	go func () {
		for event := range events {

			// increment index
			l.lastIndex++

			// write event to file
			_, err := fmt.Fprintf(
				l.file,
				outFormatString,
				l.lastIndex, event.Type, event.Key, event.Value,
			)

			if err != nil {
				errors <- err
				return
			}
		}
	}()
}


func (l *FileTransactionLogger) ReplayEvents() (<-chan misc.Event, <-chan error){
	scanner := bufio.NewScanner(l.file)
	events := make(chan misc.Event)
	errors := make(chan error, 1)

	go func () {

		var event misc.Event
		defer func() {
			close(events)
			close(errors)
		}()

		for scanner.Scan() {
			line := scanner.Text()

			_, err := fmt.Sscanf(line, inFormatString, &event.Index, &event.Type, &event.Key, &event.Value)
			if err != nil {
				errors <- fmt.Errorf("log parse error: %w, line: %s", err, line)
				return
			}

			if event.Index <= l.lastIndex {
				errors <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			l.lastIndex = event.Index // update the highest index
			events <- event // push (stream) the event on the channel
		}

		if err := scanner.Err(); err != nil {
			errors <- fmt.Errorf("failed to read transaction log: %w", err)
			return
		}

		fmt.Printf("--> Read %d lines from log\n", l.lastIndex)
	}()

	return events, errors

}