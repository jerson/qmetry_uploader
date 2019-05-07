package commands

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"qmetry_uploader/modules/config"
	"sort"

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/disintegration/imaging"
)

// Compress ...
func Compress() error {
	baseDir := config.Vars.Dir.Input
	return CompressDir(baseDir)
}

// CompressDir ...
func CompressDir(baseDir string) error {

	scenarios, err := ReportDir(baseDir)
	if err != nil {
		return err
	}

	caseGroup := map[string]Case{}
	for _, scenario := range scenarios {
		for _, caseItem := range scenario.Cases {
			key := fmt.Sprintf("%s_%s (%d steps)", caseItem.Name, caseItem.Device, len(caseItem.Steps))
			caseGroup[key] = caseItem

		}
	}

	var caseOptions []string
	for key := range caseGroup {
		caseOptions = append(caseOptions, key)
	}
	sort.Strings(caseOptions)

	var caseOptionsSelected []string
	prompt := &survey.MultiSelect{
		Message:  "Choose cases",
		Options:  caseOptions,
		PageSize: 10,
		Default:  caseOptions,
	}
	err = survey.AskOne(prompt, &caseOptionsSelected, nil)
	if err != nil {
		return err
	}

	for i, caseOption := range caseOptionsSelected {
		caseItem := caseGroup[caseOption]

		fmt.Println(fmt.Sprintf("[%d/%d] processing: %s_%s (%d steps)", i, len(caseOptionsSelected), caseItem.Name, caseItem.Device, len(caseItem.Steps)))

		var filePaths []string
		for _, step := range caseItem.Steps {
			filePaths = append(filePaths, step.Path)
		}

		_ = os.MkdirAll(config.Vars.Dir.Output, 0777)
		output := fmt.Sprintf("%s/%s_%s.png", config.Vars.Dir.Output, caseItem.Device, caseItem.Name)
		err := mergeImages(filePaths, output)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(fmt.Sprintf("output: %s\n", output))

	}

	return nil
}

// ImageFixed ...
type ImageFixed struct {
	Image *image.NRGBA
	Point image.Point
}

func mergeImages(paths []string, output string) error {

	resizeHeight := 600
	x := 0
	y := 0
	width := 0
	height := resizeHeight

	var images []ImageFixed
	for i, path := range paths {
		img, err := decode(path)
		if err != nil {
			return err
		}

		src := imaging.Resize(img, 0, resizeHeight, imaging.Lanczos)
		images = append(images, ImageFixed{
			Point: image.Pt(x, y),
			Image: src,
		})
		x += src.Bounds().Max.X
		if x > width {
			width = x
		}
		if math.Mod(float64(i+1), 3) == 0 {
			x = 0
			y += resizeHeight
			if i < len(paths)-1 {
				height += resizeHeight
			}
		}

	}

	target := imaging.New(width, height, color.NRGBA{0, 0, 0, 0})
	for _, imageFixed := range images {
		target = imaging.Paste(target, imageFixed.Image, imageFixed.Point)
	}
	err := imaging.Save(target, output)

	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	return nil
}

func decode(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}
