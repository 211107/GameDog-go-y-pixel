package scenes

import (
	_ "image/jpeg"
	_ "image/png"

	"github.com/faiface/pixel"
)

//  representa un perro en el juego
type Dog struct {
	X, Y, Vy float64
	Sprite   *pixel.Sprite
}

// setDogSprite establece el sprite del perro según el estado
func setDogSprite(dog *Dog, spriteMap map[string]*pixel.Sprite, state string) {
	switch state {
	case "Right1":
		dog.Sprite = spriteMap["dog2"]
	case "Right2":
		dog.Sprite = spriteMap["dog3"]
	case "Left1":
		dog.Sprite = spriteMap["dog5"]
	case "Left2":
		dog.Sprite = spriteMap["dog6"]
	case "JumpingRight":
		dog.Sprite = spriteMap["dog7"]
	case "JumpingLeft":
		dog.Sprite = spriteMap["dog8"]
	case "NeutralLeft":
		dog.Sprite = spriteMap["dog4"]
	default:
		dog.Sprite = spriteMap["dog1"]
	}
}
