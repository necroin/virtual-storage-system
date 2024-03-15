package config

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"runtime"
	"time"
	"vss/src/lan"
	"vss/src/settings"
	"vss/src/utils"
)

type StorageRole struct {
	Enable bool `yaml:"enable"`
}

type RunnerRole struct {
	Enable   bool   `yaml:"enable"`
	Platform string `yaml:"platform"`
}

type RouterRole struct {
	Enable bool `yaml:"enable"`
}

type Roles struct {
	Storage StorageRole `yaml:"storage"`
	Runner  RunnerRole  `yaml:"runner"`
	Router  RouterRole  `yaml:"router"`
}

type Log struct {
	Enable bool   `yaml:"enable"`
	Path   string `yaml:"path"`
	Level  string `yaml:"level"`
}

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
}

type Config struct {
	Url         string `yaml:"url"`
	ListenPort  string `yaml:"listen_port"`
	Roles       Roles  `yaml:"roles"`
	Log         Log    `yaml:"log"`
	User        User   `yaml:"user"`
	Certificate tls.Certificate
	RootCAs     *x509.CertPool
}

func Load() (*Config, error) {
	url := flag.String("url", lan.GetMyLanAddr(), "server url")
	listenPort := flag.String("listen-port", settings.DefaultListenPort, "server topology listen port")

	storageRoleEnable := flag.Bool("storage", true, "enables 'storage' role")
	runnerRoleEnable := flag.Bool("runner", false, "enables 'runner' role")
	routerRoleEnable := flag.Bool("router", false, "enables 'router' role")
	platform := flag.String("platform", runtime.GOOS, "OS platform ('windows', 'linux', 'darwin', etc.)")

	logEnable := flag.Bool("log-enable", false, "enbales logs")
	logPath := flag.String("log-path", "logs/log_"+time.Now().Format("2006-01-02T15:04:05")+".txt", "path to logs file")
	logLevel := flag.String("log-level", "info", "logs level (error, info, verbose, debug)")

	username := flag.String("username", settings.DefaultUsername, "authentication username")
	password := flag.String("password", settings.DefaultPassword, "authentication password")
	token := flag.String("token", utils.GenerateSecureToken(10), "security token")

	certificateCrt := flag.String("cert-crt", "certificates/vss.crt", "path to vss.crt")
	certificateKey := flag.String("cert-key", "certificates/vss.key", "path to vss.key")

	flag.Parse()

	certificate, err := tls.LoadX509KeyPair(*certificateCrt, *certificateKey)
	if err != nil {
		return nil, fmt.Errorf("[Config] failed load certificate: %s", err)
	}

	rootCAs := x509.NewCertPool()
	certificateAuthority, err := x509.ParseCertificate(certificate.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("[Config] failed parse certificate: %s", err)
	}
	rootCAs.AddCert(certificateAuthority)

	config := &Config{
		Url:        *url,
		ListenPort: *listenPort,
		Roles: Roles{
			Storage: StorageRole{
				Enable: *storageRoleEnable,
			},
			Runner: RunnerRole{
				Enable:   *runnerRoleEnable,
				Platform: *platform,
			},
			Router: RouterRole{
				Enable: *routerRoleEnable,
			},
		},
		Log: Log{
			Enable: *logEnable,
			Path:   *logPath,
			Level:  *logLevel,
		},
		User: User{
			Username: *username,
			Password: *password,
			Token:    *token,
		},
		Certificate: certificate,
		RootCAs:     rootCAs,
	}

	return config, nil
}
