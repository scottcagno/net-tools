package web

import "html/template"

type Model interface {
	ID() int
}

type Service interface {
	Render(template template.Template)
	FindAll() []Model
	FindOne(id int) Model
	Save(model Model) Model
	Delete(model Model)
}

type Repository interface {
}
