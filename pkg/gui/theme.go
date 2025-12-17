package gui

import "github.com/gdamore/tcell"

// Theme defines the color palette for the UI
type Theme struct {
	Fg           tcell.Color // Default text color
	Bg           tcell.Color // Default background color
	Border       tcell.Color // Border color
	Title        tcell.Color // Panel title color
	Header       tcell.Color // Table header color
	SelectedFg   tcell.Color // Selected row text color
	SelectedBg   tcell.Color // Selected row background color
	Keybinding   tcell.Color // Navigation keybinding text color
	InfoLabel    tcell.Color // Labels in info panel
	InfoValue    tcell.Color // Values in info panel
	StatusUp     tcell.Color // Status: Up/Running
	StatusDown   tcell.Color // Status: Down/Exited
	StatusWarn   tcell.Color // Status: Paused/Warning
	ListItem     tcell.Color // Default list item color
	Images       tcell.Color // Image specific color
	Volumes      tcell.Color // Volume specific color
	Networks     tcell.Color // Network specific color
	CleanupItems tcell.Color // Cleanup item color
	Tasks        tcell.Color // Task text color
}

// CurrentTheme holds the current theme configuration
var CurrentTheme = &Theme{
	Fg:           tcell.ColorWhite,
	Bg:           tcell.ColorDefault,
	Border:       tcell.ColorDarkSlateGray, // Darker border for "geek" feel
	Title:        tcell.ColorLightCyan,     // Cyberpunk Cyan for titles
	Header:       tcell.ColorLightSkyBlue,  // Distinct header color
	SelectedFg:   tcell.ColorWhite,
	SelectedBg:   tcell.ColorDarkViolet, // Deep purple for selection
	Keybinding:   tcell.ColorYellow,     // High contrast for help
	InfoLabel:    tcell.ColorLightGreen, // Matrix green for labels
	InfoValue:    tcell.ColorWhite,
	StatusUp:     tcell.ColorGreenYellow, // Bright green for success
	StatusDown:   tcell.ColorCrimson,     // Deep red for errors/stopped
	StatusWarn:   tcell.ColorOrange,
	ListItem:     tcell.ColorLightGray,
	Images:       tcell.ColorPlum,        // Soft purple for images
	Volumes:      tcell.ColorLightSalmon, // Soft pink/orange for volumes
	Networks:     tcell.ColorLightBlue,   // Blue for networks
	CleanupItems: tcell.ColorSandyBrown,
	Tasks:        tcell.ColorLightGreen,
}

// GetStatusColor returns the color for a container state
func GetStatusColor(state string) tcell.Color {
	switch state {
	case "exited", "dead":
		return CurrentTheme.StatusDown
	case "running":
		return CurrentTheme.StatusUp
	case "paused":
		return CurrentTheme.StatusWarn
	case "restarting":
		return tcell.ColorLightSkyBlue
	case "removing":
		return tcell.ColorDarkMagenta
	case "created":
		return tcell.ColorLightCyan
	default:
		return tcell.ColorWhite
	}
}
