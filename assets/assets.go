package assets

import "io/fs"

type SEAssets interface {
	// FS returns a filesystem containing the assets associated with the service
	// extension.
	FS() (fs.FS, error)
	// ReadFile returns the file contents of a file from the service extension
	// assets.
	ReadFile(string) ([]byte, error)
}
