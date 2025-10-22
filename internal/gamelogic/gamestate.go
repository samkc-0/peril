package gamelogic

import (
	"sync"
)

type GameState struct {
	Player Player
	Paused bool
	mu     *sync.RWMutext
}

func NewGameState(username string) *GameState {
	return &GameState{
		Player: Player{Username: username, Units: map[int]Unit{}},
		Paused: false,
		mu:     &sync.RWMutex{},
	}
}

func (gs *GameState) resumeGame() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Paused = false
}

func (gs *GameState) pauseGame() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Paused = true
}

func (gs *GameState) isPaused() bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.Paused
}

func (gs *GameState) addUnit(u Unit) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Player.Units[u.ID] = u
}

func (gs *GameState) removeUnitsInLocation(loc Location) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	for unitID, unit := range gs.Player.Units {
		if unit.Location == loc {
			delete(gs.Player.Units, unitID)
		}
	}
}

func (gs *GameState) UpdateUnit(u Unit) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Player.Units[u.ID] = u
}

func (gs *GameState) GetUsername() string {
	return gs.Player.Username
}

func (gs *GameState) getUnitsSnap() []Unit {
	gs.mu.Lock()
	defer gs.mu.RUnlock()
	units := []Unit{}
	for _, unit := range gs.Player.Units {
		units = append(units, unit)
	}
	return units
}

func (gs *GameState) GetUnit(id int) (Unit, bool) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	unit, ok := gs.Player.Units[id]
	return unit, ok
}

func (gs *GameState) GetPlayerSnap() Player {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	units := map[int]Unit{}
	for unitID, unit := range gs.Player.Units {
		units[unitID] = unit
	}
	return Player{
		Username: gs.Player.Username,
		Units:    units,
	}
}
