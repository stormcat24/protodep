package version

import "fmt"

var (
	gitCommit     = "unspecified"
	gitCommitFull = "unspecified"
	buildDate     = "unspecified"
	version       = "unspecified"
)

type Info struct {
	GitCommit     string
	GitCommitFull string
	BuildDate     string
	Version       string
}

func Get() Info {
	return Info{
		GitCommit:     gitCommit,
		GitCommitFull: gitCommitFull,
		BuildDate:     buildDate,
		Version:       version,
	}
}

func (i Info) String() string {
	return fmt.Sprintf(
		`{"Version": "%s", "GitCommit": "%s", "GitCommitFull": "%s", "BuildDate": "%s"}`,
		i.Version,
		i.GitCommit,
		i.GitCommitFull,
		i.BuildDate,
	)
}
