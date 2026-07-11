# lem-in — paths/best.go

## CanSendThisPath

```go
func CanSendThisPath(set *core.PathsSet, i int, ants int) bool {
	if i == 0 {
		return true
	}

	sum := 0
	for j := i - 1; j >= 0; j-- {
		if set.Lengths[j] < set.Lengths[i] {
			sum += set.Lengths[i] - set.Lengths[j]
		}
	}

	return ants > sum
}
```
Decides whether it is worth sending ants down path `i` given how many ants are left.
- The first path (index 0, shortest) always gets ants
- For longer paths, calculates the total extra steps compared to all shorter paths
- If the remaining ants outnumber the total extra steps, it is worth using this path
- Returns false if using this path would not save moves

> This is the key optimization. A longer path is only worth using if there are enough ants to justify it. For example if path 0 has length 2 and path 1 has length 5, path 1 only helps if there are more than 3 ants left — otherwise all ants are better off on path 0.

---

## checkLongestMove

```go
func checkLongestMove(d *core.Data, set *core.PathsSet, antsToPath []int) {
	longest := 0

	for i := 0; i < set.PathsAmount; i++ {
		val := antsToPath[i] + set.Lengths[i]
		if val > longest {
			longest = val
		}
	}

	if d.BestSet == nil || longest < d.BestSpeed {
		d.BestSet = set
		d.BestSpeed = longest
	}
}
```
Calculates the total number of moves this path set would take and checks if it is better than the current best.
- For each path, adds the number of ants assigned to it plus the path length
- The total moves is the maximum of all these values — because the slowest path determines when we finish
- If this is better than the current best (fewer moves), saves it as the new best

---

## CheckIfCurrentIsBest

```go
func CheckIfCurrentIsBest(d *core.Data, set *core.PathsSet) {
	ants := d.Ants
	antsToPath := make([]int, set.PathsAmount)

	for ants > 0 {
		for i := 0; i < set.PathsAmount && ants > 0; i++ {
			if CanSendThisPath(set, i, ants) {
				antsToPath[i]++
				ants--
			}
		}
	}

	checkLongestMove(d, set, antsToPath)
}
```
Distributes all ants across the current paths optimally and checks if this set is the fastest.
- Keeps assigning ants one by one to paths in order
- Only assigns an ant to a path if `CanSendThisPath` says it is worth it
- Repeats until all ants are assigned
- Calls `checkLongestMove` to compare this distribution against the current best

---

## BestPathsSetOperations

```go
func BestPathsSetOperations(d *core.Data) {
	set := BuildPathsSetStructure(d)
	if set.PathsAmount == 0 {
		return
	}

	CalculatePathsLengths(d, set)
	AllocatePathsArrays(set)
	SaveCurrentPathsSet(d, set)
	SortPathsShortToLong(set)
	CheckIfCurrentIsBest(d, set)
}
```
Runs all path operations after each successful BFS run.
- Builds the path set structure
- Calculates lengths, allocates arrays, saves the paths, sorts them
- Checks if this set is the best one found so far
- If no paths exist yet, returns immediately

> This is called after every BFS run. Each call considers one more path than the previous call. The best combination across all calls is stored in `d.BestSet` and used by the simulator.
---
