package main

import (
	"fmt"
	"math"
	"sort"
)

const (
	MAX_ORDERS      = 5
	MAX_INGREDIENTS = 10
	MAX_BREWED      = 3
)

type Point int

type Witch struct {
	Ingredients
	Points Point
	Me     bool
	Brewed int
}

type Ingredients [4]int

func NewIngredients(i0, i1, i2, i3 int) Ingredients {
	return [4]int{i0, i1, i2, i3}
}
func (igts Ingredients) Add(i0, i1, i2, i3 int) Ingredients {
	return NewIngredients(
		igts[0]+i0,
		igts[1]+i1,
		igts[2]+i2,
		igts[3]+i3,
	)
}
func (igts Ingredients) AddIngredients(other Ingredients) Ingredients {
	return igts.Add(other[0], other[1], other[2], other[3])
}
func (igts Ingredients) IsLegit() bool {
	sum := 0
	for _, q := range igts {
		if q < 0 {
			return false
		}
		sum += q
	}
	return sum <= MAX_INGREDIENTS
}
func (igts Ingredients) Balance() float64 {
	const balance = MAX_INGREDIENTS / 4.0
	b0 := float64(igts[0]) - balance
	b1 := float64(igts[1]) - balance
	b2 := float64(igts[2]) - balance
	b3 := float64(igts[3]) - balance
	return math.Sqrt(b0*b0 + b1*b1 + b2*b2 + b3*b3)
}
func (igts Ingredients) IsMoreBalancedThan(other Ingredients) bool {
	return igts.Balance() > other.Balance()
}
func (igts Ingredients) Sum() (sum int) {
	for _, q := range igts {
		sum += q
	}
	return sum
}
func (igts Ingredients) Complexity() (sum int) {
	for i, q := range igts {
		sum += q * (i + 1)
	}
	return sum
}

type ActionType string

const (
	BREW          = ActionType("BREW")
	CAST          = ActionType("CAST")
	OPPONENT_CAST = ActionType("OPPONENT_CAST")
	LEARN         = ActionType("LEARN")
	REST          = ActionType("REST")
)

type Action struct {
	ID int
	Ingredients
	Points     Point
	Type       ActionType
	Castable   bool
	Repeatable bool
	TaxCount   int // the amount of taxed tier-0 ingredients you gain from learning this spell
	TomeIndex  int // the index in the tome if this is a tome spell, equal to the above tax
}

func (a Action) String() string {
	switch a.Type {
	case BREW, CAST:
		return fmt.Sprintf("%s %d", a.Type, a.ID)
	case REST:
		return string(REST)
	default:
		return string(a.Type)
	}
}

func (a Action) IsPossibleFor(witch Witch) bool {
	switch a.Type {
	case BREW:
		return witch.AddIngredients(a.Ingredients).IsLegit()
	case CAST:
		return witch.Me && a.Castable && witch.AddIngredients(a.Ingredients).IsLegit()
	case OPPONENT_CAST:
		return !witch.Me && a.Castable && witch.AddIngredients(a.Ingredients).IsLegit()
	case REST:
		return true
	default:
		return false
	}
}

func (a Action) IsNeededFor(witch Witch, target Action) bool {
	switch target.Type {
	case BREW:
		switch a.Type {
		case CAST, OPPONENT_CAST:
			debug("cast: ", a.Ingredients)
			for i, q := range witch.Ingredients.AddIngredients(target.Ingredients) {
				if q < 0 && a.Ingredients[i] > 0 {
					debug("missing ingredient index ", q, " * ", i)
					return true
				}
			}
			return false
		default:
			return false
		}
	default:
		return true
	}
}

// IsLessThan return true if a is less considered than other
func (a Action) IsLessThan(other Action) bool {
	return other.Type != REST && (false ||
		a.Type != other.Type && a.Type != BREW ||
		a.Points < other.Points ||
		!a.IsMoreBalancedThan(other.Ingredients))
}

type ActionLesser func(ai, aj Action) bool // return true if ai is less considered than aj

type ActionSlice struct {
	Slice  []Action
	Lesser ActionLesser
}

func (as ActionSlice) Len() int { return len(as.Slice) }
func (as ActionSlice) Less(i, j int) bool {
	ai, aj := as.Slice[i], as.Slice[j]
	if as.Lesser == nil {
		return ai.IsLessThan(aj)
	}
	return as.Lesser(ai, aj)
}
func (as ActionSlice) Swap(i, j int) { as.Slice[i], as.Slice[j] = as.Slice[j], as.Slice[i] }

