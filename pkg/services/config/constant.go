package config

const (
	Version  = "0.1.0"
	Name     = "dpx"
	FileName = Name + ".yml"
	PathEnv  = "PATH"
	PathData = "PATH=%s:$PATH"
	BaseDir  = ".dpx"
	MaxLevel = 10

	// alias file
	AliasFile     = Name + "-alias"
	AliasData     = `dpx exec -c $(dirname ${0%/*/*})/` + FileName + ` ${0##*/} "$@"`
	AliasFileMode = 0744

	// bin file (symlink)
	BinDir      = BaseDir + "/bin"
	BinFileMode = AliasFileMode

	// others
	ProcessIdLength = 12
	IndentLevel     = 2
)
