package hashicon

import (
	"fmt"
	"image/color"
	"math"
	"os"
)

var (
	ErrInvalidLen = fmt.Errorf("invalid len")
)

type Hashicon struct {
	Stride int
	Pix    []float64
}

// New creates and saves a new user
func New(b []byte) (*Hashicon, error) {
	// at least 256-bit, only powers of two
	bits := len(b) * 8
	fmt.Println("len", len(b))
	fmt.Println("bits", bits)
	if len(b) < 32 || !((len(b) & (len(b) - 1)) == 0) {
		return nil, ErrInvalidLen
	}

	stride := len(b) * 2
	i, frac := math.Modf(math.Sqrt(float64(stride)))
	if frac == 0 {
		stride = int(i)
	} else {
		i, _ := math.Modf(math.Sqrt(float64(stride / 2)))
		stride = int(i)
	}

	fmt.Println("stride", stride)

	h := &Hashicon{
		Stride: stride,
		Pix:    make([]float64, stride*stride),
	}

	fmt.Println("BEFORE PARSE")
	h.parse(b)

	return h, nil
}

func (h *Hashicon) parse(b []byte) {
	//bpc := h.Stride >> 1
	//fmt.Println("bpc", bpc)

	//x := int(b[0] & 0b11100000 >> 5)
	//y := int(b[0] & 0b00011100 >> 2)

	//x := h.Stride / 2
	//y := h.Stride / 2

	//x := 0

	//stride := int(math.Sqrt(float64(len(b) * 4)))
	//fmt.Println(x, y)

	bpc := int(math.Log2(float64(h.Stride)))
	fmt.Println("bpc", bpc)

	fmt.Printf("%08b\n", b[0])
	x := 0
	for i := 0; i < bpc; i++ {
		x = x + getBit(b[i/8], i%8)<<(bpc-i-1)
	}

	y := 0
	for i := 0; i < bpc; i++ {
		y = y + getBit(b[(i+bpc)/8], (i+bpc)%8)<<(bpc-i-1)
	}

	fmt.Println(x, y)

	//locating := true
	for i := bpc * 2; i < len(b)*8; i += 2 { // loop all bits
		b1 := getBit(b[i/8], i%8) * 2
		b2 := getBit(b[(i+1)/8], (i+1)%8)

		switch b1 + b2 {
		case 0:
			if y != 0 {
				// move top
				y = y - 1
			}
			break
		case 1:
			if x != h.Stride-1 {
				// move right
				x = x + 1
			}
			break
		case 2:
			if y != h.Stride-1 {
				// move bottom
				y = y + 1
			}
			break
		case 3:
			if x != 0 {
				// move left
				x = x - 1
			}
			break
		}

		h.Pix[x*h.Stride+y] += 1
	}

	max := 0.0
	for _, p := range h.Pix {
		if p > max {
			max = p
		}
	}
	for i, p := range h.Pix {
		h.Pix[i] = p / max
	}

	//fmt.Println(h.Pix)

}

//func (h *Hashicon) parseOld(b []byte) {
//	for i := 0; i < len(b)*8; i++ { // loop all bits
//		byt := b[i/8]
//		bit := getBit(byt, i%8)
//		h.Pix[i/h.BitsPerPixel] += bit / float64(h.BitsPerPixel)
//	}
//}

//func weightToColor(w float64, step float64) color.RGBA {
//	switch {
//	case w > .75:
//		return color.RGBA{
//			R: 32,
//			G: 233,
//			B: 183,
//			A: uint8(w * 255),
//		}
//	default:
//		return color.RGBA{
//			R: 58,
//			G: 141,
//			B: 153,
//			A: uint8(w * 255),
//		}
//	}
//}

func weightToColor(w float64) color.RGBA {
	switch {
	case w > .66:
		return color.RGBA{
			R: 32,
			G: 233,
			B: 183,
			A: uint8(w * 255),
		}
	default:
		return color.RGBA{
			R: 58,
			G: 141,
			B: 153,
			A: uint8(math.Pow(w, 1.66) * 255),
		}
	}
}

func (h *Hashicon) ToSVG() string {
	mult := 256 / h.Stride
	svg := fmt.Sprintf(`<svg width="%d" height="%d" version="1.1" xmlns="http://www.w3.org/2000/svg">`, h.Stride*mult, h.Stride*mult)
	svg += `<rect width="100%" height="100%" fill="#181b21"/>`
	for i, p := range h.Pix {
		x := i % h.Stride
		y := i / h.Stride
		clr := weightToColor(p)

		svg += fmt.Sprintf(
			`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" opacity="%f" />`,
			x*mult,
			y*mult,
			mult,
			mult,
			fmt.Sprintf("#%X%X%X", clr.R, clr.G, clr.B),
			float64(clr.A)/255,
		)
	}
	svg += `</svg>`
	return svg
}

func (h *Hashicon) SaveSVGIncremental(name string) error {

	f, err := os.Create(fmt.Sprintf("./tmp/hashicons/hashicon-%s.svg", name))
	if err != nil {
		return err
	}

	mult := 256 / h.Stride
	f.Write([]byte(fmt.Sprintf(`<svg width="%d" height="%d" version="1.1" xmlns="http://www.w3.org/2000/svg">`, h.Stride*mult, h.Stride*mult)))
	f.Write([]byte(`<rect width="100%" height="100%" fill="#181b21"/>`))
	for i, p := range h.Pix {
		x := i % h.Stride
		y := i / h.Stride
		clr := weightToColor(p)

		f.Write([]byte(fmt.Sprintf(
			`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" opacity="%f" />`,
			x*mult,
			y*mult,
			mult,
			mult,
			fmt.Sprintf("#%X%X%X", clr.R, clr.G, clr.B),
			float64(clr.A)/255,
		)))
	}
	f.Write([]byte(`</svg>`))

	return f.Close()

}

func (h *Hashicon) Save(name string) error {
	f, err := os.Create(fmt.Sprintf("./tmp/hashicons/hashicon-%s.svg", name))
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(h.ToSVG()))
	if err != nil {
		return err
	}
	return f.Close()
}

func getBit(b byte, i int) int {
	if i < 0 || i > 7 {
		return 0
	}
	return int(b >> (8 - i - 1) & 1)
}

func getBitOld(b byte, i int) float64 {
	if i < 0 || i > 7 {
		return 0
	}
	return float64(b >> (8 - i - 1) & 1)
}
