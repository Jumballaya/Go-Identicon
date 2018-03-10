package icon

// @TODO: I need to refactor buildHalfMap it still feels clumsy and I
//				am not so sure that It needs to be a method anyway. It
//        should be refactored into smaller functions.
//        This also needs to be tighter when it comes to typing

// @TODO: I should mirror the half_map and set PixelMap to 6x5
//        before getting the the Render method.
//
//        Simplify the Render method by moving most of the meat to
//        functions and better utilize SquareSize, Width and Height.
//        Set the color to the color type. Change the Hex property or
//        get rid of it, I am not sure I really need it.
//
//        Allow the CLI tool to override the color, square size, height,
//        width or padding.

// @TODO: Really need to refactor the creation of the img.Pix array. I need to really
//        learn the way it is built and how to actually form the slice (array?)
//

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"image"
	"image/color"
	"image/png"
	"strconv"
	"strings"
)

// Identicon Struct
type Identicon struct {
	Input      string
	Hex        string
	Color      []string
	Pixelmap   [3][5]bool
	Padding    int
	SquareSize int
	Width      int
	Height     int
}

// makeHex()
// Creates 32 character hex string from the Identicon's input
func makeHex(input string) string {
	data := []byte(input)
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

// Build the pixel byte map
// Currently just all filled-in
func makePixels(square_size int) []byte {
	pixels := make([]byte, square_size)
	for i := 0; i < square_size; i++ {
		pixels[i] = 1
	}
	return pixels
}

// Create colors based on a hex string
func makeColors(hex string) []color.NRGBA {
	// Primary
	primary_red, _ := strconv.ParseInt(hex[0:2], 16, 32)
	primary_blue, _ := strconv.ParseInt(hex[6:8], 16, 32)
	primary_green, _ := strconv.ParseInt(hex[12:14], 16, 32)
	color_primary := color.NRGBA{R: uint8(primary_red), G: uint8(primary_green), B: uint8(primary_blue), A: 0xff}
	white := color.NRGBA{0xf1, 0xf1, 0xf1, 0xff}
	colors := []color.NRGBA{white, color_primary}
	return colors
}

// Turn a slice of letters into a nested array of booleans [3][5]bool
// @TODO: Refactor this
func makeHalfMap(hex string) [3][5]bool {
	letters := strings.Split(hex, "")
	hex_array := []bool{}

	// Pair the letters by two's and convert the hex (e.g. f8) to integer
	for i := 0; i < len(letters); i = i + 2 {
		digit := letters[i : i+2]
		num, _ := strconv.ParseInt(digit[0], 16, 32)
		hex_array = append(hex_array, num%2 == 0)
	}
	hex_array = hex_array[:15]

	half_grid := [3][5]bool{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			half_grid[i][j] = hex_array[(i)+(j)]
		}
	}
	return half_grid
}

// Mirror Half Grid
// Turn the [3][5] array into a [6][5]
func mirrorHalfGrid([3][5]bool) [6][5]bool {
	return [6][5]bool{}
}

// Render the image and return the bytes
func (icon *Identicon) Render() []byte {
	// @TODO:
	// Create the hex and the map
	// Change this to:
	// halfMap, hex := createMap(icon.Input)
	// icon.Color = createColor(hex)
	// icon.PixelMap = mirrorHalfMap(halfMap)
	// ... go through the pixelmap and build the image bytes
	// png.Encode and return the bytes

	icon.Hex = makeHex(icon.Input)

	// icon.Pixelmap = mirrorHalfGrid(makeHalfMap(icon.Input))
	icon.Pixelmap = makeHalfMap(icon.Input)

	// Build the colors, rectagle, palette and the image
	colors := makeColors(icon.Hex)
	rect := image.Rect(0, 0, icon.Width, icon.Height)
	palette := color.Palette{colors[0], colors[1]}
	img := image.NewPaletted(rect, palette)

	// Loop throught the pixelmap and build each square
	for i := 0; i < len(icon.Pixelmap); i++ {
		for j := 0; j < len(icon.Pixelmap[i]); j++ {

			// If the square is true then color it
			if icon.Pixelmap[i][j] {

				// Go column by column by column and fill the square
				for k := 0; k < (icon.SquareSize); k++ {

					// create the x and y coordinates to start
					x := icon.Padding + (i * icon.SquareSize)
					y := icon.Padding + (j * icon.SquareSize) + k
					offset := img.PixOffset(x, y)
					copy(img.Pix[offset:], makePixels(icon.SquareSize))

					// Mirror the process on the right side
					x = icon.Padding + (5-i)*icon.SquareSize // The 5 comes from the 5 row, this effect the column count
					offset = img.PixOffset(x, y)
					copy(img.Pix[offset:], makePixels(icon.SquareSize))
				}
			}
		}
	}

	// Write the bytes to a buffer and return it the bytes
	var buffer bytes.Buffer
	png.Encode(&buffer, img)
	return buffer.Bytes()
}
