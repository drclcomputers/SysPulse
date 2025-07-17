package plugins

import (
	"fmt"
	"os/exec"
	"strings"
	"syspulse/internal/utils"
	"time"

	"github.com/rivo/tview"
)

type DockerPlugin struct {
	config     PluginConfig
	containers []DockerContainer
	images     []DockerImage
	stats      DockerStats
}

type DockerContainer struct {
	ID      string
	Name    string
	Image   string
	Status  string
	Ports   string
	Created string
}

type DockerImage struct {
	ID      string
	Repo    string
	Tag     string
	Size    string
	Created string
}

type DockerStats struct {
	ContainersRunning int
	ContainersStopped int
	Images            int
	SystemDFUsage     string
}

func NewDockerPlugin() *DockerPlugin {
	return &DockerPlugin{
		containers: make([]DockerContainer, 0),
		images:     make([]DockerImage, 0),
	}
}

func (p *DockerPlugin) Name() string {
	return "Docker Monitor"
}

func (p *DockerPlugin) Version() string {
	return "1.0.0"
}

func (p *DockerPlugin) Description() string {
	return "Monitor Docker containers, images, and system statistics"
}

func (p *DockerPlugin) Author() string {
	return "drclcomputers @ SysPulse"
}

func (p *DockerPlugin) Initialize(config PluginConfig) error {
	p.config = config

	if !p.isDockerAvailable() {
		return fmt.Errorf("Docker is not available or not running")
	}

	return nil
}

func (p *DockerPlugin) Shutdown() error {
	return nil
}

func (p *DockerPlugin) CreateWidget() (tview.Primitive, error) {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle(p.config.Layout.Title)
	textView.SetDynamicColors(true)
	textView.SetScrollable(true)

	if p.config.Layout.BorderColor != "" {
		textView.SetBorderColor(utils.GetColorFromName(p.config.Layout.BorderColor))
	}

	if p.config.Layout.ForegroundColor != "" {
		textView.SetTitleColor(utils.GetColorFromName(p.config.Layout.ForegroundColor))
	}

	if err := p.UpdateWidget(textView); err != nil {
		return nil, err
	}

	return textView, nil
}

func (p *DockerPlugin) UpdateWidget(widget tview.Primitive) error {
	textView, ok := widget.(*tview.TextView)
	if !ok {
		return fmt.Errorf("widget is not a TextView")
	}

	if err := p.updateDockerData(); err != nil {
		textView.SetText(fmt.Sprintf("[red]Error updating Docker data: %v[white]", err))
		return err
	}

	content := p.buildDisplayContent()
	textView.SetText(content)

	return nil
}

func (p *DockerPlugin) GetWidgetConfig() WidgetConfig {
	return p.config.Layout
}

func (p *DockerPlugin) CollectData() (map[string]interface{}, error) {
	if err := p.updateDockerData(); err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"containers":   p.containers,
		"images":       p.images,
		"stats":        p.stats,
		"last_updated": time.Now(),
	}

	return data, nil
}

func (p *DockerPlugin) ExportData() map[string]interface{} {
	data, _ := p.CollectData()
	return data
}

func (p *DockerPlugin) UpdateInterval() time.Duration {
	return 5 * time.Second
}

func (p *DockerPlugin) isDockerAvailable() bool {
	cmd := exec.Command("docker", "version", "--format", "{{.Client.Version}}")
	return cmd.Run() == nil
}

func (p *DockerPlugin) updateDockerData() error {
	if err := p.updateContainers(); err != nil {
		return err
	}

	if err := p.updateImages(); err != nil {
		return err
	}

	if err := p.updateStats(); err != nil {
		return err
	}

	return nil
}

func (p *DockerPlugin) updateContainers() error {
	cmd := exec.Command("docker", "ps", "-a", "--format", "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}\t{{.CreatedAt}}")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\n")
	p.containers = make([]DockerContainer, 0)

	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}

		parts := strings.Split(lines[i], "\t")
		if len(parts) >= 6 {
			container := DockerContainer{
				ID:      parts[0],
				Name:    parts[1],
				Image:   parts[2],
				Status:  parts[3],
				Ports:   parts[4],
				Created: parts[5],
			}
			p.containers = append(p.containers, container)
		}
	}

	return nil
}

func (p *DockerPlugin) updateImages() error {
	cmd := exec.Command("docker", "images", "--format", "table {{.ID}}\t{{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\n")
	p.images = make([]DockerImage, 0)

	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}

		parts := strings.Split(lines[i], "\t")
		if len(parts) >= 5 {
			image := DockerImage{
				ID:      parts[0],
				Repo:    parts[1],
				Tag:     parts[2],
				Size:    parts[3],
				Created: parts[4],
			}
			p.images = append(p.images, image)
		}
	}

	return nil
}

func (p *DockerPlugin) updateStats() error {
	runningCount := 0
	stoppedCount := 0

	for _, container := range p.containers {
		if strings.Contains(container.Status, "Up") {
			runningCount++
		} else {
			stoppedCount++
		}
	}

	cmd := exec.Command("docker", "system", "df", "--format", "table {{.Type}}\t{{.Total}}\t{{.Active}}\t{{.Size}}")
	output, err := cmd.Output()
	dfUsage := "N/A"
	if err == nil {
		dfUsage = string(output)
	}

	p.stats = DockerStats{
		ContainersRunning: runningCount,
		ContainersStopped: stoppedCount,
		Images:            len(p.images),
		SystemDFUsage:     dfUsage,
	}

	return nil
}

func (p *DockerPlugin) buildDisplayContent() string {
	content := "[green]Docker Monitor[white]\n\n"

	content += "[yellow]Statistics:[white]\n"
	content += fmt.Sprintf("Running Containers: [green]%d[white]\n", p.stats.ContainersRunning)
	content += fmt.Sprintf("Stopped Containers: [red]%d[white]\n", p.stats.ContainersStopped)
	content += fmt.Sprintf("Images: [blue]%d[white]\n\n", p.stats.Images)

	content += "[yellow]Recent Containers:[white]\n"
	containerCount := len(p.containers)
	if containerCount > 5 {
		containerCount = 5
	}

	for i := 0; i < containerCount; i++ {
		container := p.containers[i]
		statusColor := "[red]"
		if strings.Contains(container.Status, "Up") {
			statusColor = "[green]"
		}

		content += fmt.Sprintf("• %s%s[white] (%s)\n", statusColor, container.Name, container.Image)
	}

	if len(p.containers) > 5 {
		content += fmt.Sprintf("... and %d more\n", len(p.containers)-5)
	}

	content += "\n[yellow]Recent Images:[white]\n"
	imageCount := len(p.images)
	if imageCount > 5 {
		imageCount = 5
	}

	for i := 0; i < imageCount; i++ {
		image := p.images[i]
		content += fmt.Sprintf("• [blue]%s:%s[white] (%s)\n", image.Repo, image.Tag, image.Size)
	}

	if len(p.images) > 5 {
		content += fmt.Sprintf("... and %d more\n", len(p.images)-5)
	}

	return content
}
