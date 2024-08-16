package matchmaking

import (
	"sort"

	"github.com/ZdravkoGyurov/wipgame/matchmaking-server/internal/types"
)

var queue = []types.Player{}

var ratingRanges = map[string]int{} // player id -> player ratingRange

func insertPlayer(slice []types.Player, value types.Player) []types.Player {
	index := sort.Search(len(slice), func(i int) bool {
		return slice[i].Rating >= value.Rating
	})
	slice = append(slice, types.Player{})
	copy(slice[index+1:], slice[index:])
	slice[index] = value
	return slice
}

func playerIdx(slice []types.Player, value types.Player) int {
	low := 0
	high := len(slice) - 1
	for low <= high {
		mid := (low + high) / 2
		if slice[mid].Rating == value.Rating {
			return mid
		}
		if slice[mid].Rating > value.Rating {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}

func removePlayer(slice []types.Player, idx int) []types.Player {
	return append(slice[:idx], slice[idx+1:]...)
}
