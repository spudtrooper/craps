package craps

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

type outcome string

const (
	craps outcome = "craps"
	come  outcome = "come"
)

type state string

const (
	off       state = "off"
	comingOut state = "coming out"
	on        state = "on"
)

type Dice struct {
	left, right int
}

func (d *Dice) Val() int { return d.left + d.right }

type Game struct {
	point   int
	player  Player
	dice    Dice
	history []Game
	outcome outcome
	state   state
}

func (g Game) Point() int { return g.point }

func (g Game) SetPlayer(name string, money Money, s Strategy) Game {
	p := Player{name: name, strat: s, money: []Money{money}}
	g.player = p
	return g
}

func (g Game) String() string {
	return fmt.Sprintf("Game{point: %d, dice: %d, players: %+v}", g.point, g.dice.Val(), g.player)
}

func rollDice() Dice {
	left, right := rand.Intn(6)+1, rand.Intn(6)+1
	return Dice{left: left, right: right}
}

func (g Game) Roll() Game {
	dice := rollDice()
	log("roll: %d", dice.Val())
	var player Player
	point := g.point
	var outcome outcome
	var state state
	if g.isOn() {
		if g.point == dice.Val() {
			log("on: %s", color.GreenString("point"))
			player = g.pointWhenOn(dice)
			point = 0
			state = off
		} else {
			state = on
			switch dice.Val() {
			case 7:
				log("on: %s", color.RedString("craps lose"))
				outcome = craps
				player = g.crapsWhenOnLose()
				point = 0
				state = off
			case 2, 3, 11, 12:
				log("on: %s", color.YellowString("craps win"))
				outcome = craps
				player = g.crapsWhenOnWin(dice)
			default:
				log("on: %s", color.GreenString("come"))
				outcome = come
				player = g.comeWhenOn(dice)
			}
		}
	} else { // off
		switch dice.Val() {
		case 7, 11:
			log("off: %s", color.GreenString("craps win"))
			outcome = craps
			player = g.crapsWhenOffWin()
			state = off
		case 2, 3, 12:
			log("off: %s", color.RedString("craps lose"))
			outcome = craps
			player = g.crapsWhenOffLose()
			state = off
		default:
			log("off: %s", color.BlueString("set point"))
			outcome = come
			player = g.player.AddMoney(0)
			point = dice.Val()
			state = comingOut
		}
	}

	res := Game{
		point:   point,
		dice:    dice,
		history: append(g.history, g),
		player:  player,
		outcome: outcome,
		state:   state,
	}

	log("next game: %+v", res)

	return res
}

func (g Game) pointWhenOn(dice Dice) Player {
	p := g.player
	bet := p.Bet(g)
	var money Money
	if bet.Pass > 0 {
		money += bet.Pass
		switch dice.Val() {
		case 4, 10:
			money += 9 * bet.PassOdds / 5
		case 5, 9:
			money += 7 * bet.PassOdds / 5
		case 6, 8:
			money += 7 * bet.PassOdds / 6
		}
	}
	return p.AddMoney(money)
}

func (g Game) crapsWhenOnLose() Player {
	p := g.player
	bet := p.Bet(g)
	var money Money
	money -= bet.Pass
	money -= bet.PassOdds
	money += bet.Come
	money -= bet.DontPass
	money += bet.DontCome
	money -= bet.Place4
	money -= bet.Place5
	money -= bet.Place6
	money -= bet.Place7
	money -= bet.Place8
	money -= bet.Place9
	money -= bet.Place10
	money -= bet.Place4Odds
	money -= bet.Place5Odds
	money -= bet.Place6Odds
	money -= bet.Place7Odds
	money -= bet.Place8Odds
	money -= bet.Place9Odds
	money -= bet.Place10Odds
	return p.AddMoney(money)
}

func (g Game) crapsWhenOnWin(dice Dice) Player {
	p := g.player
	bet := p.Bet(g)
	var money Money
	// TODO
	switch dice.Val() {
	case 7, 11:
		money += bet.Come
		money -= bet.DontCome
	case 2, 3, 12:
		money -= bet.Come
		money += bet.DontCome
	}
	return p.AddMoney(money)
}

