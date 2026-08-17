package model

type Item struct {
	ID          int    `toml:"id"`
	Title       string `toml:"title"`
	Description string `toml:"description"`
	State       string `toml:"state"`
	Created     string `toml:"created"`
	Updated     string `toml:"updated"`
}