package usecases

import (
	"bytes"
	"text/template"

	"github.com/proelbtn/vnet/pkg/entities"
)

const blueprint = `graph {
	graph [ layout = "neato" ]
	
	{
		node [ shape = "circle", fontsize = 16, margin = 0 ]
		{{- range .containers }}
		"{{ . }}";
		{{- end }}
	}
	
	{
		node [ shape = "diamond", fontsize = 8, margin = 0, style = "filled", color = "#111827", fontcolor = "#111827", fillcolor = "#D1D5DB" ]
		{{- range .networks }}
		"{{ . }}";
		{{- end }}
	}

	{
		edge [ fontsize = 8, color = "#D1D5DBC0" ]
		{{- range .links }}
		"{{ .from }}" -- "{{ .to }}" [ taillabel = "{{.portName}}" ];
		{{- end }}
	}
}`

func renderTopology(lab *entities.Laboratory) (string, error) {
	var buf bytes.Buffer
	tmpl := template.Must(template.New("topology").Parse(blueprint))

	data := map[string]interface{}{
		"containers": getContainers(lab.Containers),
		"networks":   getNetworks(lab.Networks),
		"links":      getLinks(lab.Containers),
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getContainers(containers []*entities.Container) []string {
	data := []string{}
	for _, container := range containers {
		data = append(data, container.Name)
	}
	return data
}

func getNetworks(networks []*entities.Network) []string {
	data := []string{}
	for _, network := range networks {
		data = append(data, network.Name)
	}
	return data
}

func getLinks(containers []*entities.Container) []map[string]string {
	data := []map[string]string{}
	for _, container := range containers {
		for _, port := range container.Ports {
			data = append(data, map[string]string{
				"from":     container.Name,
				"to":       port.Network.Name,
				"portName": port.Name,
			})
		}
	}
	return data
}
