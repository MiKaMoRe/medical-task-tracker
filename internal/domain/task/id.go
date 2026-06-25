package task

type ID int

func (id ID) Int() int {
	return int(id)
}

func IDFromInt(v int) ID {
	return ID(v)
}
