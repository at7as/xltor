package main

// Config struct which contains
type Config struct {
	File  string        `json:"file"`
	Sheet []ConfigSheet `json:"sheet"`
}

// ConfigSheet struct which contains
type ConfigSheet struct {
	Name string      `json:"name"`
	Ref  string      `json:"ref"`
	Row  []ConfigDim `json:"r"`
	Col  []ConfigDim `json:"c"`
}

// ConfigDim struct which contains
type ConfigDim struct {
	Index  bool         `json:"index"`
	Title  string       `json:"title"`
	Length []int        `json:"length"`
	Root   int          `json:"root"`
	Fixed  bool         `json:"fixed"`
	Count  int          `json:"count"`
	Hier   []ConfigHier `json:"hier"`
}

// ConfigHier struct which contains
type ConfigHier struct {
	Fixed  bool    `json:"fixed"`
	Count  int     `json:"count"`
	Random []int   `json:"random"`
	Ratio  float64 `json:"ratio"`
}

type sheet struct {
	cfgSheet *ConfigSheet
	row      []dim
	col      []dim
}

type dim struct {
	cfgDim *ConfigDim
	root   []dimItem
}

type dimItem struct {
	title    string
	children []dimItem
}
