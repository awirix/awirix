package context

import (
	"context"
	"time"
)

// Context is a wrapper around the standard context.Context interface.
// It is used to provide a cancelable context.Context to the Lua VM.
type Context struct {
	done     chan struct{}
	err      error
	values   map[any]any
	deadline time.Time
}

// New returns a new Context.
func New() *Context {
	return &Context{
		done:   make(chan struct{}),
		values: make(map[any]any),
	}
}

func (c *Context) Set(key, value any) {
	c.values[key] = value
}

// Cancel cancels the context.
func (c *Context) Cancel() {
	c.err = context.Canceled
	close(c.done)
}

// Reset resets the context.
// This is useful for reusing the same context after it has been canceled.
func (c *Context) Reset() {
	c.err = nil
	c.done = make(chan struct{})
}

func (c *Context) SetDeadline(t time.Time) {
	c.deadline = t
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (c *Context) Done() <-chan struct{} {
	return c.done
}

func (c *Context) Err() error {
	return c.err
}

func (c *Context) Value(key any) any {
	return c.values[key]
}