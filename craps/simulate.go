package craps

import (
	"flag"
	"fmt"
	"math"
	"sort"
)

type GameCtr func() Game

func median(data []float64) float64 {
	dataCopy := make([]float64, len(data))
	copy(dataCopy, data)

	sort.Float64s(dataCopy)

	var median float64
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = (dataCopy[l/2-1] + dataCopy[l/2]) / 2
	} else {
		median = dataCopy[l/2]
	}

	return median
}

func mean(data []float64) float64 {
	var total float64
	for _, v := range data {
		total += v
	}
	return math.Round(total / float64(len(data)))
}

func min(data []float64) float64 {
	res := data[0]
	for _, v := range data[1:] {
		res = math.Min(res, v)
	}
	return res
}

func max(data []float64) float64 {
	res := data[0]
	for _, v := range data[1:] {
		res = math.Max(res, v)
	}
	return res
}

func wins(g Game, data []float64) int {
	var res int
	for _, v := range data {
		if v > float64(g.player.Start()) {
			res++
		}
	}
	return res
}

func losses(g Game, data []float64) int {
	var res int
	for _, v := range data {
		if v < float64(g.player.Start()) {
			res++
		}
	}
	return res
}

func ties(g Game, data []float64) int {
	var res int
	for _, v := range data {
		if v == float64(g.player.Start()) {
			res++
		}
	}
	return res
}

func runGame(g Game, games, rolls int) Money {
	for i := 0; i <= rolls; i++ {
		g = g.Roll()
	}
	if games <= 1 {
		fmt.Println()
		fmt.Println("Printing summary of 1 game...")
		g.PrintResults()
	}
	return g.Player().Final()
}

func SimulateGames(gameCtr GameCtr, games, rolls int) {
	flag.Parse()

	g := gameCtr()

	var money []float64
	for i := 0; i < games; i++ {
		m := runGame(g, games, rolls)
		money = append(money, float64(m))
	}

	if games > 1 {
		fmt.Println()
		fmt.Printf("Printing stats of %d games...\n", games)
		mean := mean(money)
		median := median(money)
		max := max(money)
		min := min(money)
		wins := wins(g, money)
		losses := losses(g, money)
		ties := ties(g, money)
		fmt.Printf("COUNT  : %d\n", len(money))
		fmt.Printf("MEDIAN : %.1f\n", median)
		fmt.Printf("MEAN   : %.1f\n", mean)
		fmt.Printf("MIN    : %.1f\n", min)
		fmt.Printf("MAX    : %.1f\n", max)
		fmt.Printf("WINS   : %d (%.2f%%)\n", wins, 100*float64(wins)/float64(games))
		fmt.Printf("LOSSES : %d (%.2f%%)\n", losses, 100*float64(losses)/float64(games))
		fmt.Printf("TIES   : %d (%.2f%%)\n", ties, 100*float64(ties)/float64(games))
	}
}
