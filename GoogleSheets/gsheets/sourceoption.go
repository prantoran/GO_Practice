package gsheets

// SourceOption is the options interface which contains the method Apply()
type SourceOption interface {
	Apply(*Source)
}

// SourceOptionFunc is the function
type SourceOptionFunc func(*Source)

// Apply is the function
func (f SourceOptionFunc) Apply(s *Source) {
	f(s)
}

// WithHeader is the function
func WithHeader() SourceOption {
	return SourceOptionFunc(func(s *Source) {
		s.header = true
	})
}
