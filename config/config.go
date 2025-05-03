package config

type ApiConfig struct {
	Port       string `json:"port"`
	MongoDBUri string `json:"mongodb_uri"`
}
