package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	PaddingY = 1
	PaddingX = 2
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("248")).
	Padding(PaddingY, PaddingX)

// AddTitleToBorder replaces dashes used for border with a title.
func AddTitleToBorder(in, title string) string {
	return strings.Replace(
		in,
		strings.Repeat("─", 2+len(title)),
		" "+title+" ",
		1,
	)
}

// View renders the DPS meter.
func (s *State) View() string {
	// Grab the terminal size for all the fun math stuff.
	width, height, err := terminal.GetSize(0)
	if err != nil {
		return fmt.Sprintf("Error: %v\n", err)
	}

	out := "\nHello, Player!\n"
	if s.PrimaryPlayer != nil {
		out = fmt.Sprintf("\nHello, %v!\n", s.PrimaryPlayer.Name)
	}

	out = baseStyle.
		Bold(true).
		Width(32). // Enough to cover title
		Render(out)

	out = AddTitleToBorder(out, "FFXIV - TUI - DAMAGE METER")

	out += "\n"

	// Encounter Details
	if s.CombatData != nil {
		enc := s.CombatData.Encounter

		p := message.NewPrinter(language.English)

		damage, _ := strconv.Atoi(enc.Damage)
		dps, _ := strconv.ParseFloat(enc.DPS, 32)

		encounterStats := fmt.Sprintf("Duration: %v\n", enc.Duration)
		encounterStats += fmt.Sprintf("Damage:   %v\n", p.Sprintf("%d", damage))
		encounterStats += fmt.Sprintf("DPS:      %v", p.Sprintf("%.0f", dps))
		encounterStats = baseStyle.MarginLeft(2).Render(encounterStats)

		out = lipgloss.JoinHorizontal(
			lipgloss.Top,
			out,
			AddTitleToBorder(encounterStats, "ENCOUNTER"),
		)

		combatants := s.GetSortedCombatants(CombatantSortOptions{
			IncludeLimitBreak: false,
		})

		combatantBoxHeight := height - lipgloss.Height(out) - 4
		combatantBoxWidth := lipgloss.Width(out) + 10
		maxNumPlayers := (combatantBoxHeight - 4) / 3

		// Attempt to scale # of players based on available height
		if len(combatants) > maxNumPlayers {
			combatants = combatants[:maxNumPlayers]
		}

		combatantStats := ""

		meterWidth := combatantBoxWidth - 4

		for i, c := range combatants {
			stats := fmt.Sprintf(
				"%v: %v (%v) - %v (%v) [%v]\n",
				i+1,
				c.Name,
				c.Job,
				c.DPS,
				c.Damage,
				c.DamagePercent,
			)

			roleColor := GetRoleColor(c.Job)

			perc, _ := strconv.ParseFloat(strings.Replace(c.DamagePercent, "%", "", 1), 32)
			perc /= 100

			meter := ""

			if perc > 0 {
				bars := float64(meterWidth) * perc

				meter = lipgloss.NewStyle().
					Background(lipgloss.Color(roleColor)).
					Render(strings.Repeat(" ", int(bars)))
			}

			combatantStats += lipgloss.NewStyle().
				Bold(c.Name == "YOU").
				Foreground(lipgloss.Color(roleColor)).
				Render(stats + meter)

			if i < len(combatants)-1 {
				combatantStats += "\n\n"
			}
		}

		out += AddTitleToBorder(
			baseStyle.
				Height(combatantBoxHeight).
				Width(combatantBoxWidth).
				MarginTop(1).
				Render(combatantStats),
			"COMBATANTS",
		)
	}

	return baseStyle.
		UnsetBorderStyle().
		Width(width - (PaddingX * 2)).
		Height(height - (PaddingY * 2)).
		MaxHeight(height). // Keep top visible when list is longer than viewport
		Align(lipgloss.Center).
		Render(out)
}
