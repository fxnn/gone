// Package router provides basic means of switching between the wiki's
// different page types.
// Those types are, for example, viewing a page, editing a page, logging in
// etc., and are called "modes".
//
// This package decides when to invoke the HTTP handler for which mode and
// provides logic for redirecting betweend those modes.
//
// Template URLs
//
// A special case are URLs to files that are part of gone's templates.
// These may reside in another namespace than the files delivered regularly,
// e.g. when template data from inside the gone binary is be used.
//
// The question is here, how to distinguish those files from a given URL.
// This is either possible by a special path component, e.g. "/+template",
// or by a query parameter, e.g. "?template".
//
// While making this decision, we should never break the spec in that a
// template can be a different resource than a regularly delivered file with
// the same name.
// Therefore, our means should be able to identify different resources.
//
// This is possible with paths (of course), but as the spec states, also with
// query parameters. See RFC 3986 Sec. 3.4:
// "The query component contains non-hierarchical data that, along with
// data in the path component (Section 3.3), serves to identify a
// resource within the scope of the URI's scheme and naming authority
// (if any)."
//
// As in gone, the path shall solely identify a file inside the content root,
// and we already use the notion of "modes" in query parameters, we decide to
// also use modes for templates.
//
package router
