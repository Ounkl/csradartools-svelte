package main

import (
	"bufio"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"

	//r "runtime"
	ex "github.com/markus-wa/demoinfocs-golang/v3/examples"
	demoinfocs "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"

	//"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	//events "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
	//"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/msg"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx  context.Context
	demo Replay
}

type Replay struct {
	gameTicks   []Tick
	mapMetadata ex.Map
	mapName     string
}

type Tick struct {
	players []Player
}

type Player struct {
	Name          string
	X             float64
	Y             float64
	ViewDirection float32
}

type Point struct {
	X float32
	Y float32
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	fmt.Println(intersectPoint(Point{113.46939, -678.7755}, Point{98.47115, -678.5454}, Point{64, -742}, Point{72, -654}))
}

// Accesses provided demo through path, updates radar image based on metadata
func (a *App) GetDemo() {

	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Demo files (*.dem)",
				Pattern:     "*.dem",
			},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(path)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	parser := demoinfocs.NewParser(f)

	// Parse header (contains map-name etc.)
	header, err := parser.ParseHeader()
	if err != nil {
		panic(err)
	}

	var (
		mapRadarImg image.Image
	)

	parser.RegisterNetMessageHandler(func(msg *msg.CSVCMsg_ServerInfo) {

		fmt.Printf("signal")

		// Get metadata for the map that the game was played on for coordinate translations
		a.demo.mapMetadata = ex.GetMapMetadata(header.MapName, msg.GetMapCrc())

		a.demo.mapName = header.MapName

		//Load map overview image
		mapRadarImg = ex.GetMapRadar(header.MapName, msg.GetMapCrc())

		imgFile, err := os.Create("frontend/src/assets/images/radar.png")
		defer imgFile.Close()
		if err != nil {
			fmt.Println("Cannot create file:", err)
		}
		png.Encode(imgFile, mapRadarImg)
	})

	a.parseMatch(parser)
}

func (a *App) parseMatch(dem demoinfocs.Parser) {

	fmt.Println("Parsing...")

	readable, err := dem.ParseNextFrame()

	for readable && err == nil {

		var players []Player

		for _, player := range dem.GameState().Participants().Playing() {
			players = append(players, Player{player.Name, player.Position().X, player.Position().Y, player.ViewDirectionX()})
		}

		a.demo.gameTicks = append(a.demo.gameTicks, Tick{players})

		//fmt.Println(len(dem.GameState().Participants().All()))

		readable, err = dem.ParseNextFrame()
	}

	fmt.Println("Parsing Completed")

	fmt.Println(len(a.demo.gameTicks))

	if err != nil {
		log.Println(err)
	}
}

func (a *App) GetPlayerPositions(tick int) []float64 {

	var players []float64

	//fmt.Println(len(a.demo.gameTicks[tick].state.Participants().All()))
	//fmt.Println(tick)

	for _, player := range a.demo.gameTicks[tick].players {

		//fmt.Println("Getting Positions")

		x, y := a.demo.mapMetadata.TranslateScale(player.X, player.Y)

		players = append(players, x)
		players = append(players, y)
	}

	return players
}

