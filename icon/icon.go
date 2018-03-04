package icon


import (
  "encoding/hex"
  "crypto/md5"
  "strings"
  "strconv"
  "image"
  "image/color"
  "image/png"
  "bytes"
)


// Identicon Struct
type Identicon struct {
  Input string
  Hex string
  Color []string
  Pixelmap [3][5]bool
  Padding int
  SquareSize int
  Width int
  Height int
}


// Identicon.CreateHex()
// Creates 32 character hex string from the Identicon's input
func (icon *Identicon) CreateHex() {
  data := []byte(icon.Input)
  sum := md5.Sum(data)
  hex := hex.EncodeToString(sum[:])
  icon.Hex = hex
}


// Identicon.CreateMap()
// Creates a nested boolean array
func (icon *Identicon) CreateMap() {
  letters := strings.Split(icon.Hex, "")
  hex_array := []bool{}
  // Pair the letters by two's and convert the hex (e.g. f8) to integer
  for i := 0; i < len(letters); i = i + 2 {
    digit := letters[i: i + 2]
    num, _ := strconv.ParseInt(digit[0], 16, 32)

    // This is the criteria for which squares get color
    // Could be anything else e.g.
    // (num < 4 || num > 14) || (num % 3 == 0)
    hex_array = append(hex_array, num % 2 == 0)
  }
  hex_array = hex_array[:15]

  // Build the grid, we will mirror it later
  half_grid := [3][5]bool{}
  for i := 0; i < 3; i++ {
    for j := 0; j < 5; j++ {
      half_grid[i][j] = hex_array[(i) + (j)]
    }
  }
  icon.Pixelmap = half_grid
}

// Build the pixel byte map
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
  primary_red, _:= strconv.ParseInt(hex[0:2], 16, 32)
  primary_blue, _:= strconv.ParseInt(hex[6:8], 16, 32)
  primary_green, _:= strconv.ParseInt(hex[12:14], 16, 32)
  color_primary := color.NRGBA{
    R: uint8(primary_red),
    G: uint8(primary_green),
    B: uint8(primary_blue),
    A: 0xff,
  }

  // White
  white := color.NRGBA{0xf1, 0xf1, 0xf1, 0xff}

  colors := []color.NRGBA{ white, color_primary }
  return colors
}


// Render the image and return the bytes
func (icon *Identicon) Render() []byte {
  // Create the hex and the map
  icon.CreateHex()
  icon.CreateMap()

  // Build the colors, rectagle, palette and the image
  colors := makeColors(icon.Hex)
  rect := image.Rect(0, 0, icon.Width, icon.Height)
  palette := color.Palette{ colors[0], colors[1] }
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
          x = icon.Padding + (5 - i) * icon.SquareSize // The 5 comes from the 5 row, this effect the column count
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
