package main

import (
	"fmt"
	"github.com/Jumballaya/identicon/icon"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
)

var output string

func makeIdenticon(input string) {
	icn := &icon.Identicon{Input: input, Padding: 4, SquareSize: 82, Width: 500, Height: 418}
	png := icn.Render()
	ioutil.WriteFile(output+"/"+input+".png", png, os.FileMode(0644))
}

func main() {
	output = "./"

	app := cli.NewApp()
	app.Name = "Identicon"
	app.Usage = "Create Identicon images"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "output, o",
			Value:       "./",
			Usage:       "output destination for the identicon images",
			Destination: &output,
		},
	}

	app.Action = func(c *cli.Context) error {
		input := c.Args().Get(0)
		if len(input) > 0 {
			makeIdenticon(input)
			fmt.Printf("\nIdenticon Created\n\n")
		} else {
			fmt.Printf("\nError: Input required to make image\n\n")
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
