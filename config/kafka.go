package config

type Producer struct {
	URL       string  `yaml:"url"`        // kafka broker url
	ClientID  string  `yaml:"client_id"`  // produce client id
	Acks      string  `yaml:"acks"`       // 0, 1, all
	Topic     string  `yaml:"topic"`      // produce topic
	BatchTime float64 `yaml:"batch_time"` // batch time
}
