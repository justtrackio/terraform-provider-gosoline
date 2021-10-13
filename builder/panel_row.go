package builder

func NewPanelRow(title string, collapsed bool) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			collapsed: collapsed,
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
