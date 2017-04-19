package runtime

import "strings"

// EqualMatcher performs a case-sensitive equality match for request metadata keys
func EqualMatcher(s string) ServeMuxOption {
	return func(mux *ServeMux) {
		mux.headerMatchers = append(mux.headerMatchers, func(h string) bool {
			return h == s
		})
	}
}

// EqualFoldMatcher performs a case-insensitive equality match for request metadata keys
func EqualFoldMatcher(s string) ServeMuxOption {
	return func(mux *ServeMux) {
		mux.headerMatchers = append(mux.headerMatchers, func(h string) bool {
			return strings.EqualFold(h, s)
		})
	}
}

// PrefixMatcher performs a case-sensitive prefix match for request metadata keys
func PrefixMatcher(prefix string) ServeMuxOption {
	return func(mux *ServeMux) {
		mux.headerMatchers = append(mux.headerMatchers, func(h string) bool {
			return strings.HasPrefix(h, prefix)
		})
	}
}

// PrefixFoldMatcher performs a case-insensitive prefix match for request metadata keys
func PrefixFoldMatcher(prefix string) ServeMuxOption {
	return func(mux *ServeMux) {
		mux.headerMatchers = append(mux.headerMatchers, func(h string) bool {
			return len(h) >= len(prefix) && strings.EqualFold(prefix, h[:len(prefix)])
		})
	}
}
