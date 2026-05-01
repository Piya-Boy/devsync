package transport

import (
	"fmt"
	"io"
	"os"
)

// prefixWriter wraps an io.Writer and prepends a prefix to every write.
type prefixWriter struct {
	prefix string
	w      io.Writer
}

func newPrefixWriter(prefix string, w io.Writer) *prefixWriter {
	if w == nil {
		w = os.Stdout
	}
	return &prefixWriter{prefix: prefix, w: w}
}

func (p *prefixWriter) Write(b []byte) (int, error) {
	_, err := fmt.Fprintf(p.w, "%s%s", p.prefix, b)
	return len(b), err
}
