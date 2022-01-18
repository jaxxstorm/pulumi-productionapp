package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	neturl "net/url"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"github.com/jaxxstorm/pulumi-productionapp/cli/pkg/namegenerator"
	"github.com/jaxxstorm/pulumi-productionapp/sdk/go/productionapp"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/alecthomas/kingpin.v2"
)

const columnWidth = 50

// Style definitions.
var (
	app        = kingpin.New("productionapp", "A command-line application deployment tool using pulumi.")
	deployCmd  = app.Command("deploy", "Deploy a container image.")
	destroyCmd = app.Command("destroy", "Destroy a deployment")
	name       = app.Flag("name", "Deployment name to use").String()

	// kingpin vars
	image = deployCmd.Flag("image", "Image to deploy.").Required().String()
	port  = deployCmd.Flag("port", "port container listens on").Default("80").Int()

	subtle  = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	special = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	list = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, true, false, false).
		BorderForeground(subtle).
		MarginRight(2).
		Height(8).
		Width(columnWidth + 1)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	listItem = lipgloss.NewStyle().PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	listDone = func(s string) string {
		return checkMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}

	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

// pulumiProgram is the Pulumi program itself where resources are declared. It deploys a simple static website to S3.
func pulumiProgram(ctx *pulumi.Context) error {

	application, err := productionapp.NewDeployment(ctx, "productionapp", &productionapp.DeploymentArgs{
		Image: pulumi.String(*image),
		Port:  pulumi.Int(*port),
	})
	if err != nil {
		return fmt.Errorf("error creating application: %v", err)
	}

	ctx.Export("url", application.Url)

	return nil
}

// runPulumiUpdate runs the update or destroy commands based on input.
// It takes as arguments a flag to determine update or destroy, a channel to receive log messages
// and another to receive structured events from the Pulumi Engine.
func runPulumiUpdate(destroy bool, logChannel chan<- logMessage, eventChannel chan<- events.EngineEvent) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		projectName := "productionapp-cli"
		// we use a simple stack name here, but recommend using auto.FullyQualifiedStackName for maximum specificity.

		var nameToUse string
		if *name == "" {
			nameToUse = namegenerator.Generate()
		} else {
			nameToUse = *name
		}

		stackName := nameToUse
		// stackName := auto.FullyQualifiedStackName("myOrgOrUser", projectName, stackName)

		// create or select a stack matching the specified name and project.
		// this will set up a workspace with everything necessary to run our inline program (deployFunc)
		s, err := auto.UpsertStackInlineSource(ctx, stackName, projectName, pulumiProgram)
		if err != nil {
			app.Fatalf("Failed to get stack: %v", err)
		}

		logChannel <- logMessage{msg: fmt.Sprintf("Created/Selected stack %q\n", stackName)}

		w := s.Workspace()

		logChannel <- logMessage{msg: "Installing the Kubernetes plugin"}

		// for inline source programs, we must manage plugins ourselves
		err = w.InstallPlugin(ctx, "kubernetes", "v3.5.2")
		if err != nil {
			app.Fatalf("Failed to install program plugins: %v\n", err)

		}

		logChannel <- logMessage{msg: "Successfully installed Kubernetes plugin"}

		logChannel <- logMessage{msg: "Successfully set config"}
		logChannel <- logMessage{msg: "Running refresh..."}

		_, err = s.Refresh(ctx)
		if err != nil {
			app.Fatalf("Failed to refresh stack: %v\n", err)
		}

		logChannel <- logMessage{msg: "Refresh succeeded!"}

		if destroy {
			logChannel <- logMessage{msg: "Running destroy..."}

			// destroy our stack and exit early
			_, err := s.Destroy(ctx, optdestroy.EventStreams(eventChannel))
			if err != nil {
				fmt.Printf("Failed to destroy stack: %v", err)
			}
			logChannel <- logMessage{msg: "Stack successfully destroyed"}
			return logMessage{msg: "Success"}
		}

		logChannel <- logMessage{msg: "Running update..."}

		res, err := s.Up(ctx, optup.EventStreams(eventChannel))
		if err != nil {
			app.Fatalf("Failed to update stack: %v\n\n", err)
		}

		logChannel <- logMessage{msg: "Update succeeded!"}

		// get the URL from the stack outputs
		url, ok := res.Outputs["url"].Value.(string)
		if !ok {
			fmt.Println("Failed to unmarshal output URL")
			os.Exit(1)
		}

		logChannel <- logMessage{msg: fmt.Sprintf("URL: %s\n", url)}
		return logMessage{msg: url}
	}
}

