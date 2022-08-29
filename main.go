package main

import (
	"flag"

	"github.com/spudtrooper/craps/craps"
)

var (
	start       = flag.Float64("start", 1000, "starting money")
	rolls       = flag.Int("rolls", 10, "number of rolls")
	games       = flag.Int("games", 10, "number of games")
	passBet     = flag.Int("pass", 25, "pass bet")
	place6Bet   = flag.Float64("place6", 60, "place6 bet")
	place8Bet   = flag.Float64("place8", 60, "place8 bet")
	betPassOdds = flag.Bool("pass_odds", false, "bet place odds")
)

type strat struct{}

func (*strat) Initial(g craps.Game, p craps.Player) craps.Bet {
	var passOdds int
	pass := *passBet
	if *betPassOdds {
		switch g.Point() {
		case 4, 10:
			passOdds = int(3 * pass)
		case 5, 9:
			passOdds = int(4 * pass)
		case 6, 8:
			passOdds = int(5 * pass)
		}
	}
	return craps.Bet{
		Pass:     craps.Money(pass),
		PassOdds: craps.Money(passOdds),
		Place6:   craps.Money(*place6Bet),
		Place8:   craps.Money(*place8Bet),
	}
}

func main() {
	flag.Parse()

	craps.SimulateGames(func() craps.Game {
		var g craps.Game
		g = g.SetPlayer("Jeff", craps.Money(*start), &strat{})
		return g
	}, *games, *rolls)
}
