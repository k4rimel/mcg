package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
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

const WIDTH = 800
const HEIGHT = 400
const friction = 0.5

var COLORS = []string{"blue", "green", "yellow", "red", "black", "purple", "pink"}

type Velocity struct {
	vX float64
	vY float64
}

type Player struct {
	Id       string
	X        float64
	Y        float64
	dX       float64
	dY       float64
	time     float64
	color    string
	velocity *Velocity
	mass     float64
}

func newPlayer() *Player {

	rb := make([]byte, 32)
	_, err := rand.Read(rb)

	if err != nil {
		fmt.Println(err)
	}
	idColor := random(0, len(COLORS)-1)
	id := base64.URLEncoding.EncodeToString(rb)

	return &Player{Id: id, X: WIDTH / 2, Y: HEIGHT / 2, dX: 1, dY: 1, time: 0.2, velocity: &Velocity{vX: 0, vY: 0}, color: COLORS[idColor], mass: 500}
}

type Game struct {
	players map[*melody.Session]*Player
}

func newGame() *Game {
	return &Game{players: make(map[*melody.Session]*Player)}
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
			for _, player := range this.players {

				xPrime := player.X + player.time*(player.dX-player.X) //+ player.velocity.vX
				yPrime := player.Y + player.time*(player.dY-player.Y) //+ player.velocity.vY

				collision := false

				for _, otherPlayer := range this.players {

					if otherPlayer == player {
						continue
					}
					xLen := xPrime - otherPlayer.X
					yLen := yPrime - otherPlayer.Y
					dist := math.Sqrt(xLen*xLen + yLen*yLen)

					if dist < 20 {
						collision = true
						otherPlayer.velocity.vX = 0.8
						otherPlayer.velocity.vY = 0.6
					}
				}

				if collision {
					continue
				}

				player.X = xPrime
				player.Y = yPrime

				if player.X > WIDTH {
					player.X = WIDTH
					player.time = 0
				}
				if player.Y > HEIGHT {
					player.Y = HEIGHT
					player.time = 0
				}

				//if player.time < 1 {
				player.time = 0.01
				//}

				if player.velocity.vX > 0 {
					player.dX += 0.01 * (-player.velocity.vX)
				}

				if player.velocity.vY > 0 {
					player.dY += 0.01 * (-player.velocity.vY)
				}

				for s, _ := range this.players {
					s.Write([]byte("player:" + player.Id + ":" + strconv.FormatFloat(player.X, 'f', 3, 64) + "," + strconv.FormatFloat(player.Y, 'f', 3, 64)))
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
		mouseCoord := getCoords(string(msg))
		p := game.GetPlayer(s)
		p.time = 0
		p.velocity.vX *= friction
		p.velocity.vY *= friction
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
func dot(x, y []int) (r int, err error) {
	if len(x) != len(y) {
		return 0, errors.New("incompatible lengths")
	}
	for i := range x {
		r += x[i] * y[i]
	}
	return
}
func collide(p1 Player, p2 Player) {
	if(true)
	return

	collision 

}