func (a *App) WriteBoundary(boundary []int) {

	f, err := os.Create(a.demo.mapName + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(a.demo.mapName + "\n")

	if err2 != nil {
		log.Fatal(err2)
	}

	for _, coord := range boundary {

		_, err3 := f.WriteString(strconv.Itoa(coord) + "\n")

		if err3 != nil {
			log.Fatal(err2)
		}
	}

	fmt.Println("Wrote.")
}

func (a *App) GetViewDirections(tick int) []float32 {

	var directionRays []float32

	for _, player := range a.demo.gameTicks[tick].players {

		legb := 15 * math.Cos(float64(player.ViewDirection) * (math.Pi / 180))
		//y axis is inverted in the html canvas, so the coordinate needs to be inverted
		lega := 15 * math.Sin(float64(player.ViewDirection) * (math.Pi / 180))


		/**

		if player.ViewDirection > 180 {
			lega = -lega
		}

		if player.ViewDirection > 90 && player.ViewDirection < 270 {
			legb = -legb
		}
		*/

		fmt.Println(player.ViewDirection, player.X, player.Y)

		x, y := a.demo.mapMetadata.TranslateScale(player.X, player.Y)

		directionRays = append(directionRays, float32(x), float32(y), float32(x+legb), float32(y+lega))

	}

	return directionRays
}

func (a *App) GetRayCast(tick int) []float32 {

	var boundary []Point

	//load wall vectors into memory
	file, err := os.Open("de_inferno.txt")

	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()

	var PointX int
	var PointY int

	//gets the coordinates of the vectors, uses modulo to record every two values
	for scanner.Scan() {

		PointX, _ = strconv.Atoi(scanner.Text())

		scanner.Scan()

		PointY, _ = strconv.Atoi(scanner.Text())

		boundary = append(boundary, Point{float32(PointX), -float32(PointY)})

	}

	var vectorRays []float32

	for _, player := range a.demo.gameTicks[tick].players {
		//b is arbitrary (i think it can be any value honestly, depends on how intersectPoint's math works
		//leg := 50 * math.Tan(float64(player.ViewDirection) * (math.Pi / 180))

		legb := 15 * math.Cos(float64(player.ViewDirection) * (math.Pi / 180))
		//y axis is inverted in the html canvas, so the coordinate needs to be inverted
		lega := 15 * math.Sin(float64(player.ViewDirection) * (math.Pi / 180))

		playerX, playerY := a.demo.mapMetadata.TranslateScale(player.X, player.Y)

		playerVectorPointA := Point{float32(playerX), -float32(playerY)}

		playerVectorPointB := Point{float32(playerX) + float32(legb), -float32(playerY) + float32(lega)}

		//fmt.Println(playerVectorPointA, playerVectorPointB)
		//fmt.Println(player.ViewDirection)

		var smallestDistance float32
		var smallestX float32
		var smallestY float32
		initialFound := false
		//compare the player vector point a/b with every wall vector and pick the shortest distance as the point to return
		for i := 0; i < len(boundary); i = i + 2 {
			dist, intersects := intersectPoint(playerVectorPointA, playerVectorPointB, boundary[i], boundary[i+1])

			//fmt.Println(dist, intersects)

			if intersects {

				distanceX := dist.X - playerVectorPointA.X
				distanceY := dist.Y - playerVectorPointA.Y

				//pythagorean to find dist
				distance := math.Sqrt(float64(distanceX*distanceX + distanceY*distanceY))
				if !initialFound || float32(distance) < smallestDistance {
					smallestDistance = float32(distance)
					smallestX = dist.X
					smallestY = dist.Y
					initialFound = true
				}
			}
		}

		if initialFound {
			//endPointY := smallestX * float32(math.Tan(float64(player.ViewDirection) * (math.Pi / 180)))
			//fmt.Println(smallestX, endPointY, player.ViewDirection, math.Tan(float64(player.ViewDirection) * (math.Pi / 180)))

			//fmt.Println(playerVectorPointA.X, " ", -playerVectorPointA.Y, " ", smallestX, " ", -smallestY)
			vectorRays = append(vectorRays, playerVectorPointA.X, -playerVectorPointA.Y, smallestX, -smallestY)
		}
	}
	return vectorRays
}

func intersectPoint(a Point, b Point, c Point, d Point) (Point, bool) {

	//fmt.Println(a)
	//fmt.Println(b)
	//fmt.Println(c)
	//fmt.Println(d)


	denom := (c.X-d.X)*(a.Y-b.Y) - (c.Y-d.Y)*(a.X-b.X)
	//(64 - 72 * -678 - -678) - (-742 - -654 * 113 - 98)
	// = 0 - - 1320
	// = 1320

	if denom == 0 {
		return Point{}, false
	}

	t := ((c.X-a.X)*(a.Y-b.Y) - (c.Y-a.Y)*(a.X-b.X)) / denom
	// (64 - 113 * -678 - -678) - (-742 - -678 * 113 - 98) / 1320
	//  = 0 - -960 / 1320
	// = 960 / 1320
	//

	u := ((c.X-a.X)*(c.Y-d.Y) - (c.Y-a.Y)*(c.X-d.X)) / denom
	//- (64 - 72 * -742 - -678) - (-742 - -678 * 64 - 113) / denom
	// = - (-512 - 3136) / denom
	// = 3648 / 1320
	// = > 0

	if t > 0 && t < 1 && u > 0 {

		//fmt.Println(t, " ", u)

		return Point{c.X + t*(d.X-c.X), c.Y + t*(d.Y-c.Y)}, true
		//64 + 0.72 * 72 - 64, -742 + 0.72 * -654 - -742
		//69.76, -678.64

	} else {
		return Point{}, false
	}
}

func (a *App) GetTickCount() int {
	return len(a.demo.gameTicks)
}

func (a *App) runDemo() string {
	return "hello"
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
