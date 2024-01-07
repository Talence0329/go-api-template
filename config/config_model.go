package config

import "time"

type Env struct {
	SystemENV
	JWTENV
	CookiesENV
	FunctionalENV
	ServiceENV
}

// SystemENV : 系統環境變數
type SystemENV struct {
	Project      string        `env:"project" validate:"required"`
	Environment  string        `env:"environment" validate:"required"`
	Port         int           `env:"port" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"readTimeout" envDefault:"1m0s"`
	WriteTimeout time.Duration `env:"writeTimeout" envDefault:"1m0s"`
	IsTwoStep    bool          `env:"isTwostep" envDefault:"true"`
	IsUseLogAll  bool          `env:"isUseLogAll" envDefault:"false"`
}

type FunctionalENV struct {
	AllowOrigins       string `env:"allowOrigins" envDefault:""`
	AllowCredentials   bool   `env:"allowCredentials" envDefault:"false"`
	LogLevel           string `env:"logLevel" envDefault:"error"`
	PasswordHash       string `env:"passwordHash" validate:"required"`
	GoogleClientID     string `env:"googleClientId" validate:"required"`
	GoogleClientSecret string `env:"googleClientSecret" validate:"required"`
	FBClientID         string `env:"fbClientId" validate:"required"`
	FBClientSecret     string `env:"fbClientSecret" validate:"required"`
	TerminalServer     string `env:"terminalServer" envDefault:"http://211.149.198.94:8080" validate:"required"`
	TerminalUsername   string `env:"terminalUsername" envDefault:"benz" validate:"required"`
	TerminalPassword   string `env:"terminalPassword" envDefault:"admin123" validate:"required"`
}

type JWTENV struct {
	JwtSecretkey        string `env:"jwtSecretkey" envDefault:"33EFA739GGA2443VC2646294A404D6351"`
	JwtTwoStepSecretkey string `env:"jwtTwoStepSecretkey" envDefault:"81QFF527GUE1421GC25766794A405D651"`
	JwtExpires          int    `env:"jwtExpires" envDefault:"7200"`
	JwtTokenKey         string `env:"jwtTokenKey" envDefault:"jwtToken"`
	JwtUserKey          string `env:"jwtUserKey" envDefault:"currentUser"`
}

type CookiesENV struct {
	Secure bool `env:"secure" envDefault:"true"`
	MaxAge int  `env:"maxAge" envDefault:"7200" validate:"min=1"`
}

type ServiceENV struct {
	REDIS_URL string `envDefault:"http://${REDIS_URL}" envExpand:"true"`
}
