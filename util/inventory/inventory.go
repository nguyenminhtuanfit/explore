package inventory

import (
	"github.com/jasdel/explore/entity/thing"
	"github.com/jasdel/explore/util/command"
	"sync"
)

const (
	notFound = -1
)

type Interface interface {
	Add(thing.Interface) bool
	Remove(thing.Interface) bool
	Contains(thing.Interface) bool
	List(omit ...thing.Interface) []thing.Interface
}

type Inventory struct {
	things []thing.Interface

	mutex sync.RWMutex
}

func New(initialCap int) *Inventory {
	return &Inventory{
		things: make([]thing.Interface, 0, initialCap),
	}
}

// Adds the given thing to the inventory, ignoring the
// command if the exact thing already exists.
func (i *Inventory) Add(t thing.Interface) bool {
	if find(t, i.things) == notFound {
		i.things = append(i.things, t)
		return true
	}
	return false
}

// Removes the given item from the inventory if it exists
func (i *Inventory) Remove(t thing.Interface) bool {
	if idx := find(t, i.things); idx != notFound {
		i.things = append(i.things[:idx], i.things[idx+1:]...)
		return true
	}
	return false
}

// Returns if the thing is within the inventory
func (i *Inventory) Contains(t thing.Interface) bool {
	return find(t, i.things) != notFound
}

// List returns a slice of thing.Interface in the Inventory, possibly with
// specific items omitted. An example of when you want to omit something is when
// a Player does something - you send a specific message to the player:
//
//  You pick up a ball.
//
// A different message is sent to any observers:
//
//  You see Diddymus pick up a ball.
//
// However when broadcasting the message to the location you want to omit the
// 'actor' who has the specific message.
//
// Note that locations implement an inventory to store what mobiles/players and
// things are present which is why this works.
func (i *Inventory) List(omit ...thing.Interface) []thing.Interface {
	things := make([]thing.Interface, 0, len(i.things))

	for _, t := range i.things {
		if find(t, omit) != notFound {
			continue
		}
		things = append(things, t)
	}

	return things
}

// Returns the number of things in the inventory
func (i *Inventory) Len() int {
	return len(i.things)
}

// Processes the given command against each item in the inventory
// until it is handled, or all items are processed
func (i *Inventory) Process(cmd *command.Command) bool {
	for _, t := range i.things {
		// Don't process the command issuer - gets very recursive!
		if t.IsAlso(cmd.Issuer) {
			continue
		}

		if t, ok := t.(command.Processor); ok {
			if t.Process(cmd) {
				return true
			}
		}
	}
	return false
}

// Searches through the list of things for the target.
// returning its index if found.
func find(t thing.Interface, things []thing.Interface) int {
	for i := 0; i < len(things); i++ {
		if things[i].IsAlso(t) {
			return i
		}
	}
	return notFound
}
