package l7policies

type RedirectUrlConfig struct {
	StatusCode string `json:"status_code"`
	Protocol   string `json:"protocol"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Path       string `json:"path"`
	Query      string `json:"query"`
}
