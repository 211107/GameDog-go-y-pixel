package scenes

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var mu sync.Mutex
var gameEnded bool = false
var gameWon bool = false
var walkingFrame int = 1
var frameCounter int = 0
var winMessageShown bool = false

const tiempoMaximo = 2035 // Tiempo máximo
var tiempoTranscurrido int

func loadPicture(path string) (pixel.Picture, error) {
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

func Ini() {
	cfg := pixelgl.WindowConfig{

		Title:  "Juego del perrito",
		Bounds: pixel.R(0, 0, 800, 600),
	}

	win, err := pixelgl.NewWindow(cfg)

	if err != nil {

		panic(err)

	}

	var wg sync.WaitGroup

	birdSprites, err := loadBirdSprites()

	if err != nil {

		log.Fatalf("Error cargando imágenes de pájaros: %v", err)

	}

	bird := Bird{
		X:         rand.Float64() * win.Bounds().Max.X,
		Y:         win.Bounds().Max.Y - 50, // Ajusta la posición vertical
		Direction: 1,                       // Empieza volando hacia la derecha
	}

	bird.Frames = birdSprites

	dog := Dog{X: 100, Y: 150, Vy: 0}
	bones := make([]Bone, 0)
	score := 0

	bgPic, _ := loadPicture("assets/im.jpg")
	bgSprite := pixel.NewSprite(bgPic, bgPic.Bounds())

	bonePic, _ := loadPicture("assets/bones.png")
	boneSprite := pixel.NewSprite(bonePic, bonePic.Bounds())

	spriteMap := make(map[string]*pixel.Sprite)

	spriteFiles := []string{"dog1", "dog2", "dog3", "dog4", "dog5", "dog6", "dog7", "dog8"}
	for _, fname := range spriteFiles {
		pic, _ := loadPicture("assets/" + fname + ".png")
		spriteMap[fname] = pixel.NewSprite(pic, pic.Bounds())
	}

	setDogSprite(&dog, spriteMap, "NeutralRight")

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	// Goroutine movimiento del perro y recoleccion de huesos
	wg.Add(1)
	go func() {
		defer wg.Done()
		gravity := -0.5
		for !win.Closed() {
			mu.Lock()
			tiempoTranscurrido++

			frameCounter++
			if frameCounter >= 10 {
				walkingFrame = 3 - walkingFrame
				frameCounter = 0
			}

			if tiempoTranscurrido >= tiempoMaximo {
				gameEnded = true
				winMessageShown = true
			}

			if !gameEnded {
				if dog.Y <= 150 || (dog.Y <= 350 && dog.Y >= 300 && dog.X >= 300 && dog.X <= 500) {
					if win.JustPressed(pixelgl.KeySpace) {
						dog.Vy = 15
					}
				}

				dog.Y += dog.Vy
				if dog.Y > 150 {
					dog.Vy += gravity
				} else {
					dog.Y = 150
					dog.Vy = 0
				}

				toRemove := -1
				for i, bone := range bones {
					if (bone.X-dog.X)*(bone.X-dog.X)+(bone.Y-dog.Y)*(bone.Y-dog.Y) < 2500 {
						score++
						toRemove = i
						break
					}
				}

				if toRemove != -1 {
					bones = append(bones[:toRemove], bones[toRemove+1:]...)
				}

				if score >= 20 && tiempoTranscurrido < tiempoMaximo {
					// El jugador ha ganado
					gameEnded = true
					gameWon = true
					winMessageShown = true
				}

			} else {
				if win.JustPressed(pixelgl.KeyR) {
					// Reiniciar el juego
					gameEnded = false
					gameWon = false
					score = 0
					dog.X = 100
					dog.Y = 150
					bones = nil
					winMessageShown = false
					tiempoTranscurrido = 0
				}
			}

			mu.Unlock()
			time.Sleep(time.Millisecond * 16)
		}
	}()

	//gorou movimiento de las teclas del control

	wg.Add(1)
	go func() {
		defer wg.Done()
		for !win.Closed() { //inc
			mu.Lock()
			if !gameEnded {
				if win.Pressed(pixelgl.KeyLeft) && dog.X >= 25 {
					dog.X -= 10
				}
				if win.Pressed(pixelgl.KeyRight) && dog.X <= 775 {
					dog.X += 10
				}
			}
			mu.Unlock()                       //bloquea
			time.Sleep(time.Millisecond * 16) //dormir
		}
	}()

	// Goroutine se encarga de dibujr en la pantalla y manejo de mensajes
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !win.Closed() {
			win.Clear(colornames.Greenyellow)

			bgSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.8).Moved(win.Bounds().Center()))

			mu.Lock()

			frameCounter++
			if frameCounter >= 10 {
				walkingFrame = 3 - walkingFrame // Cambia entre 1 y 2
				frameCounter = 0
			}

			if dog.Vy != 0 {
				if dog.X > 100 {
					setDogSprite(&dog, spriteMap, "JumpingRight")
				} else {
					setDogSprite(&dog, spriteMap, "JumpingLeft")
				}
			} else if win.Pressed(pixelgl.KeyRight) {
				setDogSprite(&dog, spriteMap, fmt.Sprintf("Right%d", walkingFrame))
			} else if win.Pressed(pixelgl.KeyLeft) {
				setDogSprite(&dog, spriteMap, fmt.Sprintf("Left%d", walkingFrame))
			} else {
				setDogSprite(&dog, spriteMap, "NeutralRight")
			}

			for _, bone := range bones {
				boneSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.2).Moved(pixel.V(bone.X, bone.Y)))
			}

			dog.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.2).Moved(pixel.V(dog.X, dog.Y)))

			txt := text.New(pixel.V(50, 550), atlas)
			txt.Color = colornames.Black
			fmt.Fprintf(txt, "Huesos obtenidos: %d", score)

			if gameEnded {
				// Mostrar mensaje de victoria o derrota y cómo reiniciar
				bigTxt := text.New(pixel.V(150, 300), atlas)
				bigTxt.Color = colornames.Black
				bigTxt.Draw(win, pixel.IM.Scaled(bigTxt.Orig, 4))
				if gameWon {
					fmt.Fprintln(bigTxt, "¡Has ganado! Presiona R para reiniciar el juego")
				} else {
					fmt.Fprintln(bigTxt, "Has perdido, presiona R para reiniciar el juego")
				}
				bigTxt.Draw(win, pixel.IM.Scaled(bigTxt.Orig.Add(pixel.V(0, -100)), 2))
			} else {
				// Mostrar la cantidad de huesos obtenidos en la pantalla
				txt := text.New(pixel.V(150, 530), atlas)
				txt.Color = colornames.Black
				fmt.Fprintf(txt, "Come 20 huesos para ganar! ")
				fmt.Fprintf(txt, "Huesos obtenidos: %d", score)
				txt.Draw(win, pixel.IM)
			}

			// Mostrar el tiempo transcurrido dentro de la imagen de fondo
			txt.Clear() // Limpiar el contenido anterior de txt
			txt.Color = colornames.Black
			fmt.Fprintf(txt, "Tiempo restante: %d segundos", tiempoMaximo-tiempoTranscurrido)
			txt.Draw(win, pixel.IM)

			mu.Unlock()

			win.Update()
			time.Sleep(time.Millisecond * 16)
		}
	}()

	// Goroutine para el pájaro
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !win.Closed() {
			mu.Lock()

			// Cambia la posición del pájaro en la dirección actual
			bird.X += 0.03 * float64(bird.Direction)

			// Cambia la dirección cuando el pájaro llega a los límites de la pantalla
			if bird.X >= win.Bounds().Max.X || bird.X <= 0 {
				bird.Direction *= -1
			}

			// Cambia el frame del pájaro
			bird.Frame++
			if bird.Frame >= 10 {
				bird.CurrentPic = (bird.CurrentPic + 1) % len(bird.Frames)
				bird.Frame = 0
			}

			// Dibuja el pájaro en la pantalla
			bird.Frames[bird.CurrentPic].Draw(win, pixel.IM.Scaled(pixel.ZV, 0.2).Moved(pixel.V(bird.X, bird.Y)))

			mu.Unlock()

		}
	}()

	// Goroutine generacios de  huesos aleatorios
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !win.Closed() {
			mu.Lock()
			if !gameEnded && len(bones) < 10 {
				bone := Bone{X: rand.Float64()*150 + 125, Y: rand.Float64()*25 + 175}
				if rand.Float64() > 0.5 {
					bone.X = rand.Float64()*150 + 525
				}
				if rand.Float64() > 0.7 {
					bone.X = rand.Float64()*150 + 325
					bone.Y = rand.Float64()*25 + 375
				}
				bones = append(bones, bone)
			}
			mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()
}
