package assets

import (
	"bytes"
	_ "embed" // for the go:embed directive
)

// GQLPlaygroundHTML stores the content of the playground/index.html file.
//go:embed playground/index.html
var GQLPlaygroundHTML []byte

// InitGQLPlaygroundHTML prepares the playground and puts the right port in there.
func InitGQLPlaygroundHTML(port []byte, endpoint []byte) {
	GQLPlaygroundHTML = bytes.Replace(GQLPlaygroundHTML, []byte("{{.port}}"), port, 1)
	GQLPlaygroundHTML = bytes.Replace(GQLPlaygroundHTML, []byte("{{.endpoint}}"), endpoint, 1)
}
