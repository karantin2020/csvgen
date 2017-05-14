package test

type Foo struct {
	a string
	b int64
	c int
	s *Zoo
}

type Boo struct {
	a  *float64
	b  *string
	z  Zoo
	d  int16
	f  float32
	ff float64
	u  uint8
}

// type Boo struct {
// 	a *floatas64
// 	b strdsaing
// 	z Zoo
// }

type Zoo string

func (this *Zoo) UnmarshallCSV(in string) error {
	*this = Zoo(in)
	return nil
}

func (this *Zoo) MarshallCSV() (string, error) {
	return string(*this), nil
}
