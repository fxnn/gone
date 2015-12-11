// Package filestore implements the storage via filesystem.
// Provided features are reading and writing content, creating and deleting
// files as well as realization of access control.
//
// It is important to note, that the error handling is implemented using a
// saved error value.
// As soon as an error occured, all functions get no-ops.
// This way, implementation gets a lot easier and more readable.
// However, callers and developers always have to ensure to correctly check
// for errors, as the type system won't ensure this anymore.
package filestore
