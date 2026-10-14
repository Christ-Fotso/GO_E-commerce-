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

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ecommerce-cli/internal/tui/admin"
	"ecommerce-cli/internal/tui/client"
)

type launcherModel struct {
	choice      int // 0 = menu, 1 = client, 2 = admin
	clientModel tea.Model
	adminModel  tea.Model
	width       int
	height      int
}

func newLauncherModel(initialChoice int) *launcherModel {
	m := &launcherModel{choice: initialChoice}
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
	return nil
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

	case tea.KeyMsg:
		if m.choice == 0 {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "1", "c":
				m.choice = 1
				m.clientModel = client.NewModel()
				var cmd tea.Cmd
				if m.width > 0 && m.height > 0 {
					m.clientModel, _ = m.clientModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				}
				cmd = m.clientModel.Init()
				return m, cmd
			case "2", "a":
				m.choice = 2
				m.adminModel = admin.NewModel()
				var cmd tea.Cmd
				if m.width > 0 && m.height > 0 {
					m.adminModel, _ = m.adminModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				}
				cmd = m.adminModel.Init()
				return m, cmd
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

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		MarginBottom(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 3).
		Margin(1, 0)

	itemStyle := lipgloss.NewStyle().
		PaddingLeft(2)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	var content strings.Builder
	content.WriteString(titleStyle.Render("★ PORTAIL ACCÈS DISTANT SSH - E-COMMERCE ★"))
	content.WriteString("\n\n")
	content.WriteString("Veuillez sélectionner l'interface que vous souhaitez exécuter :\n\n")
	content.WriteString(itemStyle.Render("[1] Interface Boutique & Achats (Client TUI)\n"))
	content.WriteString(itemStyle.Render("[2] Interface Administration & Gestion (Admin TUI)\n\n"))
	content.WriteString(itemStyle.Render("[q] Quitter la session SSH\n\n"))
	content.WriteString(hintStyle.Render("Appuyez sur 1 ou 2 pour lancer l'interface."))

	return boxStyle.Render(content.String())
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
				_, _, active := s.Pty()
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

				m := newLauncherModel(initialChoice)
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
