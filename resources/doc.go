// Package resources manages non-code data used by the application.
// Notably, there is a set of files bundled in the gone executable.
// Those files are kept in the "static" subdirectory.
//
// Use the generate-code.sh script to map those files into compilable code.
package resources

// Convert files into go code
//go:generate esc -pkg resources -prefix static -o static.go static

// List all converted files, as esc is currently not capable of doing so
//go:generate /bin/sh -c "echo \"package resources\nvar AllFileNames = []string{\" >allfilenames.go"
//go:generate /bin/sh -c "find static -type f -printf \"\\t\\\"/%P\\\",\\n\" >>allfilenames.go"
//go:generate /bin/sh -c "echo \"}\" >>allfilenames.go"
