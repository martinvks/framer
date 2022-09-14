package types

type Header struct {
	Name  string
	Value string
}

type Headers []Header

func (hs Headers) Get(name string) (value string, ok bool) {
	for _, header := range hs {
		if header.Name == name {
			return header.Value, true
		}
	}

	return "", false
}
