package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"ecommerce-cli/internal/tui/admin"
	"ecommerce-cli/internal/tui/client"
)

var (
	// Palette de couleurs professionnelles
	colorIndigo  = lipgloss.Color("#6366F1") // Accent principal
	colorSky     = lipgloss.Color("#38BDF8") // Accent secondaire cyan
	colorEmerald = lipgloss.Color("#10B981") // Statut en ligne
	colorSlate   = lipgloss.Color("#94A3B8") // Texte secondaire
	colorDark    = lipgloss.Color("#1E293B") // Fond de badge
	colorBorder  = lipgloss.Color("#334155") // Bordure inactive
	colorWhite   = lipgloss.Color("#F8FAFC") // Texte clair
	colorAmber   = lipgloss.Color("#F59E0B") // Accent d'avertissement
)

type clockTickMsg time.Time

func tickClock() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return clockTickMsg(t)
	})
}

type launcherModel struct {
	choice       int // 0 = menu portail, 1 = client, 2 = admin
	menuIndex    int // 0 = client, 1 = admin
	clientModel  tea.Model
	adminModel   tea.Model
	width        int
	height       int
	spinner      spinner.Model
	startTime    time.Time
	remoteUser   string
	remoteIP     string
	termType     string
}

func newLauncherModel(initialChoice int, user, remoteIP, termType string) *launcherModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(colorSky)

	if user == "" {
		user = "invite"
	}
	if termType == "" {
		termType = "xterm-256color"
	}

	m := &launcherModel{
		choice:     initialChoice,
		menuIndex:  0,
		spinner:    sp,
		startTime:  time.Now(),
		remoteUser: user,
		remoteIP:   remoteIP,
		termType:   termType,
	}

	if initialChoice == 1 {
		m.clientModel = client.NewModel()
	} else if initialChoice == 2 {
		m.adminModel = admin.NewModel()
	}
	return m
}

func (m *launcherModel) Init() tea.Cmd {
	if m.choice == 1 && m.clientModel != nil {
		return m.clientModel.Init()
	}
	if m.choice == 2 && m.adminModel != nil {
		return m.adminModel.Init()
	}
	return tea.Batch(m.spinner.Tick, tickClock())
}

