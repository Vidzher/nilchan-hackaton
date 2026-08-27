package progress

type LevelInfo struct {
	Level              int
	ActiveBacklogLimit int
}

type levelThreshold struct {
	RequiredXP   int64
	Level        int
	BacklogLimit int
}

var levels = []levelThreshold{
	{
		RequiredXP:   0,
		Level:        1,
		BacklogLimit: 5,
	},
	{
		RequiredXP:   120,
		Level:        2,
		BacklogLimit: 6,
	},
	{
		RequiredXP:   300,
		Level:        3,
		BacklogLimit: 7,
	},
	{
		RequiredXP:   600,
		Level:        4,
		BacklogLimit: 8,
	},
	{
		RequiredXP:   1000,
		Level:        5,
		BacklogLimit: 10,
	},
}

func FromXP(xp int64) LevelInfo {
	result := levels[0]
	for _, threshold := range levels {
		if xp < threshold.RequiredXP {
			break
		}
		result = threshold
	}
	return LevelInfo{
		Level:              result.Level,
		ActiveBacklogLimit: result.BacklogLimit,
	}
}
