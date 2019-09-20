package model

type Command struct {
	Id         string   `json:"id"`
	Enabled    bool     `json:"enabled"`
	Workers    int      `json:"workers"`
	KeepAlive  bool     `json:"keep_alive"`
	RetryDelay Duration `json:"retry_delay"`
	WorkDir    string   `json:"work_dir"`
	Exec       string   `json:"exec"`
	Args       []string `json:"args"`
}
