package options

// Opts represent struct for all program options
type Opts struct {
	Cache   Cache  `group:"cache" namespace:"cache" env-namespace:"CACHE"`
	Log     Log    `group:"log" namespace:"log" env-namespace:"LOG"`
	APIPort string `long:"api-port" env:"API_PORT" default:"8080" description:"what port listen"`
	Swagger bool   `long:"swagger" env:"SWAGGER" description:"host swagger docs"`
	Schema  string `long:"SCHEMA" env:"SCHEMA" description:"providers specs"`
	MaxSide int    `long:"MAX_SIDE" env:"MAX_SIDE" default:"10" description:"max square side"`
}

// Cache represent struct for Cache options
type Cache struct {
	Enable bool   `long:"enable" env:"ENABLE" description:"enable cache"`
	Path   string `long:"path" env:"PATH" default:"./data/cache" description:"a path for cache dir"`
	Alive  int    `long:"alive" env:"ALIVE" default:"14400" description:"cache alive in minutes"`
}

// Log represent struct for Log options
type Log struct {
	Save       bool   `long:"save" env:"SAVE" description:"enable logs save"`
	Path       string `long:"path" env:"PATH" default:"./data/logs/log.jsonl" description:"a path for logs dir"`
	MaxBackups int    `long:"max-backups" env:"MAX_BACKUPS" default:"3" description:"max backups"`
	MaxSize    int    `long:"max-size" env:"MAX_SIZE" default:"1" description:"max logs size in megabytes"`
	MaxAge     int    `long:"max-age" env:"MAX_AGE" default:"7" description:"max logs age"`
}
