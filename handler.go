package prettyslog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

type PrettyHandler struct {
	opts           slog.HandlerOptions
	preformatted   []byte   // data from WithGroup and WithAttrs
	unopenedGroups []string // groups from WithGroup that haven't been opened
	openedCount    int
	mu             *sync.Mutex
	out            io.Writer
}

// NewHandler creates a [PrettyHandler] that writes to w, using the given options. If opts is nil, the default options are used.
func NewHandler(out io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	h := &PrettyHandler{out: out, mu: new(sync.Mutex)}
	if opts != nil {
		h.opts = *opts
	}

	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	bufp := allocBuf()
	buf := *bufp
	defer func() {
		*bufp = buf
		freeBuf(bufp)
	}()

	rep := h.opts.ReplaceAttr

	// Time
	if !r.Time.IsZero() {
		a := slog.Time(slog.TimeKey, r.Time)
		if rep != nil {
			a = rep(nil, a)
			if a.Equal(slog.Attr{}) {
				a = slog.Attr{}
			}
		}
		buf = h.appendAttr(buf, a, 0)
	}

	// Level
	a := slog.Any(slog.LevelKey, levels[r.Level])
	if rep != nil {
		replaced := rep(nil, slog.Any(slog.LevelKey, r.Level))
		if replaced.Equal(slog.Attr{}) {
			a = slog.Attr{}
		} else if _, ok := levels[r.Level]; !ok || replaced.Value.Kind() != slog.KindAny {
			a = replaced
		}
	}
	buf = h.appendAttr(buf, a, 0)

	// Source
	if h.opts.AddSource {
		src := r.Source()
		if src == nil {
			src = &slog.Source{}
		}
		//Optimize to minimize allocation.
		srcbufp := allocBuf()
		defer freeBuf(srcbufp)
		*srcbufp = fmt.Appendf(*srcbufp, "%s:%d", src.File, src.Line)
		buf = h.appendAttr(buf, slog.String(slog.SourceKey, filepath.Base(string(*srcbufp))), 0)
	}

	// Message
	buf = h.appendAttr(buf, slog.String(slog.MessageKey, r.Message), 0)

	// Insert preformatted attributes just after built-in ones.
	buf = append(buf, h.preformatted...)
	if r.NumAttrs() > 0 {
		// Magenta opening group text
		buf = h.appendUnopenedGroups(buf, "\x1B[35m")
		r.Attrs(func(a slog.Attr) bool {
			buf = h.appendAttr(buf, a, len(h.unopenedGroups))
			return true
		})
		// Magenta closing group text
		buf = h.closeGroups(buf, len(h.unopenedGroups), "\x1B[35m)\x1B[0m")
	}
	// Blue closing group text
	buf = h.closeGroups(buf, h.openedCount, "\x1B[34m)\x1B[0m")
	// End of line
	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

func (h *PrettyHandler) appendAttr(buf []byte, a slog.Attr, openedCount int) []byte {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return buf
	}

	switch a.Value.Kind() {
	case slog.KindTime:
		buf = a.Value.Time().AppendFormat(buf, time.DateTime)
	case slog.KindGroup:
		attrs := a.Value.Group()
		if len(attrs) == 0 {
			return buf
		}
		if a.Key != "" {
			// Cyan opening group text
			buf = fmt.Appendf(buf, "\x1B[36m%s(\x1B[0m", a.Key)
			openedCount++
		}
		for _, ga := range attrs {
			buf = h.appendAttr(buf, ga, openedCount)
		}
		// Cyan closing group text
		buf = h.closeGroups(buf, 1, "\x1B[36m)\x1B[0m")
	default:
		if a.Key == slog.SourceKey || a.Key == slog.MessageKey || a.Key == slog.LevelKey || a.Key == slog.TimeKey {
			buf = append(buf, a.Value.String()...)
		} else {
			// White key text
			buf = fmt.Appendf(buf, "\x1B[1;37m%s=\x1B[0m\"%s\"", a.Key, a.Value.String())
		}
	}
	buf = append(buf, ' ')

	return buf
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := *h
	// Add an unopened group to h2 without modifying h.
	h2.unopenedGroups = make([]string, len(h.unopenedGroups)+1)
	copy(h2.unopenedGroups, h.unopenedGroups)
	h2.unopenedGroups[len(h2.unopenedGroups)-1] = name
	return &h2
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := *h
	// Force an append to copy the underlying array.
	pre := slices.Clip(h.preformatted)
	// Add all groups from WithGroup that haven't already been added.
	// Blue opening group text
	h2.preformatted = h2.appendUnopenedGroups(pre, "\x1B[34m")
	// Each of those groups increased the count by 1.
	h2.openedCount += len(h2.unopenedGroups)
	// Now all groups have been opened.
	h2.unopenedGroups = nil
	// Pre-format the attributes.
	for _, a := range attrs {
		h2.preformatted = h2.appendAttr(h2.preformatted, a, h2.openedCount)
	}
	return &h2
}

func (h *PrettyHandler) appendUnopenedGroups(buf []byte, color string) []byte {
	for _, g := range h.unopenedGroups {
		buf = fmt.Appendf(buf, "%s%s(\x1B[0m", color, g)
	}
	return buf
}

func (h *PrettyHandler) closeGroups(buf []byte, n int, close string) []byte {
	// Trim trailing space
	if len(buf) > 0 && buf[len(buf)-1] == ' ' {
		buf = buf[:len(buf)-1]
	}
	for i := 0; i < n; i++ {
		buf = append(buf, close...)
	}
	return buf
}
