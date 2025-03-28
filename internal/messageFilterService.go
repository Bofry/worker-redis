package internal

type MessageFilter map[string]MessageFilterAffinity

func (f MessageFilter) Add(stream string, filter MessageFilterAffinity) {
	f[stream] = filter
}

func (f MessageFilter) Match(m *Message) bool {
	if m == nil {
		return false
	}
	if len(f) == 0 {
		return true
	}

	var proc = f[m.Stream]
	if proc != nil {
		return proc.Filter(m)
	}
	return true
}
