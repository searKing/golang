package class

type Classer interface {
	GetDerived() Classer
	SetDerived(derived Classer)
}

type Class struct {
	derived Classer
}

func NewClass() *Class {
	return &Class{}
}

func (task *Class) GetDerived() Classer {
	return task.derived
}

func (task *Class) SetDerived(derived Classer) {
	task.derived = derived
}
