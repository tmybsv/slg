package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"

	"os"
)

const enASCII = 97

func generateAlphabet(lang string) []string {
	alp := []string{}
	switch lang {
	case "en":
		for i := 0; i < 26; i++ {
			alp = append(alp, string(rune(enASCII+i)))
		}
	}
	return alp
}

func parseImageNames(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	names := []string{}
	for _, e := range entries {
		names = append(names, dir+e.Name())
	}
	return names, nil
}

func generateLettersToImagesMap(style, lang, dir string) (map[string]string, error) {
	names, err := parseImageNames(dir)
	if err != nil {
		return nil, err
	}

	m := map[string]string{}
	switch style {
	case "ugly":
		alp := generateAlphabet(lang)
		for i, l := range alp {
			m[l] = names[i]
		}
	default:
		return nil, errors.New("unknown style")
	}
	return m, nil
}

var (
	style     string
	lang      string
	assetsDir string
	outDir    string
)

func initFlags() {
	flag.StringVar(&style, "style", "ugly", "generated word style")
	flag.StringVar(&style, "s", "ugly", "generated word style")
	flag.StringVar(&lang, "lang", "en", "generated language")
	flag.StringVar(&lang, "l", "en", "generated language")
}

func main() {
	initFlags()
	flag.Parse()

	assetsDir, exists := os.LookupEnv("ASSETS_DIR")
	if !exists {
		assetsDir = "assets/en/ugly/"
	}

	lettersToImages, err := generateLettersToImagesMap(style, lang, assetsDir)
	if err != nil {
		fmt.Printf("slg: failed to generate map: %v\n", err)
		os.Exit(1)
	}

	if flag.NArg() < 1 {
		fmt.Println("slg: failed to parse word")
		os.Exit(1)
	}
	word := flag.Arg(0)
	imgs := []image.Image{}
	for _, l := range word {
		imgFile, err := os.Open(lettersToImages[string(l)])
		if err != nil {
			fmt.Printf("slg: failed to open image: %v\n", err)
			os.Exit(1)
		}
		defer imgFile.Close()

		img, err := png.Decode(imgFile)
		if err != nil {
			fmt.Printf("slg: failed to decode image: %v\n", err)
			os.Exit(1)
		}

		imgs = append(imgs, img)
	}

	var (
		outWidth  int
		outHeight int
	)
	for _, i := range imgs {
		outWidth += i.Bounds().Dx()

		if outHeight < i.Bounds().Dy() {
			outHeight = i.Bounds().Dy()
		}
	}

	outImg := image.NewRGBA(image.Rect(0, 0, outWidth, outHeight))
	draw.Draw(outImg, imgs[0].Bounds(), imgs[0], image.Point{}, draw.Src)
	offset := image.Point{}
	for _, i := range imgs[1:] {
		offset = offset.Add(image.Pt(i.Bounds().Dx(), 0))
		draw.Draw(outImg, i.Bounds().Add(offset), i, image.Point{}, draw.Src)
	}

	outImgFile, err := os.Create(fmt.Sprintf("%s.png", word))
	if err != nil {
		fmt.Printf("slg: failed to output image: %v\n", err)
		os.Exit(1)
	}
	defer outImgFile.Close()

	if err := png.Encode(outImgFile, outImg); err != nil {
		fmt.Printf("slg: failed to encode image: %v\n", err)
		os.Exit(1)
	}
}
