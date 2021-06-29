package text

type Lookuper interface {
	Lookup(string) (string, bool)
}
