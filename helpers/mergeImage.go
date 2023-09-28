package helpers

import (
	gim "github.com/ozankasikci/go-image-merge"
	"image/png"
	"log"
	"os"
)

func MergeImage(basePath string, imagePaths []string) (outputPath string, err error) {
	grids := make([]*gim.Grid, 0)
	for i, imagePath := range imagePaths {
		grid := &gim.Grid{
			ImageFilePath: imagePath,
		}
		grids[i] = grid
	}

	rgba, err := gim.New(grids, 1, 6).Merge()
	if err != nil {
		log.Printf("merge image faild, err [%+v]", err)
		return "", err
	}

	file, err := os.Create(basePath + "/output.png")
	if err != nil {
		log.Printf("create output.png  faild, err [%+v]", err)
		return "", err
	}
	err = png.Encode(file, rgba)
	outputPath = basePath + "/output.png"
	return outputPath, err
}
