package main

import (
	//"bufio"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	//"strconv"
	"encoding/json"

	//r "runtime"
	ex "github.com/markus-wa/demoinfocs-golang/v3/examples"
	demoinfocs "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"

	//"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	//events "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
	//"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/msg"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	NoTeams byte = 0
	AllTeams byte = 1
	TeamTerrorists byte = 2
	TeamCounterTerrorists byte = 3
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
	controlVectors []controlVector
	boundaries []Vector
}

type Tick struct {
	players []Player
}

type Player struct {
	Name          string `json:"name"`
	Position	  Point	 `json:"position"`
	ViewDirection float32`json:"viewangle"`
	Team 		  byte   `json:"team"`
	Alive		  bool   `json:"alive"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type controlVector struct {
	vec Vector
	state byte
}

type Vector struct {
	A Point `json:"a"`
	B Point `json:"b"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

		fmt.Println("signal")

		// Get metadata for the map that the game was played on for coordinate translations
		a.demo.mapMetadata = ex.GetMapMetadata(header.MapName, msg.GetMapCrc())
		a.demo.mapName = header.MapName
		
		//load map boundaries into boundary array from json
		boundarydata, _ := os.ReadFile(header.MapName + ".json")

		_ = json.Unmarshal([]byte(boundarydata), &a.demo.boundaries)

		//for whatever reason this isnt setting the y coordinates to negative
		for i := 0; i < len(a.demo.boundaries); i++ {
			a.demo.boundaries[i].A.Y = -a.demo.boundaries[i].A.Y
			a.demo.boundaries[i].B.Y = -a.demo.boundaries[i].B.Y
		}

		fmt.Println(a.demo.boundaries)

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
			players = append(players, Player{player.Name, Point{player.Position().X, player.Position().Y}, player.ViewDirectionX(), byte(player.Team), player.IsAlive()})
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

func (a *App) GetPlayers(tick int) []Player {
	//translate coordinates to html canvas coordinates
	var players []Player

	for _, player := range a.demo.gameTicks[tick].players {
		x, y := a.demo.mapMetadata.TranslateScale(player.Position.X, player.Position.Y)
		pos := Point{x, y}
		players = append(players, Player{player.Name, pos, player.ViewDirection, player.Team, player.Alive})
	}

	//fmt.Println(a.demo.gameTicks[tick].players)
	
	return players
}

func (a *App) GetPlayerPosition(player Player) Point {

	x, y := a.demo.mapMetadata.TranslateScale(player.Position.X, player.Position.Y)

	return Point{x,y}
}

func (a *App) WriteBoundary(boundary []Vector) {

	fmt.Println(boundary)

	boundaryjson, err := json.Marshal(boundary)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = os.WriteFile(a.demo.mapName + ".json", boundaryjson, 0644)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (a *App) GetTicksOfInterest() []int {

	fmt.Println("Finding ticks")

	var ticksOfInterest []int

	for i := 0; i < len(a.demo.gameTicks); i++ {

		if i % 5000 == 0 {
			fmt.Println("Tick: ", i)
		}

		if a.CheckControlVectors(i) {
			ticksOfInterest = append(ticksOfInterest, i)
		}
	}

	fmt.Println("Done Finding Ticks")

	return ticksOfInterest
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