func (m *launcherModel) launchChoice(choice int) (tea.Model, tea.Cmd) {
	m.choice = choice
	var cmd tea.Cmd
	if choice == 1 {
		m.clientModel = client.NewModel()
		if m.width > 0 && m.height > 0 {
			m.clientModel, _ = m.clientModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
		cmd = m.clientModel.Init()
	} else if choice == 2 {
		m.adminModel = admin.NewModel()
		if m.width > 0 && m.height > 0 {
			m.adminModel, _ = m.adminModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		}
		cmd = m.adminModel.Init()
	}
	return m, cmd
}

func (m *launcherModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		var cmds []tea.Cmd
		if m.clientModel != nil {
			var cmd tea.Cmd
			m.clientModel, cmd = m.clientModel.Update(msg)
			cmds = append(cmds, cmd)
		}
		if m.adminModel != nil {
			var cmd tea.Cmd
			m.adminModel, cmd = m.adminModel.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case spinner.TickMsg:
		if m.choice == 0 {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case clockTickMsg:
		if m.choice == 0 {
			return m, tickClock()
		}

	case tea.KeyMsg:
		if m.choice == 0 {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.menuIndex > 0 {
					m.menuIndex--
				}
				return m, nil
			case "down", "j":
				if m.menuIndex < 1 {
					m.menuIndex++
				}
				return m, nil
			case "enter", " ":
				if m.menuIndex == 0 {
					return m.launchChoice(1)
				}
				return m.launchChoice(2)
			case "1", "c":
				return m.launchChoice(1)
			case "2", "a":
				return m.launchChoice(2)
			}
			return m, nil
		}
	}

	if m.choice == 1 && m.clientModel != nil {
		var cmd tea.Cmd
		m.clientModel, cmd = m.clientModel.Update(msg)
		return m, cmd
	}

	if m.choice == 2 && m.adminModel != nil {
		var cmd tea.Cmd
		m.adminModel, cmd = m.adminModel.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *launcherModel) View() string {
	if m.choice == 1 && m.clientModel != nil {
		return m.clientModel.View()
	}
	if m.choice == 2 && m.adminModel != nil {
		return m.adminModel.View()
	}

	sessionDuration := time.Since(m.startTime).Truncate(time.Second)

	// Largeur du conteneur
	boxWidth := 74
	if m.width > 0 && m.width < 76 {
		boxWidth = m.width - 2
	}

	// Styles principaux
	bannerBorder := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(colorIndigo).
		Padding(0, 2).
		Width(boxWidth).
		Align(lipgloss.Center)

	bannerTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Render("E - C O M M E R C E   S E C U R E   G A T E W A Y")

	bannerSub := lipgloss.NewStyle().
		Foreground(colorSlate).
		Render("Plateforme transactionnelle en ligne de commande - Architecture Charmbracelet")

	statusLine := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(colorEmerald).Bold(true).Render("[EN LIGNE] "),
		lipgloss.NewStyle().Foreground(colorSlate).Render("| Port 2222 | Chiffrement ED25519 | Session : "),
		lipgloss.NewStyle().Foreground(colorSky).Bold(true).Render(fmt.Sprintf("%02dm %02ds", int(sessionDuration.Minutes()), int(sessionDuration.Seconds())%60)),
		lipgloss.NewStyle().Render(" "),
		m.spinner.View(),
	)

	// Cartes d'application
	cardStyleActive := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorIndigo).
		Background(colorDark).
		Padding(1, 2).
		Width(boxWidth - 4)

	cardStyleInactive := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1, 2).
		Width(boxWidth - 4)

	// Contenu Carte 1 : Boutique
	var card1 string
	c1Title := lipgloss.NewStyle().Bold(true).Foreground(colorSky).Render("[1] Interface Boutique & Achats (Client TUI)")
	c1Desc := lipgloss.NewStyle().Foreground(colorSlate).Render("Parcourir le catalogue, rechercher des articles, gerer son panier et regler par carte bancaire.")
	c1Badge := lipgloss.NewStyle().Foreground(colorEmerald).Render("[DISPONIBLE] Acces public catalogue et commandes")
	if m.menuIndex == 0 {
		cursor := lipgloss.NewStyle().Foreground(colorIndigo).Bold(true).Render("> Appuyez sur Entree pour lancer la boutique")
		card1 = cardStyleActive.Render(fmt.Sprintf("%s\n%s\n\n%s   %s", c1Title, c1Desc, c1Badge, cursor))
	} else {
		cursor := lipgloss.NewStyle().Foreground(colorSlate).Render("Touche [1] ou fleches pour selectionner")
		card1 = cardStyleInactive.Render(fmt.Sprintf("%s\n%s\n\n%s   %s", c1Title, c1Desc, c1Badge, cursor))
	}

	// Contenu Carte 2 : Administration
	var card2 string
	c2Title := lipgloss.NewStyle().Bold(true).Foreground(colorAmber).Render("[2] Console d'Administration (Admin TUI)")
	c2Desc := lipgloss.NewStyle().Foreground(colorSlate).Render("Superviser les comptes, editer le catalogue, piloter les stocks et le cycle de vie des commandes.")
	c2Badge := lipgloss.NewStyle().Foreground(colorAmber).Render("[RESTREINT] Acces gestionnaires - Authentification Bearer")
	if m.menuIndex == 1 {
		cursor := lipgloss.NewStyle().Foreground(colorIndigo).Bold(true).Render("> Appuyez sur Entree pour lancer l'administration")
		card2 = cardStyleActive.Render(fmt.Sprintf("%s\n%s\n\n%s   %s", c2Title, c2Desc, c2Badge, cursor))
	} else {
		cursor := lipgloss.NewStyle().Foreground(colorSlate).Render("Touche [2] ou fleches pour selectionner")
		card2 = cardStyleInactive.Render(fmt.Sprintf("%s\n%s\n\n%s   %s", c2Title, c2Desc, c2Badge, cursor))
	}

	// Infos système en bas
	metaStyle := lipgloss.NewStyle().Foreground(colorSlate).Width(boxWidth).Align(lipgloss.Left)
	infoText := fmt.Sprintf("Utilisateur SSH: %s | Client: %s | Terminal: %s", m.remoteUser, m.remoteIP, m.termType)

	// Aide navigation
	helpKey := lipgloss.NewStyle().Foreground(colorIndigo).Bold(true)
	helpMuted := lipgloss.NewStyle().Foreground(colorSlate)
	footerHelp := fmt.Sprintf("%s %s   %s %s   %s %s   %s %s",
		helpKey.Render("[Haut/Bas] :"), helpMuted.Render("Changer"),
		helpKey.Render("[Entree] :"), helpMuted.Render("Lancer"),
		helpKey.Render("[1 / 2] :"), helpMuted.Render("Acces direct"),
		helpKey.Render("[q] :"), helpMuted.Render("Deconnexion"),
	)

	// Construction complète
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(bannerBorder.Render(fmt.Sprintf("%s\n%s", bannerTitle, bannerSub)))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().PaddingLeft(2).Render(statusLine))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().PaddingLeft(2).Render(card1))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().PaddingLeft(2).Render(card2))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().PaddingLeft(2).Render(metaStyle.Render(infoText)))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().PaddingLeft(2).Render(footerHelp))
	b.WriteString("\n")

	return b.String()
}

func main() {
	port := os.Getenv("SSH_PORT")
	if port == "" {
		port = "2222"
	}
	addr := fmt.Sprintf(":%s", port)

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				pty, _, active := s.Pty()
				if !active {
					log.Println("Aucun terminal interactif (PTY) alloué")
					return nil, nil
				}

				initialChoice := 0
				cmd := s.Command()
				if len(cmd) > 0 {
					switch strings.ToLower(cmd[0]) {
					case "client", "shop", "1":
						initialChoice = 1
					case "admin", "panel", "2":
						initialChoice = 2
					}
				}

				remoteAddr := s.RemoteAddr().String()
				user := s.User()
				term := pty.Term

				m := newLauncherModel(initialChoice, user, remoteAddr, term)
				return m, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
			}),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalf("Erreur d'initialisation du serveur SSH: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Serveur SSH E-commerce démarré sur le port %s\n", port)
	go func() {
		if err := s.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			log.Fatalf("Erreur du serveur SSH: %v", err)
		}
	}()

	<-done
	log.Println("Arrêt du serveur SSH en cours...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Erreur lors de l'arrêt du serveur SSH: %v", err)
	}
	log.Println("Serveur SSH arrêté proprement.")
}
