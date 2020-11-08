package lualib

type ContentType map[string][]string

var HttpType ContentType = make(ContentType)

func InitType() {
	HttpType["text/plain"] = []string{"txt"}
	HttpType["text/css"] = []string{"css"}
	HttpType["text/html"] = []string{"lua", "htm", "html", "stm"}
	HttpType["image/gif"] = []string{"gif"}
	HttpType["image/png"] = []string{"png"}
	HttpType["image/jpeg"] = []string{"jpe", "jpeg", "jpg"}
	HttpType["audio/mpeg"] = []string{"mp3"}
	HttpType["application/zip"] = []string{"zip"}
	HttpType["application/envoy"] = []string{"evy"}
	HttpType["application/fractals"] = []string{"fif"}
	HttpType["application/x-javascript"] = []string{"js"}
	HttpType["application/octet-stream"] = []string{"*", "exe", "bin", "class", "apk"}
}

func (this *ContentType) Type(name string) string {
	for t, arr := range *this {
		for _, n := range arr {
			if name == n {
				return t
			}
		}
	}
	return "application/octet-stream"
}

func (this *ContentType) Name(typ string) []string {
	for t, name := range *this {
		if typ == t {
			return name
		}
	}
	return []string{}
}

func (this *ContentType) Is(name string, typ string) bool {
	if this.Type(name) != typ {
		return false
	}
	return true
}