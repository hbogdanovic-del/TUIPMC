package main

import (
	"fmt"
	"maps"
	"os/exec"
	"strings"

	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const cols = 2

type model struct {
	cats     []string
	logins   map[string][]string
	cols     int
	rows     int
	row      int
	col      int
	w        int
	h        int
	selected string
}

func initialModel(logins map[string][]string) model {
	cats := []string{}
	for cat := range maps.Keys(logins) {
		cats = append(cats, cat)
	}
	cols := len(cats)
	t_max := -1
	for _, cat := range cats {
		cat_logins := logins[cat]
		nb_log := len(cat_logins)
		t_max = max(t_max, nb_log)
	}
	rows := t_max
	return model{
		cats:   cats,
		logins: logins,
		cols:   cols,
		rows:   rows,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Current() []string {
	return m.logins[m.cats[m.col]]
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if s, ok := msg.(tea.WindowSizeMsg); ok {
		m.w, m.h = s.Width, s.Height
	}
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.row > 0 {
				m.row--
			}
		case "down", "j":
			if m.row < len(m.Current())-1 {
				m.row++
			}
		case "left", "h":
			if m.col > 0 {
				m.col--
				m.row = min(m.row, len(m.Current())-1)
			}
		case "right", "l":
			if m.col < len(m.cats)-1 {
				m.col++
				m.row = min(m.row, len(m.Current())-1)
			}
		case "enter":
			m.selected = m.Current()[m.row]
			return m, tea.Quit
		}
	}
	return m, nil
}

/*
 * 80col * 24row
 * ┌─┐
 * │ │
 * └─┘
 * ┌──────────────────────────────────────────────────────────────────────────────┐
 * │                                                                              │
 * │                                  Acad : XXX                                  │
 * │                                                                              │
 * │──────────────────────────────────────────────────────────────────────────────│
 * │                                      │                                       │
 * │               PUBLIC                 │               PRIVE                   │
 * │                                      │                                       │
 * │──────────────────────────────────────────────────────────────────────────────│
 * │                                      │                                       │
 * │ - EPP                          CC    │ - E3P                          CC     │
 * │                                      │                                       │
 * │ - AGAPEYYY                     CC    │ - AGAPRIYYY                    CC     │
 * │                                      │                                       │
 * │ - AGAPEYYY                     CC    │                                       │
 * │                                      │                                       │
 * │                                      │                                       │
 * │                                      │                                       │
 * │                                      │                                       │
 * │                                      │                                       │
 * │                                      │                                       │
 * │                                      │                                       │
 * │                                      │                                       │
 * └──────────────────────────────────────────────────────────────────────────────┘
 */
func (m model) View() tea.View {
	/*
		remaining_line := 2
		s := "┌──────────────────────────────────────────────────────────────────────────────┐\n"
		s += "│                                                                              │\n"
		s += "│                                   ACA XXX                                    │\n"
		s += "│                                                                              │\n"
		s += "│──────────────────────────────────────────────────────────────────────────────│\n"
		s += "│                                      │                                       │\n"
		s += "│               PUBLIC                 │                PRIVE                  │\n"
		s += "│                                      │                                       │\n"
		s += "│──────────────────────────────────────────────────────────────────────────────│\n"
		s += "│                                      │                                       │\n"
		s += "│                                      │                                       │\n"
		for range remaining_line {
			s += "│e                                                                             │\n"
		}
		s += "│ Appuyer sur Q pour quitter                                                   │\n"
		s += "└──────────────────────────────────────────────────────────────────────────────┘\n"
		return tea.NewView(s)
	*/
	if m.w == 0 {
		return tea.NewView("")
	}
	if m.w < 80 || m.h < 24 {
		return tea.NewView(fmt.Sprintf("La fenêtre est trop petite: %dx%d", m.w, m.h))
	}
	var foreground = lipgloss.RGBColor{R: 187, G: 187, B: 187}
	var bold_foreground = lipgloss.RGBColor{R: 255, G: 255, B: 255}
	var background = lipgloss.RGBColor{R: 0, G: 0, B: 0}
	heigth := 24
	width := 80
	n := len(m.cats)
	inner := width - 1 - 1
	colW := (inner - (n - 1)) / n
	body := heigth - 2 - 2
	header := lipgloss.NewStyle().Bold(true).
		Underline(true).
		Foreground(bold_foreground)
	item := lipgloss.NewStyle().Width(colW).MaxWidth(colW).
		PaddingLeft(2)
	selected := lipgloss.NewStyle().Bold(true).
		Foreground(bold_foreground).
		Background(background)
	base := lipgloss.NewStyle().
		Height(body)
	sep := lipgloss.NewStyle().Foreground(foreground).
		Render(strings.TrimSuffix(strings.Repeat("|\n", body), "\n"))
	var parts []string
	for c, cat := range m.cats {
		lines := []string{header.Render(strings.ToUpper(cat)), ""}
		for r, login := range m.logins[cat] {
			if c == m.col && r == m.row {
				lines = append(lines, selected.Render(fmt.Sprintf("> %s", login)))
			} else {
				lines = append(lines, item.Render(fmt.Sprintf("%s", login)))
			}
		}
		style := base
		if c > 0 {
			parts = append(parts, sep)
		}
		parts = append(parts, style.Render(strings.Join(lines, "\n")))
	}
	grid := lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	frame := lipgloss.NewStyle().
		Border(lipgloss.ASCIIBorder()).
		BorderBackground(background).
		Padding(0, 1).
		Width(width - 2).
		Height(body)
	help := lipgloss.NewStyle().Foreground(foreground).
		Render("Utiliser les fleches ou hjkl pour se deplacer. Q pour quitter")
	screen := lipgloss.JoinVertical(lipgloss.Left, frame.Render(grid), help)

	return tea.NewView(lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center, screen))
}

func getCoeurUsers(users []string) map[string][]string {
	prefixes := []string{"epp", "agape", "agapri", "prive"}
	groups := make(map[string][]string)
	for _, user := range users {
		for _, prefix := range prefixes {
			if strings.HasPrefix(user, prefix) {
				switch prefix {
				case "epp":
					groups["Public"] = append(groups["Public"], user)
				case "agape":
					groups["Public"] = append(groups["Public"], user)
				case "prive":
					groups["Prive"] = append(groups["Prive"], user)
				case "agapri":
					groups["Prive"] = append(groups["Prive"], user)
				}

			}
		}
	}
	return groups
}

func getUsersFromPasswd() []string {
	out, err := exec.Command("cat", "/etc/passwd").Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	l_out := strings.Split(string(out), "\n")
	l_out = l_out[:len(l_out)-1]
	users := []string{}
	for _, line := range l_out {
		user := strings.Split(line, ":")[0]
		users = append(users, user)
	}
	return users
}

func main() {
	users := getUsersFromPasswd()
	users_coeur := getCoeurUsers(users)
	/*
		for _, user := range users_coeur {
			fmt.Printf("%s\n", user)
		}
	*/

	p := tea.NewProgram(initialModel(users_coeur))
	final, err := p.Run()
	if err != nil {
		fmt.Printf("Mince! Y a eu une erreu: %v", err)
		os.Exit(1)
	}
	choix := final.(model).selected
	if choix == "" {
		return
	}
	if err := launch(buildArgv(choix)); err != nil {
		fmt.Fprintf(os.Stderr, "Echec du lancement de la commande :%v", err)
		os.Exit(1)
	}
}
