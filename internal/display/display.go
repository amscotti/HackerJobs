package display

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// DisplayResults formats and prints the search results
func DisplayResults(links map[string]string) {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Bold(true)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Padding(0, 1).
		Bold(true)

	linkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3498db")).
		Underline(true).
		Width(55)

	snippetStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2ecc71")).
		Width(75)

	dividerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#95a5a6"))

	// Create and style the header
	header := titleStyle.Render(fmt.Sprintf("Found %d postings", len(links)))
	fmt.Println(header)

	// Add a newline after the header
	fmt.Println()

	// Create column layout with background colors
	linkHeader := headerStyle.Width(55).Render("Link")
	snippetHeader := headerStyle.Width(75).Render("Posting Snippet")
	columns := lipgloss.JoinHorizontal(
		lipgloss.Top,
		linkHeader,
		lipgloss.NewStyle().Width(2).Render(""), // Spacing between columns
		snippetHeader,
	)
	fmt.Println(columns)

	divider := dividerStyle.Render(strings.Repeat("─", 132))
	fmt.Println(divider)

	// Display results
	for id, text := range links {
		link := linkStyle.Render(truncateString(fmt.Sprintf("https://news.ycombinator.com/item?id=%s", id), 52))
		snippet := snippetStyle.Render(truncateString(strings.Split(text, "\n")[0], 72))
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			link,
			lipgloss.NewStyle().Width(2).Render(""), // Spacing between columns
			snippet,
		)
		fmt.Println(row)
	}
}

// truncateString cuts off a string if it exceeds maxLength and adds an ellipsis
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}
