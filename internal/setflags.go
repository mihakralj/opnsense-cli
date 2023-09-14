package internal

var (
	verbose    int
	force      bool
	host       string
	configfile string
	nocolor    bool
	depth      int
	xmlFlag    bool
	yamlFlag   bool
	jsonFlag   bool
)

func SetFlags(v int, f bool, h string, config string, nc bool, dpt int, x bool, y bool, j bool) {
	if v < 1 || v > 5 {
		Log(1, "invalid verbosity level %d. It should be between 1 and 5", v)
	}
	verbose = v
	force = f
	host = h
	configfile = config
	nocolor = nc
	depth = dpt
	xmlFlag = x
	yamlFlag = y
	jsonFlag = j
	Log(5, "flags:\tverbose=%d, host=%s, config=%s", verbose, host, configfile)
	if nc {
		for key := range c {
			delete(c, key)
		}
		c["ell"] = "..."
		c["arw"] = " -> "
	}
}
