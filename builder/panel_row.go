package builder

func NewPanelRow(title string) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			GridPos: PanelGridPos{
				H: 1,
				W: DashboadWidth,
				X: 0,
				Y: settings.gridPos.Y,
			},
			Title: title,
			Type:  "row",
		}
	}
}
