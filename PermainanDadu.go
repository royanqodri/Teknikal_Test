package main

import (
	"fmt"
	"math/rand"
)

// Dice struct merepresentasikan dadu berenam sisi
type Dice struct {
	topSideVal int
}

// GetTopSideVal mengembalikan nilai dari sisi atas dadu
func (d *Dice) GetTopSideVal() int {
	return d.topSideVal
}

// Roll menggulung dadu dan mengatur sisi atasnya menjadi nilai acak antara 1 dan 6
func (d *Dice) Roll() *Dice {
	d.topSideVal = rand.Intn(6) + 1
	return d
}

// SetTopSideVal mengatur nilai sisi atas dadu
func (d *Dice) SetTopSideVal(topSideVal int) *Dice {
	d.topSideVal = topSideVal
	return d
}

// Player struct merepresentasikan pemain dengan cangkir dadu
type Player struct {
	diceInCup []Dice
	name      string
	position  int
	point     int
}

// GetDiceInCup mengembalikan dadu di cangkir pemain
func (p *Player) GetDiceInCup() []Dice {
	return p.diceInCup
}

// GetName mengembalikan nama pemain
func (p *Player) GetName() string {
	return p.name
}

// GetPosition mengembalikan posisi pemain
func (p *Player) GetPosition() int {
	return p.position
}

// AddPoint menambahkan poin ke pemain
func (p *Player) AddPoint(point int) {
	p.point += point
}

// GetPoint mengembalikan poin pemain
func (p *Player) GetPoint() int {
	return p.point
}

// Play menggulung dadu untuk setiap pemain di dalam cangkir
func (p *Player) Play() {
	for i := range p.diceInCup {
		p.diceInCup[i].Roll()
	}
}

// RemoveDice menghapus dadu dari cangkir pemain berdasarkan indeks
func (p *Player) RemoveDice(index int) {
	if index >= 0 && index < len(p.diceInCup) {
		p.diceInCup = append(p.diceInCup[:index], p.diceInCup[index+1:]...)
	}
}

// InsertDice menyisipkan dadu ke dalam cangkir pemain
func (p *Player) InsertDice(dice Dice) {
	p.diceInCup = append(p.diceInCup, dice)
}

// Game struct merepresentasikan permainan dengan beberapa pemain
type Game struct {
	players               []Player
	round                 int
	numberOfPlayer        int
	numberOfDicePerPlayer int
	REMOVED_WHEN_DICE_TOP int
	MOVE_WHEN_DICE_TOP    int
}

// NewGame membuat permainan baru dengan jumlah pemain dan dadu per pemain yang ditentukan
func NewGame(numberOfPlayer, numberOfDicePerPlayer int) *Game {
	game := &Game{
		round:                 0,
		numberOfPlayer:        numberOfPlayer,
		numberOfDicePerPlayer: numberOfDicePerPlayer,
		REMOVED_WHEN_DICE_TOP: 6,
		MOVE_WHEN_DICE_TOP:    1,
	}

	// Menginisialisasi pemain dan dadu mereka
	for i := 0; i < game.numberOfPlayer; i++ {
		player := Player{
			position: i,
			name:     string('A' + i),
		}
		for j := 0; j < game.numberOfDicePerPlayer; j++ {
			player.diceInCup = append(player.diceInCup, Dice{})
		}
		game.players = append(game.players, player)
	}

	return game
}

// DisplayRound menampilkan ronde saat ini dari permainan
func (g *Game) DisplayRound() *Game {
	fmt.Printf("Giliran %d\r\n", g.round)
	return g
}

// DisplayTopSideDice menampilkan nilai sisi atas dadu untuk setiap pemain
func (g *Game) DisplayTopSideDice(title string) *Game {
	fmt.Printf("%s:\n", title)
	for _, player := range g.players {
		fmt.Printf("Pemain #%s: ", player.GetName())
		var diceTopSide string
		for _, dice := range player.GetDiceInCup() {
			diceTopSide += fmt.Sprintf("%d, ", dice.GetTopSideVal())
		}
		// Tambahkan pengecekan jika panjang string diceTopSide cukup
		if len(diceTopSide) > 2 {
			fmt.Printf("%s\r\n", diceTopSide[:len(diceTopSide)-2])
		} else {
			fmt.Printf("\r\n")
		}
	}
	fmt.Printf("\r\n")
	return g
}

// DisplayWinner menampilkan pemenang permainan
func (g *Game) DisplayWinner(player Player) *Game {
	fmt.Printf("Pemenang\r\n")
	fmt.Printf("Pemain %s\r\n", player.GetName())
	return g
}

// Start memulai permainan
func (g *Game) Start() {
	fmt.Printf("Pemain = %d, Dadu = %d\r\n", g.numberOfPlayer, g.numberOfDicePerPlayer)

	// Loop sampai ada pemenang
	for {
		g.round++
		diceCarryForward := make(map[int][]Dice)

		// Menggulung dadu untuk setiap pemain
		for i := range g.players {
			g.players[i].Play()
		}

		// Menampilkan sebelum memindahkan/menghapus dadu
		g.DisplayRound().DisplayTopSideDice("Lempar Dadu")

		// Memeriksa sisi atas dadu untuk setiap pemain
		for i, player := range g.players {
			var tempDiceArray []Dice

			for j, dice := range player.GetDiceInCup() {
				// Memeriksa kemunculan angka 6
				if dice.GetTopSideVal() == g.REMOVED_WHEN_DICE_TOP {
					g.players[i].AddPoint(1)
					g.players[i].RemoveDice(j)
				}

				// Memeriksa kemunculan angka 1
				if dice.GetTopSideVal() == g.MOVE_WHEN_DICE_TOP {
					// Menentukan posisi pemain
					// Pemain maksimum berada di sebelah kanan
					// Pindahkan dadu ke sebelah kiri
					if player.GetPosition() == (g.numberOfPlayer - 1) {
						g.players[0].InsertDice(dice)
						g.players[i].RemoveDice(j)
					} else {
						tempDiceArray = append(tempDiceArray, dice)
						g.players[i].RemoveDice(j)
					}
				}
			}

			diceCarryForward[i+1] = tempDiceArray

			if val, ok := diceCarryForward[i]; ok && len(val) > 0 {
				// Menyisipkan dadu
				for _, d := range val {
					g.players[i].InsertDice(d)
				}
				// Reset
				diceCarryForward[i] = nil
			}
		}

		// Menampilkan setelah memindahkan/menghapus dadu
		g.DisplayTopSideDice("Setelah Evaluasi")

		// Mengatur jumlah pemain yang memiliki dadu
		playerHasDice := g.numberOfPlayer

		// Memeriksa apakah seorang pemain hanya memiliki satu dadu
		for _, player := range g.players {
			if len(player.GetDiceInCup()) <= 0 {
				playerHasDice--
			}
		}

		if playerHasDice == 1 {
			// Keluar dari loop jika sudah ada pemenang
			break
		}
	}

	// Menampilkan pemenang setelah keluar dari loop
	g.DisplayWinner(g.getWinner())
}

// getWinner mengembalikan pemain dengan nilai tertinggi sebagai pemenang
func (g *Game) getWinner() Player {
	var winner Player
	highscore := 0
	for _, player := range g.players {
		if player.GetPoint() > highscore {
			highscore = player.GetPoint()
			winner = player
		}
	}
	return winner
}

func main() {
	// Instance baru dari permainan
	// Tetapkan jumlah pemain dan jumlah dadu per pemain
	game := NewGame(3, 4)

	// Memulai permainan
	game.Start()
}