func (g Game) comeWhenOn(dice Dice) Player {
	p := g.player
	bet := p.Bet(g)
	var money Money
	switch dice.Val() {
	case 4:
		money += bet.Place4
		money += 9 * bet.Place6Odds / 5
	case 5:
		money += bet.Place5
		money += 7 * bet.Place5Odds / 5
	case 6:
		money += bet.Place6
		money += 7 * bet.Place6Odds / 6
	case 8:
		money += bet.Place8
		money += 7 * bet.Place8Odds / 6
	case 9:
		money += bet.Place9
		money += 7 * bet.Place9Odds / 5
	case 10:
		money += bet.Place10
		money += 9 * bet.Place10Odds / 5
	}
	return p.AddMoney(money)
}

func (g Game) crapsWhenOffWin() Player {
	p := g.player
	bet := p.Bet(g)
	var money Money
	money += bet.Pass
	money -= bet.DontPass
	return p.AddMoney(money)
}

func (g Game) crapsWhenOffLose() Player {
	p := g.player
	bet := p.Bet(g)
	var money Money
	money -= bet.Pass
	money += bet.DontPass
	return p.AddMoney(money)
}

func (g Game) Player() Player { return g.player }

func (g Game) PrintResults() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)

	sign := ""
	none := tablewriter.Colors{}
	totalCol := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor}
	if g.player.Final() > g.player.Start() {
		totalCol = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}
		sign = "+"
	} else if g.player.Final() < g.player.Start() {
		totalCol = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
		sign = "-"
	}
	table.SetHeader([]string{"#", "State", "Point", "Roll", "On", "Outcome", "Diff $", "Total $"})
	table.SetFooter([]string{"", "", "", "", "", "", "",
		fmt.Sprintf("$%0.2f (Δ: %s%0.2f/%d%%)", g.player.Final(), sign, g.player.Start()-g.player.Final(),
			int(100*(g.player.Start()-g.player.Final())/g.player.Start()))})
	table.SetFooterColor(none, none, none, none, none, none, none, totalCol)
	h := tablewriter.Colors{tablewriter.Bold, tablewriter.BgWhiteColor, tablewriter.FgBlackColor}
	table.SetHeaderColor(h, h, h, h, h, h, h, h)

	for i, g := range g.history {
		if i == 0 {
			continue
		}
		p := g.player
		total := p.Final()
		var diff Money
		if i > 0 {
			before := p.money[len(p.money)-2]
			diff = total - before
		}
		outcomeStr := string(g.outcome)
		outcomeCol := tablewriter.Colors{}
		point := g.history[i-1].point
		onStr := ""
		var onCol tablewriter.Colors
		var addBreak bool
		if point != 0 {
			if point == g.dice.Val() {
				outcomeStr = "hit point"
				outcomeCol = tablewriter.Colors{tablewriter.Bold, tablewriter.BgHiGreenColor, tablewriter.FgBlackColor}
				onCol = tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgBlackColor, tablewriter.Bold}
				onStr = "-"
				addBreak = true
			} else {
				switch g.dice.Val() {
				case 7:
					outcomeStr = "crap out"
					outcomeCol = tablewriter.Colors{tablewriter.Bold, tablewriter.BgRedColor, tablewriter.FgBlackColor}
					onCol = tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgBlackColor, tablewriter.Bold}
					onStr = "-"
					addBreak = true
				case 2, 3, 11, 12:
					outcomeCol = tablewriter.Colors{tablewriter.BgYellowColor, tablewriter.FgBlackColor}
				default:
					outcomeStr = "come"
					outcomeCol = tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor, tablewriter.FgBlackColor}
				}
			}
		} else { // off
			switch g.dice.Val() {
			case 7, 11:
				outcomeStr = "craps"
				outcomeCol = tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor, tablewriter.FgBlackColor}
			case 2, 3, 12:
				outcomeStr = "craps"
				outcomeCol = tablewriter.Colors{tablewriter.Bold, tablewriter.BgHiRedColor, tablewriter.FgBlackColor}
			default:
				outcomeStr = "set point"
				outcomeCol = tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlueColor, tablewriter.FgBlackColor}
			}
		}
		var stateCol tablewriter.Colors
		switch g.state {
		case on:
			stateCol = tablewriter.Colors{tablewriter.Bold, tablewriter.BgYellowColor, tablewriter.FgBlackColor}
			onStr = "*"
			onCol = tablewriter.Colors{tablewriter.BgBlackColor, tablewriter.FgWhiteColor, tablewriter.Bold}
		case off:
			stateCol = tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlackColor, tablewriter.FgYellowColor}
		case comingOut:
			stateCol = tablewriter.Colors{tablewriter.Bold, tablewriter.BgCyanColor, tablewriter.FgBlackColor}
			onCol = tablewriter.Colors{tablewriter.BgGreenColor, tablewriter.FgBlackColor, tablewriter.Bold}
			onStr = "+"
		}
		diffStr := "--"
		if diff != 0 {
			diffStr = fmt.Sprintf("%+0.2f", diff)
		}
		data := []string{
			strconv.Itoa(i),
			string(g.state),
			strconv.Itoa(point),
			strconv.Itoa(g.dice.Val()),
			onStr,
			outcomeStr,
			fmt.Sprintf("%+10s", diffStr),
			fmt.Sprintf("%+15s", fmt.Sprintf("%0.2f", total)),
		}
		diffCol := tablewriter.Colors{tablewriter.FgWhiteColor}
		if diff < 0 {
			diffCol = tablewriter.Colors{tablewriter.FgRedColor}
		} else if diff > 0 {
			diffCol = tablewriter.Colors{tablewriter.FgGreenColor}
		}
		colors := []tablewriter.Colors{
			{},
			stateCol,
			{},
			{},
			onCol,
			outcomeCol,
			diffCol,
			{tablewriter.Bold, tablewriter.FgWhiteColor},
		}
		table.Rich(data, colors)
		if addBreak {
			table.Append([]string{"..", "", "", "", "", "", "", ""})
		}
	}
	table.Render()
}

