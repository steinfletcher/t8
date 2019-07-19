package domain

type TemplateRenderer interface {
	Render(targetDir string, conf Config)
}
