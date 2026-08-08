module github.com/macawls/ogre/v3

go 1.25.0

require (
	golang.org/x/image v0.38.0
	golang.org/x/net v0.52.0
)

require golang.org/x/text v0.36.0

require github.com/go-text/typesetting v0.3.4

retract (
	v3.2.0 // Module path is missing the /v3 suffix.
	v3.1.0 // Module path is missing the /v3 suffix.
	v3.0.0 // Module path is missing the /v3 suffix.
)
