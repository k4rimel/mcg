package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/olahol/melody"
	"log"
	"math"
	mathrand "math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const WIDTH = 300
const HEIGHT = 200

var COLORS = []string{"blue", "green", "yellow", "red", "black", "purple", "pink"}

type Player struct {
	Id    string
	X     float64
	Y     float64
	dX    float64
	dY    float64
	time  float64
	color string
	gX    float64
	gY    float64
}

func newPlayer() *Player {

	rb := make([]byte, 32)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}
	idColor := random(0, len(COLORS)-1)
	id := base64.URLEncoding.EncodeToString(rb)

	return &Player{Id: id, X: WIDTH / 2, Y: HEIGHT / 2, dX: 1, dY: 1, time: 0.5, gX: 0, gY: 0, color: COLORS[idColor]}
}

type Game struct {
	players map[*melody.Session]*Player
}

func newGame() *Game {
	return &Game{players: make(map[*melody.Session]*Player)}
}

func (this *Game) AddPlayer(s *melody.Session) {
	p := newPlayer()

	for z, n := range this.players {
		z.Write([]byte("add:" + p.Id + ":" + strconv.FormatFloat(p.X, 'f', 3, 64) + "," + strconv.FormatFloat(p.Y, 'f', 3, 64) + "," + p.color))
		s.Write([]byte("add:" + n.Id + ":" + strconv.FormatFloat(n.X, 'f', 3, 64) + "," + strconv.FormatFloat(n.Y, 'f', 3, 64) + "," + n.color))
	}

	s.Write([]byte("add:" + p.Id + ":" + strconv.FormatFloat(p.X, 'f', 3, 64) + "," + strconv.FormatFloat(p.Y, 'f', 3, 64) + "," + p.color))
	this.players[s] = p

}
func (this *Game) GetPlayer(s *melody.Session) *Player {
	return this.players[s]
}

func (this *Game) RemovePlayer(s *melody.Session) {

	id := this.players[s].Id
	delete(this.players, s)

	for z, _ := range this.players {
		z.Write([]byte("remove:" + id))
	}
}

func (this *Game) run() {
	ticker := time.NewTicker(time.Millisecond * 1)
	go func() {
		for {
			<-ticker.C
			for _, p := range this.players {

				//p.X += p.time * (p.dX - p.X)
				//p.Y += p.time * (p.dY - p.Y)

				tempx := p.X + p.time*(p.dX-p.X) + p.gX
				tempy := p.Y + p.time*(p.dY-p.Y) + p.gY

				collision := false

				for _, fp := range this.players {

					if fp == p {
						continue
					}
					cdx := tempx - fp.X
					cdy := tempy - fp.Y
					dist := math.Sqrt(cdx*cdx + cdy*cdy)

					if dist < 20 {
						collision = true
						fp.gX = 0.8
						fp.gY = 0.6
					}
				}

				if collision {
					continue
				}

				p.X = tempx
				p.Y = tempy

				if p.X > WIDTH {
					p.X = WIDTH
					p.time = 0
				}
				if p.Y > HEIGHT {
					p.Y = HEIGHT
					p.time = 0
				}

				//if p.time < 1 {
				p.time = 0.01
				//}

				if p.gX > 0 {
					p.gX -= 0.01
				}

				if p.gY > 0 {
					p.gY -= 0.01
				}

				for s, _ := range this.players {
					// fmt.Println("player:" + p.Id + ":" + strconv.FormatInt(p.X, 10) + "," + strconv.FormatInt(p.Y, 10))

					s.Write([]byte("player:" + p.Id + ":" + strconv.FormatFloat(p.X, 'f', 3, 64) + "," + strconv.FormatFloat(p.Y, 'f', 3, 64)))
				}

			}

		}
	}()
}

func main() {
	r := gin.New()
	m := melody.New()

	size := 65536
	m.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  size,
		WriteBufferSize: size,
	}
	m.Config.MaxMessageSize = int64(size)
	m.Config.MessageBufferSize = 2048

	game := newGame()

	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		fmt.Println("Debut de connection s")
		m.HandleRequest(c.Writer, c.Request)
	})

	//var mutex sync.Mutex

	m.HandleConnect(func(s *melody.Session) {
		game.AddPlayer(s)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		//fmt.Println(fmt.Sprintf("message %s", string(msg)))
		mouseCoord := getCoords(string(msg))
		p := game.GetPlayer(s)
		p.time = 0
		p.dX = mouseCoord.X
		p.dY = mouseCoord.Y

		//fmt.Println(fmt.Sprintf("time %f", p.time))

		/*
			if deltaX := p.X - coord.X; deltaX > 0 {
				p.dX = -1
			} else if deltaX := p.X - coord.X; deltaX < 0 {
				p.dX = 1
			} else {
				p.dX = 0
			}
			if deltaY := p.Y - coord.Y; deltaY > 0 {
				p.dY = -1
			} else if deltaY := p.Y - coord.Y; deltaY < 0 {
				p.dY = 1
			} else {
				p.dY = 0
			}
		*/

	})

	m.HandleDisconnect(func(s *melody.Session) {
		game.RemovePlayer(s)
	})

	game.run()

	r.Run(":5000")
}

type coordinates struct {
	X float64
	Y float64
}

func getCoords(s string) *coordinates {
	var p coordinates
	dec := json.NewDecoder(strings.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(p.X, p.Y)
	return &p
}

func random(min, max int) int {
	mathrand.Seed(time.Now().Unix())
	return mathrand.Intn(max-min) + min
}
