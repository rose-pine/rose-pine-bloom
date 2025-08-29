package config

type BuildConfig struct {
	Template string
	Output   string
	Prefix   string
	Format   string
	Plain    bool
	Commas   bool
	Spaces   bool
}

type BuildTemplateConfig struct {
	Input   string
	Output  string
	Variant string
	Prefix  string
	Format  string
	Plain   bool
	Commas  bool
	Spaces  bool
}
