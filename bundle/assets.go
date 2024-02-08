package bundle

import "io/fs"

// SEAssets provides behaviors to interact with asset files that have been loaded
// into a config bundle. These behaviors only work if Orchestrator config is bundled
// via remote configuration.
// If the config files are on the local filesystem, then we should interact with
// those directly instead of using this interface.
type SEAssets interface {
	// FS returns a filesystem containing the assets associated with the service
	// extension.
	FS() (fs.FS, error)
	// ReadFile returns the file contents of a file from the service extension
	// assets.
	ReadFile(string) ([]byte, error)
}
