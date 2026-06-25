package tag

type Name string

func NewName(name string) (Name, error) {
	return Name(name), nil
}

func (n Name) String() string {
	return string(n)
}
