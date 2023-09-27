// bird.go
package scenes

import (
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/faiface/pixel"
)

// Bird representa un pájaro en el juego
type Bird struct {
	X, Y       float64
	Frame      int
	Frames     []*pixel.Sprite
	CurrentPic int
	Direction  int
}

func loadBirdSprites() ([]*pixel.Sprite, error) {
	var birdSprites []*pixel.Sprite

	// Ruta a la carpeta de assets
	assetsDir := "assets/pajaritos"

	// Función para cargar una imagen y convertirla en un pixel.Picture
	loadPicture := func(path string) (pixel.Picture, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			return nil, err
		}
		return pixel.PictureDataFromImage(img), nil
	}

	err := filepath.Walk(assetsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Si es un archivo en lugar de un directorio, intenta cargarlo como una imagen
			pic, err := loadPicture(path)
			if err != nil {
				// Manejar el error si no se puede cargar la imagen
				log.Printf("Error cargando imagen %s: %v\n", path, err)
			} else {
				// Crear un Sprite y agregarlo a la lista de Sprites de pájaros
				sprite := pixel.NewSprite(pic, pic.Bounds())
				birdSprites = append(birdSprites, sprite)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return birdSprites, nil
}
