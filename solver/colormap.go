package solver

type ColorMap struct {
	colorFromName map[string]Color
	nameFromColor []string
}

func NewColorMap() ColorMap {
	colorFromName := make(map[string]Color)
	nameFromColor := make([]string, 0, 20)
	newMap := ColorMap{colorFromName: colorFromName, nameFromColor: nameFromColor}
	return newMap
}

func (cm *ColorMap) ColorFromString(name string) Color {
	color, ok := cm.colorFromName[name]
	if ok {
		return color
	}
	newColor := Color(len(cm.nameFromColor))
	cm.colorFromName[name] = newColor
	cm.nameFromColor = append(cm.nameFromColor, name)
	return newColor
}

func (cm *ColorMap) StringFromColor(color Color) string {
	return cm.nameFromColor[color]
}
