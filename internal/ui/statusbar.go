package ui

import (
	"github.com/born1337/hltui/internal/style"
)

func RenderStatusBar(width int, errMsg string) string {
	if errMsg != "" {
		return style.Red.Render(errMsg)
	}
	hints := "←/→:switch  0-6:views  j/k:scroll  s:sort  r:refresh  ;:help  q:quit"
	return style.Dim.Render(hints)
}
