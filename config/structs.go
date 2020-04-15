package config

type (
	Config struct {
		Bot Bot
	}

	Bot struct {
		Token  string
		Prefix string
	}
)
