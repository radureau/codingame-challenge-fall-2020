package main

import (
	"fmt"
	"sort"
)

const (
	MAX_ORDERS = 5
)

type Point int

type Witch struct {
	Ingredients
	Points Point
	Me     bool
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

type ActionType string

const (
	BREW          = ActionType("BREW")
	CAST          = ActionType("CAST")
	OPPONENT_CAST = ActionType("OPPONENT_CAST")
	LEARN         = ActionType("LEARN")
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

func (a Action) IsLessThan(other Action) bool {
	return a.Points < other.Points
}

type ActionSlice []Action

func (as ActionSlice) Len() int           { return len(as) }
func (as ActionSlice) Less(i, j int) bool { return as[i].IsLessThan(as[j]) }
func (as ActionSlice) Swap(i, j int)      { as[i], as[j] = as[j], as[i] }

func (as ActionSlice) Last() Action { return as[len(as)-1] }
func (as ActionSlice) Pick() Action {
	sort.Sort(as)
	return as.Last()
}

func main() {

	for {
		Actions := ActionSlice{}

		var actionCount int // actionCount: the number of spells and recipes in play
		fmt.Scan(&actionCount)
		for i := 0; i < actionCount; i++ {
			var delta0, delta1, delta2, delta3, _castable, _repeatable int
			action := Action{}
			fmt.Scan(&action.ID, &action.Type, &delta0, &delta1, &delta2, &delta3, &action.Points, action.TomeIndex, action.TaxCount, &_castable, &_repeatable)
			action.Castable, action.Repeatable = _castable != 0, _repeatable != 0
			action.Ingredients = NewIngredients(delta0, delta1, delta2, delta3)

			switch ActionType(action.Type) {
			case BREW:
			case CAST, OPPONENT_CAST, LEARN:
				fallthrough
			default:
				panic(fmt.Errorf("unknown action type %s", action.Type))
			}
			Actions = append(Actions, action)
		}

		var inv0, inv1, inv2, inv3 int

		ME := Witch{Me: true}
		fmt.Scan(&inv0, &inv1, &inv2, &inv3, &ME.Points)
		ME.Ingredients = NewIngredients(inv0, inv1, inv2, inv3)

		OPNT := Witch{}
		fmt.Scan(&inv0, &inv1, &inv2, &inv3, &OPNT.Points)
		OPNT.Ingredients = NewIngredients(inv0, inv1, inv2, inv3)

		// in the first league: BREW <id> | WAIT; later: BREW <id> | CAST <id> [<times>] | LEARN <id> | REST | WAIT
		fmt.Printf("BREW %d\n", Actions.Pick().ID)
	}
}
