package builder

func NewPanelRow(title string) PanelFactory {
	return func(resourceNames ResourceNames, gridPos PanelGridPos) Panel {
		return Panel{
			GridPos: PanelGridPos{
				H: 1,
				W: DashboadWidth,
				X: 0,
				Y: gridPos.Y,
			},
			Title: title,
			Type:  "row",
		}
	}
}
