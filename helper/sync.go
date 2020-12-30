package helper

type SyncConfig struct {
	// UseHttps will force https on each proto dependencies fetch.
	UseHttps bool

	// HomeDir is the home directory, used as root to find ssh identity files.
	HomeDir string

	// TargetDir is the dependencies directory where protodep.toml files are located.
	TargetDir string

	// OutputDir is the directory where proto files will be cloned.
	OutputDir string

	// BasicAuthUsername is used if `https` mode  is enable. Optional, only if dependency repository needs authentication.
	BasicAuthUsername string

	// BasicAuthPassword is used if `https` mode is enable. Optional, only if dependency repository needs authentication.
	BasicAuthPassword string

	// IdentityFile is used if `ssh` mode is enable. Optional, it is computed like {home}/.ssh/
	IdentityFile string

	// IdentityPassword is used if `ssh` mode is enable. Optional, only if identity file needs a passphrase.
	IdentityPassword string
}
