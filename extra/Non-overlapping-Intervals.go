func eraseOverlapIntervals(intervals [][]int) int {
	if len(intervals) == 0 {
		return 0
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	nonOverlapCount := 1
	end := intervals[0][1]

	for i := 1; i < len(intervals); i++ {
		if intervals[i][0] >= end {
			nonOverlapCount++
			end = intervals[i][1]
		} else if intervals[i][1] < end {
			end = intervals[i][1]
		}
	}

	return len(intervals) - nonOverlapCount
}
