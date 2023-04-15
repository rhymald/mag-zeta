package base

var ElemList, ElemIndex = getIndexes()

func getIndexes() ([]string, map[string]int) {
	buffer := make(map[string]int)
	var fulllist []string = []string{"◌ ", "🌪 ", "🔥", "🪨", "🧊", "🌑", "🩸", "🎶", "☀️ "}
	var list []string
	for i, str := range fulllist { buffer[str] = i ; list = append(list, str) }
	return list, buffer
}