// watchForLogMessages forwards any log messages to the `Update` method
func watchForLogMessages(msg chan logMessage) tea.Cmd {
	return func() tea.Msg {
		return <-msg
	}
}

// watchForEvents forwards any engine events to the `Update` method
func watchForEvents(event chan events.EngineEvent) tea.Cmd {
	return func() tea.Msg {
		return <-event
	}
}

type logMessage struct {
	msg string
}

// model is the struct that holds the state for this program
type model struct {
	eventChannel      chan events.EngineEvent // where we'll receive engine events
	logChannel        chan logMessage         // where we'll receive log messages
	spinner           spinner.Model
	destroy           bool
	quitting          bool
	currentMessage    string
	updatesInProgress map[string]string // resources with updates in progress
	updatesComplete   map[string]string // resources with updates completed
}

// Init runs any IO needed at the initialization of the program
func (m model) Init() tea.Cmd {
	return tea.Batch(
		watchForLogMessages(m.logChannel),
		runPulumiUpdate(m.destroy, m.logChannel, m.eventChannel),
		watchForEvents(m.eventChannel),
		spinner.Tick,
	)
}

// Update acts on any events and updates state (model) accordingly
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case events.EngineEvent:
		if msg.ResourcePreEvent != nil {
			m.updatesInProgress[msg.ResourcePreEvent.Metadata.URN] = msg.ResourcePreEvent.Metadata.Type
		}
		if msg.ResOutputsEvent != nil {
			urn := msg.ResOutputsEvent.Metadata.URN
			m.updatesComplete[urn] = msg.ResOutputsEvent.Metadata.Type
			delete(m.updatesInProgress, urn)
		}
		return m, watchForEvents(m.eventChannel) // wait for next event
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case logMessage:
		if msg.msg == "Success" {
			m.currentMessage = "Succeeded!"
			return m, tea.Quit
		}
		if isUrl(msg.msg) {
			m.currentMessage = fmt.Sprintf("Succeeded! %s", msg.msg)
			return m, tea.Quit
		}
		m.currentMessage = msg.msg
		return m, watchForLogMessages(m.logChannel)
	default:
		return m, nil
	}
}

// View displays the state in the terminal
func (m model) View() string {
	var inProgVals []string
	var completedVals []string
	doc := strings.Builder{}
	if len(m.updatesInProgress) > 0 || len(m.updatesComplete) > 0 {
		for _, v := range m.updatesInProgress {
			inProgVals = append(inProgVals, listItem(v))
		}
		sort.Strings(inProgVals)
		for _, v := range m.updatesComplete {
			completedVals = append(completedVals, listDone(v))
		}
		sort.Strings(completedVals)

		inProgVals = append([]string{listHeader("Updates in progress")}, inProgVals...)
		completedVals = append([]string{listHeader("Updates completed")}, completedVals...)
		lists := lipgloss.JoinHorizontal(lipgloss.Top,
			list.Render(
				lipgloss.JoinVertical(lipgloss.Left,
					inProgVals...,
				),
			),
			list.Copy().Width(columnWidth).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					completedVals...,
				),
			),
		)
		doc.WriteString("\n")
		doc.WriteString(lists)
	}

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	s := fmt.Sprintf("\n%sCurrent step: %s%s\n", m.spinner.View(), m.currentMessage, docStyle.Render(doc.String()))
	if m.quitting {
		s += "\n"
	}
	return s
}

func main() {
	kingpin.Version("0.0.1")

	var destroy bool

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Register user
	case deployCmd.FullCommand():
		destroy = false

	// Post message
	case destroyCmd.FullCommand():
		destroy = true
		if *name == "" {
			app.FatalUsage("Must specify a name for destroys")
		}
	}

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	p := tea.NewProgram(model{
		logChannel:        make(chan logMessage),
		eventChannel:      make(chan events.EngineEvent),
		destroy:           destroy,
		spinner:           s,
		updatesInProgress: map[string]string{},
		updatesComplete:   map[string]string{},
	})

	fmt.Printf("your application has been deployed!")

	if p.Start() != nil {
		app.Fatalf("could not start program")
	}
}

func isUrl(url string) bool {
	_, err := neturl.ParseRequestURI(url)
	if err != nil {
		return false
	}

	u, err := neturl.Parse(url)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
