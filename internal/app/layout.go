package app

import (
	"github.com/blumenwagen/durandal/internal/components"
	"github.com/charmbracelet/lipgloss"
)

// calculateLayout distributes terminal dimensions for a custom asymmetric UI.
// Every pixel is accounted for.
func calculateLayout(m Model) Model {
	w := m.Width
	h := m.Height

	if w < 40 || h < 15 {
		return m
	}

	// Reserve fixed rows: header (1) + help bar (1) = 2
	usableH := h - 2

	// Process list takes full usable height
	procH := usableH

	// Width split: Left stats (35%) | Process list (65%)
	leftW := w * 35 / 100
	if leftW < 30 {
		leftW = 30
	}
	rightW := w - leftW
	if rightW < 30 {
		rightW = 30
		leftW = w - rightW
	}

	// GPU check
	gpuH := 0
	if len(m.GPU.GPUs) > 0 {
		gpuH = usableH * 15 / 100
	}

	// Vertical split for left column: SENTINEL | CPU | [GPU] | MEM | NET | DISK
	sentinelH := usableH * 16 / 100
	if sentinelH < 5 {
		sentinelH = 5
	}
	cpuH := usableH * 22 / 100
	memH := usableH * 22 / 100
	netH := usableH * 17 / 100

	// Recalculate Disk to absorb rounding so total matches exactly usableH
	diskH := usableH - (sentinelH + cpuH + gpuH + memH + netH)

	// Header & HelpBar
	m.Header.Width = w

	// Left column components
	m.Sentinel.Width = leftW
	m.Sentinel.Height = sentinelH
	m.CPU.Width = leftW
	m.CPU.Height = cpuH
	m.GPU.Width = leftW
	m.GPU.Height = gpuH
	m.Memory.Width = leftW
	m.Memory.Height = memH
	m.Network.Width = leftW
	m.Network.Height = netH
	m.Disk.Width = leftW
	m.Disk.Height = diskH

	// Inspector overlays left column
	m.Inspector.Width = leftW
	m.Inspector.Height = usableH

	// Right column
	m.Processes.Width = rightW
	m.Processes.Height = procH

	m.ProcY = 1 // Process list starts immediately after the header

	return m
}

// renderLayout composes all views into the final screen.
func renderLayout(m Model) string {
	header := m.Header.View()

	var leftCol string
	if m.InspectorOpen {
		leftCol = m.Inspector.View()
	} else {
		var panels []string
		panels = append(panels, m.Sentinel.View(), m.CPU.View())

		gpuView := m.GPU.View()
		if gpuView != "" {
			panels = append(panels, gpuView)
		}

		panels = append(panels, m.Memory.View(), m.Network.View(), m.Disk.View())
		leftCol = lipgloss.JoinVertical(lipgloss.Left, panels...)
	}

	procs := m.Processes.View()

	middleRow := lipgloss.JoinHorizontal(lipgloss.Top,
		leftCol,
		procs,
	)

	helpBar := components.HelpBar(m.Width, false)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		middleRow,
		helpBar,
	)
}
