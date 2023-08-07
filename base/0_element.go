package base

var ElemList, ElemIndex, PhysList, PhysIndex = getIndexes()

func getIndexes() ([]string, map[string]int, []string, map[string]int) {
	bufferE := make(map[string]int)
	bufferP := make(map[string]int)
	var fulllist []string = []string{"◌", "🌪", "🔥"}//, "🪨", "🧊", "🌑", "🩸", "🎶", "☀️"}
	var elist []string
	var physlist []string = []string{"◌", "🌱", "🪵", "🪨", "🛡"}
	var plist []string
	for i, str := range fulllist { bufferE[str] = i ; elist = append(elist, str) }
	for i, str := range physlist { bufferP[str] = i ; plist = append(plist, str) }
	return elist, bufferE, plist, bufferE
}

