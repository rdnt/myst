package hashicon

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/ajstarks/svgo"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidLen is returned if the passed byte slice is not a power of 2
	// or is outside of the allowed bounds (256-65536 bits inclusive)
	ErrInvalidLen = fmt.Errorf("invalid len")
)

// Hashicon represents a hash in image form, and is based on the concepts of
// identicons and randomart's drunken bishop algorithm.
type Hashicon struct {
	Stride int
	Pix    []float64
}

// New creates a new hashicon and stores its pix array and stride.
// The length of the pix slice always equals stride^2.
func New(b []byte) (*Hashicon, error) {
	if len(b) < 32 || len(b) > 8192 || !((len(b) & (len(b) - 1)) == 0) {
		return nil, ErrInvalidLen
	}

	// make sure the stride is a power of two (that will also have an integer
	// logarithm of 2)
	stride := len(b) * 2
	i, frac := math.Modf(math.Sqrt(float64(stride)))
	if frac == 0 {
		stride = int(i)
	} else {
		i, _ := math.Modf(math.Sqrt(float64(stride / 2)))
		stride = int(i)
	}

	pix := make([]float64, stride*stride)

	// bits per coordinate indicates how many bits should be used to determine
	// the starting position of the
	bpc := int(math.Log2(float64(stride)))

	// find x pos
	x := 0
	for i := 0; i < bpc; i++ {
		x = x + getBit(b[i/8], i%8)<<(bpc-i-1)
	}

	// find y pos
	y := 0
	for i := 0; i < bpc; i++ {
		y = y + getBit(b[(i+bpc)/8], (i+bpc)%8)<<(bpc-i-1)
	}

	// perform the drunken bishop algorithm (modified for movements in
	// top, left, right, bottom instead of diagonally)
	for i := bpc * 2; i < len(b)*8; i += 2 {
		b1 := getBit(b[i/8], i%8) * 2
		b2 := getBit(b[(i+1)/8], (i+1)%8)

		// perform the movement but don't exit the bounds of the grid
		switch b1 + b2 {
		case 0:
			if y != 0 {
				// move top
				y = y - 1
			}
		case 1:
			if x != stride-1 {
				// move right
				x = x + 1
			}
		case 2:
			if y != stride-1 {
				// move bottom
				y = y + 1
			}
		case 3:
			if x != 0 {
				// move left
				x = x - 1
			}
		}

		pix[x*stride+y] += 1
	}

	// find max pixel and normalize all pixels in range 0..1
	max := 0.0
	for _, p := range pix {
		if p > max {
			max = p
		}
	}
	for i, p := range pix {
		pix[i] = p / max
	}

	return &Hashicon{
		Stride: stride,
		Pix:    pix,
	}, nil
}

// WeightToColor converts a normalized weight (in the range 0..1) to a color.
// Can be overwritten to allow for custom visualizations.
func WeightToColor(w float64) color.NRGBA {
	switch {
	case w > .5:
		return color.NRGBA{
			R: 32,
			G: 233,
			B: 183,
			A: uint8(w * 255),
		}
	default:
		return color.NRGBA{
			R: 58,
			G: 141,
			B: 153,
			A: uint8(math.Pow(w, 1.66) * 255),
		}
	}
}

// ToSVG returns an SVG based on the hashicon's pixel data.
func (h *Hashicon) ToSVG() string {
	// try to be the same size in the resulting SVG regardless of stride
	size := 256 / h.Stride

	buf := new(bytes.Buffer)

	canvas := svg.New(buf)
	canvas.Start(h.Stride*size, h.Stride*size)

	canvas.Polygon([]int{0, h.Stride * size, h.Stride * size, 0}, []int{0, 0, h.Stride * size, h.Stride * size}, "fill: #181b21")

	for i, p := range h.Pix {
		x := (i % h.Stride) * size
		y := (i / h.Stride) * size

		clr := WeightToColor(p)

		xs := []int{x, x + size, x + size, x}
		ys := []int{y, y, y + size, y + size}

		canvas.Polygon(xs, ys, fmt.Sprintf("fill:#%02X%02X%02X; opacity:%f", clr.R, clr.G, clr.B, float64(clr.A)/255.0))
	}

	canvas.End()

	return buf.String()
}

// Export converts the hashicon to SVG and saves it in the specified path.
func (h *Hashicon) Export(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer f.Close()

	_, err = f.Write([]byte(h.ToSVG()))
	if err != nil {
		return errors.Wrap(err, "failed to write to file")
	}

	return nil
}

// Returns 1 or 0 depending on the bit specified in the given byte's position.
func getBit(b byte, i int) int {
	if i < 0 || i > 7 {
		return 0
	}
	return int(b >> (8 - i - 1) & 1)
}
