// +build !development

package server

// DisableVerification won't pass through invalid JWT in production mode.
var DisableVerification = true
