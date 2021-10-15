package builder

func NewPanelRow(title string) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
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

func NewPanelRowCollapsed(title string, panels []Panel) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Collapsed: true,
			GridPos: PanelGridPos{
				H: 1,
				W: DashboadWidth,
				X: 0,
				Y: gridPos.Y,
			},
			Title:  title,
			Type:   "row",
			Panels: panels,
		}
	}
}
