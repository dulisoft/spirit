package annotation

type App struct {
	Service     string
	Port        string
	Controllers []Controller
}

func (a *App) RegisterRouters() error {
	return WirteTemplate(a, getRouterTemplatePath(), getMainCodePath())
}

func NewApp(controllers []Controller) *App {
	return &App{
		Service:     APP_NAME,
		Port:        PORT,
		Controllers: controllers,
	}
}
