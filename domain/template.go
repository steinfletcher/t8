package domain

//go:generate mockgen -source=template.go -destination=../mocks/template.go -package=mocks

type TemplateRenderer interface {
	Render(targetDir string, conf Config)
}
