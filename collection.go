package errors

import "sync"

// Collection manages non nil errors with threadsafe methods.
type Collection struct {
	errs []error
	mtx  sync.RWMutex
}

// NewCollection returns a new Collection with initialErrors.
func NewCollection(initialErrors ...error) *Collection {
	return &Collection{errs: flatten(initialErrors)}
}

// Add an error, but only if it is not nil.
func (c *Collection) Add(err error) {
	if err == nil {
		return
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.errs = append(c.errs, err)
}

// Remove removes all instances of err from the collection.
func (c *Collection) Remove(err error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for i := len(c.errs) - 1; i >= 0; i-- {
		if c.errs[i] == err {
			copy(c.errs[i:], c.errs[i+1:])
			c.errs = c.errs[:len(c.errs)-1]
		}
	}
}

// Errors returns the errors in the collection.
func (c *Collection) Errors() []error {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.errs
}

// Combine returns the errors as a single error
// combination, similar to the Combine function.
// If the Collection is empty, then nil is returned.
func (c *Collection) Combine() error {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if len(c.errs) == 0 {
		return nil
	}
	return &combination{
		c.errs,
		callers(0),
	}
}