func (g Game) isOn() bool { return g.point != 0 }

type Strategy interface {
	Initial(g Game, p Player) Bet
}

type Player struct {
	name  string
	money []Money
	strat Strategy
}

func (p Player) String() string {
	return fmt.Sprintf("Player{name: %s, money: %0.2f}", p.name, p.Final())
}

func (p Player) Final() Money {
	if len(p.money) == 0 {
		return 0
	}
	return p.money[len(p.money)-1]
}

func (p Player) Start() Money {
	return p.money[0]
}

func (p Player) AddMoney(money Money) Player {
	before := p.Final()
	after := before
	after += money
	p.money = append(p.money, after)
	m := color.YellowString("$%0.2f", money)
	a := color.YellowString("$%0.2f", after)
	sign := "_"
	if after > before {
		m = color.HiGreenString("$%0.2f", money)
		a = color.GreenString("$%0.2f", after)
		sign = "+"
	} else if after < before {
		m = color.HiRedString("$%0.2f", money)
		a = color.RedString("$%0.2f", after)
		sign = "-"
	}
	log("Adding to %s: $%0.2f %s %s -> %s", p.name, before, sign, m, a)
	return p
}

func (p Player) Bet(g Game) Bet {
	bet := p.strat.Initial(g, p)
	log("Bet for %s: %s", p.name, bet)
	return bet
}

type Money float64

type Bet struct {
	Pass, DontPass                                                                      Money
	PassOdds                                                                            Money
	Field                                                                               Money
	Place4, Place5, Place6, Place7, Place8, Place9, Place10                             Money
	Place4Odds, Place5Odds, Place6Odds, Place7Odds, Place8Odds, Place9Odds, Place10Odds Money
	Come, DontCome                                                                      Money
}

func moneyFieldsReflectively(x interface{}) string {
	v := reflect.ValueOf(x)
	var buf bytes.Buffer
	for _, sf := range reflect.VisibleFields(v.Type()) {
		f := v.FieldByName(sf.Name)
		if f.Type().String() != "craps.Money" {
			continue
		}
		m := f.Interface().(Money)
		if m == 0 {
			continue
		}
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(sf.Name)
		buf.WriteString(":")
		buf.WriteString(fmt.Sprintf("%0.2f", m))
	}
	return buf.String()
}

func (b Bet) String() string {
	return moneyFieldsReflectively(b)
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
