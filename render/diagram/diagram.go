// https://tromp.github.io/cl/diagrams.html
package diagram

import (
	"image"
	"image/color"
)

type Config struct {
	VarLines color.Color
}

func (c *Config) Draw() (image.Image, error) { return nil, nil }
