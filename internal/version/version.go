package version

var (
	Version = "v0.0.0-dev"
	Commit  = "HEAD"
	Date    = "unknown"
	BuiltBy = "local"
)

type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	BuiltBy string `json:"builtBy"`
}

func Get() Info {
	return Info{
		Version: Version,
		Commit:  Commit,
		Date:    Date,
		BuiltBy: BuiltBy,
	}
}
