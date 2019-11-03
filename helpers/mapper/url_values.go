package mapper

import "net/url"

// MarshalMapToURLValues creates a url.Values from a map[string]string
func MarshalMapToURLValues(m map[string]string) url.Values {
	u := make(url.Values)

	for key, value := range m {
		u.Add(key, value)
	}

	return u
}
