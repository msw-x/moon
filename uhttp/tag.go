package uhttp

const (
	UrlTag    string = "url"
	HeaderTag string = "http"
)

func OmitEmpty(flags []string) bool {
	for _, flag := range flags {
		if flag == "omitempty" {
			return true
		}
	}
	return false
}
