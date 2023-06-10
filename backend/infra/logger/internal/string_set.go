package internal

type StringSet struct {
	data map[string]bool
}

func NewStringSet() StringSet {
	return StringSet{data: map[string]bool{}}
}

func (s *StringSet) Add(e ...string) {
	for _, elem := range e {
		s.data[elem] = true
	}
}

func (s *StringSet) Contains(e string) bool {
	_, ok := s.data[e]
	return ok
}

func (s *StringSet) Delete(e ...string) {
	for _, elem := range e {
		delete(s.data, elem)
	}
}

func (s *StringSet) Size() int {
	return len(s.data)
}

func (s *StringSet) Union(other StringSet) StringSet {
	res := NewStringSet()
	for k, ok := range s.data {
		if ok {
			res.Add(k)
		}
	}
	for k, ok := range other.data {
		if ok {
			res.Add(k)
		}
	}
	return res
}

func (s *StringSet) Intersection(other StringSet) StringSet {
	res := NewStringSet()
	for k, ok := range s.data {
		if ok && other.Contains(k) {
			res.Add(k)
		}
	}
	return res
}

func (s *StringSet) Difference(other StringSet) StringSet {
	res := NewStringSet()
	for k, ok := range s.data {
		if ok && !other.Contains(k) {
			res.Add(k)
		}
	}
	return res
}

func (s *StringSet) Equal(other StringSet) bool {
	if s.Size() != other.Size() {
		return false
	}
	for k, ok := range s.data {
		if ok && !other.Contains(k) {
			return false
		}
	}
	return true
}
