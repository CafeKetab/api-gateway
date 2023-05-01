package http

type Config struct {
	ListenPort int `koanf:"listen_port"`
	TargetUrls struct {
		Users string `koanf:"users"`
		Books string `koanf:"books"`
	}
}
