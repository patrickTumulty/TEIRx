package core

const (
	TIER_S int = 6
	TIER_A int = 5
	TIER_B int = 4
	TIER_C int = 3
	TIER_D int = 2
	TIER_F int = 1
)

func TierInt2Str(tier int) rune {
	switch tier {
	case TIER_S:
		return 'S'
	case TIER_A:
		return 'A'
	case TIER_B:
		return 'B'
	case TIER_C:
		return 'C'
	case TIER_D:
		return 'D'
	case TIER_F:
		return 'F'
	default:
		return 0
	}
}

func TierStr2Int(tier rune) int {
	if tier == 'S' {
		return TIER_S
	}
	if tier == 'A' {
		return TIER_A
	}
	if tier == 'B' {
		return TIER_B
	}
	if tier == 'C' {
		return TIER_C
	}
	if tier == 'D' {
		return TIER_D
	}
	if tier == 'F' {
		return TIER_F
	}
	return -1
}
