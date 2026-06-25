package tag

type ID string

func NewID(id string) (ID, error) {
	return ID(id), nil
}

func (id ID) String() string {
	return string(id)
}
