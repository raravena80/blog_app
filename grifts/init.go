package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/raravena80/blog_app/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
