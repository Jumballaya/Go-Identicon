package icon

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/color"
	"image/png"
	"log"
	"strconv"
	"strings"
)

// Identicon Struct
type Identicon struct {
	Input      string
	Hex        string
	Color      []string
	Pixelmap   [6][5]bool
	Padding    int
	SquareSize int
	Width      int
	Height     int
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// makeHex Creates 32 character hex string from the Identicon's input
func makeHex(input string) string {
	data := []byte(input)
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

// Build the pixel byte map
// Currently just all filled-in
func makePixels(squareSize int) []byte {
	pixels := make([]byte, squareSize)
	for i := 0; i < squareSize; i++ {
		pixels[i] = 1
	}
	return pixels
}

// Create colors based on a hex string
func makeColors(hex string) []color.NRGBA {
	red, err := strconv.ParseInt(hex[0:2], 16, 32)
	handleError(err)

	blue, err := strconv.ParseInt(hex[6:8], 16, 32)
	handleError(err)

	green, err := strconv.ParseInt(hex[12:14], 16, 32)
	handleError(err)

	squareColor := color.NRGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: 0xff}
	white := color.NRGBA{0xf1, 0xf1, 0xf1, 0xff}
	colors := []color.NRGBA{white, squareColor}
	return colors
}

// makeHalfMap Turns a string of length 32 into a nested array of booleans [3][5]bool
// 1. Pair the letters by two's and convert the hex (e.g. f8) to integer
// 2. Convert the 2 letter hex digit into a number
// 3. Limit to 15 characters for 3x5
// 4. Chunk the hex array into a 3x5 array
func makeHalfMap(hex string) [3][5]bool {
	letters := strings.Split(hex, "")
	hexArray := []bool{}
	for i := 0; i < len(letters); i = i + 2 {
		digit := letters[i : i+2]
		num, err := strconv.ParseInt(digit[0]+digit[1], 16, 32)
		handleError(err)
		hexArray = append(hexArray, num%2 == 0)
	}
	hexArray = hexArray[:15]
	halfGrid := [3][5]bool{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			halfGrid[i][j] = hexArray[(i)+(j)]
		}
	}
	return halfGrid
}

// mirrorMap Turns the [3][5] array into a [6][5]
func mirrorMap(halfMap [3][5]bool) [6][5]bool {
	pixelMap := [6][5]bool{}
	for i := 0; i < len(pixelMap); i++ {
		for j := 0; j < len(pixelMap[i]); j++ {
			if i < 3 {
				pixelMap[i][j] = halfMap[i][j]
			} else {
				pixelMap[i][j] = halfMap[5-i][j]
			}
		}
	}
	return pixelMap
}

// makeMap creates the hex and pixelMap given an input string
func makeMap(input string) (string, [6][5]bool) {
	hex := makeHex(input)
	pixelMap := mirrorMap(makeHalfMap(hex))
	return hex, pixelMap
}

// Render the image and return the bytes in png encoding
// 1. Create the hex and pixelMap values
// 2. Build the colors, rectagle, palette and the image
// 3. Loop through the pixelmap and build each square
// 4. Write the bytes to a buffer and return it the bytes
func (icon *Identicon) Render() []byte {
	hex, pixelMap := makeMap(icon.Input)
	icon.Hex = hex
	icon.Pixelmap = pixelMap

	colors := makeColors(icon.Hex)
	rect := image.Rect(0, 0, icon.Width, icon.Height)
	palette := color.Palette{colors[0], colors[1]}
	img := image.NewPaletted(rect, palette)

	for i := 0; i < len(icon.Pixelmap); i++ {
		for j := 0; j < len(icon.Pixelmap[i]); j++ {
			if icon.Pixelmap[i][j] {
				for k := 0; k < (icon.SquareSize); k++ {
					x := icon.Padding + (i * icon.SquareSize)
					y := icon.Padding + (j * icon.SquareSize) + k
					offset := img.PixOffset(x, y)
					copy(img.Pix[offset:], makePixels(icon.SquareSize))
				}
			}
		}
	}

	var buffer bytes.Buffer
	png.Encode(&buffer, img)
	return buffer.Bytes()
}