func (as ActionSlice) Last() Action { return as.Slice[len(as.Slice)-1] }
func (as ActionSlice) Pick(optional ...ActionLesser) Action {
	if len(optional) == 1 && optional[0] != nil {
		as.Lesser = optional[0]
	}
	sort.Sort(as)
	debug(as)
	if as.Len() == 0 {
		return Action{Type: REST}
	}
	return as.Last()
}

func main() {

	LastTurnME := Witch{Me: true}
	LastTurnOPNT := Witch{}
	for {
		Actions := ActionSlice{Slice: []Action{{Type: REST}}}

		var actionCount int // actionCount: the number of spells and recipes in play
		fmt.Scan(&actionCount)
		for i := 0; i < actionCount; i++ {
			var delta0, delta1, delta2, delta3, _castable, _repeatable int
			action := Action{}
			fmt.Scan(&action.ID, &action.Type, &delta0, &delta1, &delta2, &delta3, &action.Points, &action.TomeIndex, &action.TaxCount, &_castable, &_repeatable)
			action.Castable, action.Repeatable = _castable != 0, _repeatable != 0
			action.Ingredients = NewIngredients(delta0, delta1, delta2, delta3)

			switch action.Type {
			case BREW, CAST, OPPONENT_CAST:
				Actions.Slice = append(Actions.Slice, action)
			case LEARN:
				fallthrough
			default:
				panic(fmt.Errorf("unknown action type %s\n%+v", action.Type, action))
			}
		}

		var inv0, inv1, inv2, inv3 int

		ME := Witch{Me: true, Brewed: LastTurnME.Brewed}
		fmt.Scan(&inv0, &inv1, &inv2, &inv3, &ME.Points)
		ME.Ingredients = NewIngredients(inv0, inv1, inv2, inv3)
		if ME.Points > LastTurnME.Points {
			ME.Brewed++
		}
		LastTurnME = ME

		OPNT := Witch{Brewed: LastTurnOPNT.Brewed}
		fmt.Scan(&inv0, &inv1, &inv2, &inv3, &OPNT.Points)
		OPNT.Ingredients = NewIngredients(inv0, inv1, inv2, inv3)
		if OPNT.Points > LastTurnOPNT.Points {
			OPNT.Brewed++
		}
		LastTurnOPNT = OPNT

		Orders := ActionSlice{}
		for _, a := range Actions.Slice {
			if a.Type == BREW {
				Orders.Slice = append(Orders.Slice, a)
			}
		}

		targetPotion := Orders.Pick(EasierPotionFor(ME))
		debug("target potion: ", targetPotion)

		NextTurnPossibilities := ActionSlice{}
		incaseWeAreStuck := ActionSlice{}
		for _, a := range Actions.Slice {
			if a.IsPossibleFor(ME) {
				incaseWeAreStuck.Slice = append(incaseWeAreStuck.Slice, a)
				if a.Castable && !a.IsNeededFor(ME, targetPotion) {
					continue
				}
				NextTurnPossibilities.Slice = append(NextTurnPossibilities.Slice, a)
			}
		}
		if NextTurnPossibilities.Len() <= 1 {
			weAreStuck := true
			for _, a := range Actions.Slice {
				if a.Type == CAST && !a.Castable {
					weAreStuck = false
					break
				}
			}
			if weAreStuck {
				NextTurnPossibilities = incaseWeAreStuck
			}
		}

		nextAction := NextTurnPossibilities.Pick(func(ai, aj Action) bool {
			return ai.IsLessThan(aj)
		})
		// in the first league: BREW <id> | WAIT; later: BREW <id> | CAST <id> [<times>] | LEARN <id> | REST | WAIT
		fmt.Println(
			nextAction,
			fmt.Sprintf("ME: %d VS OPNT: %d", ME.Brewed, OPNT.Brewed),
		)
	}
}

func EasierPotionFor(witch Witch) ActionLesser {
	return func(ai, aj Action) bool {
		if ai.Type != aj.Type || ai.Type != BREW {
			panic("EasierPotion wrong type arguments")
		}
		// it's tricky because we actually want the less complex one at the end of the slice
		// however BREW deltas are negative so it works out in the end
		return witch.AddIngredients(ai.Ingredients).Complexity() < witch.AddIngredients(aj.Ingredients).Complexity()
	}
}
