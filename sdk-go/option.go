package buffer_pool

type Option func(*MultiValueSet)

func WithSizeLimit(size uint64) Option {
	return func(m *MultiValueSet) {
		m.size = size
	}
}
