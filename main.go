package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/golang/freetype"
)

//go:embed square.ttf
var square []byte

func confirmOverwrite() bool {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("out.png already exists. Overwrite? [y/N]: ")
	res, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	if len(res) < 2 {
		return false
	}
	return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
}

const WIDTH = 2480
const HEIGHT = 3508
const RATIO = float64(HEIGHT) / WIDTH

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("Expected single file argument")
		return
	}
	filename := flag.Arg(0)

	if _, err := os.Stat(filename); err == nil {
		if !confirmOverwrite() {
			log.Fatal("Exiting...")
		}
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var fileLines []string

	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	maxWidth := 0
	for _, s := range fileLines {
		maxWidth = max(maxWidth, len(s))
	}
	fmt.Printf("Width x Height: %dx%d\n", maxWidth, len(fileLines))

	factor := max(float64(maxWidth)*RATIO, float64(len(fileLines)))
	size := 800.0 / factor

	f, err := freetype.ParseFont(square)
	if err != nil {
		log.Println(err)
		return
	}

	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(300)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	pt := freetype.Pt(50, 50+int(c.PointToFixed(size)>>6))
	for _, s := range fileLines {
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(size)
	}

	outFile, err := os.Create("out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}
