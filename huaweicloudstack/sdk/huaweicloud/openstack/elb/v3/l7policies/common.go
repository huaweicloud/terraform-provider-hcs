package l7policies

type RedirectUrlConfig struct {
	StatusCode string `json:"status_code,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Host       string `json:"host,omitempty"`
	Port       string `json:"port,omitempty"`
	Path       string `json:"path,omitempty"`
	Query      string `json:"query,omitempty"`
}
