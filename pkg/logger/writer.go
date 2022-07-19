package logger

import (
	"io"
	"os"
)

// Writer wraps an io.Writer with a mutex, so that multiple loggers can be
// created that use the same writer.
type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// Write writes in the io.Writer. Safe for concurrent use.
func (w *Writer) Write(p []byte) (int, error) {
	if w.writer == nil {
		return 0, nil
	}
	return w.writer.Write(p)
}

// Sync flushes the writer's buffer to the fs
func (w *Writer) Sync() error {
	file, ok := w.writer.(*os.File)
	if ok {
		return file.Sync()
	}
	return nil
}

// Close closes the writer and removes the file lock
func (w *Writer) Close() error {
	file, ok := w.writer.(*os.File)
	if ok {
		return file.Close()
	}
	return nil
}

func (w *Writer) SetWriter(writer io.Writer) {
	w.writer = writer
}
