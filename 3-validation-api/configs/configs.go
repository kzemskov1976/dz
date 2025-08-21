package configs

type Config struct {
	Email    string
	Password string
	Address  string
}

func LoadConfig() *Config {
	return &Config{
		Email:    "kirill.zemskoff@gmail.com",
		Password: "123",
		Address:  "smtp.gmail.com",
	}
}
