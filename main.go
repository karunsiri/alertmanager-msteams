package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// AlertManagerPayload represents the structure of the Alertmanager webhook payload
type AlertManagerPayload struct {
	Alerts            []Alert           `json:"alerts"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	CommonLabels      map[string]string `json:"commonLabels"`
	ExternalURL       string            `json:"externalURL"`
	GroupKey          string            `json:"groupKey"`
	GroupLabels       map[string]string `json:"groupLabels"`
	Receiver          string            `json:"receiver"`
	Status            string            `json:"status"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Version           string            `json:"version"`
}

// Alert represents individual alert details within the Alertmanager payload
type Alert struct {
	Annotations  map[string]string `json:"annotations"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
	StartsAt     string            `json:"startsAt"`
	Status       string            `json:"status"`
}

const (
	templatePath = "./default.tmpl"
	timeFormat   = "2006-01-02 15:04:05 MST"
)

var excludedLabels []string

func main() {
	excludedLabelsArg := flag.String("exclude-labels", "", "Comma-separated list of labels to exclude")
	flag.Parse()

	if *excludedLabelsArg != "" {
		excludedLabels = strings.Split(*excludedLabelsArg, ",")
	}

	http.HandleFunc("/alert", alertHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Determine if a label should be included based on the excluded labels list
func shouldInclude(label string) bool {
	for _, excludedLabel := range excludedLabels {
		if label == excludedLabel {
			return false
		}
	}
	return true
}

func formatTimestamp(ts string) string {
	parsedTime, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts // If parsing fails, return the original timestamp
	}
	return parsedTime.Format(timeFormat)
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	// Read the JSON body from Alertmanager
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON into an AlertManagerPayload struct
	var payload AlertManagerPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error reading payload: %v", err)
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return
	}

	log.Printf(
		"Receive alert: %s (%s)\n",
		payload.CommonLabels["alertname"],
		payload.Status,
	)

	// Load the template file with the custom function map
	extraFunctions := template.FuncMap{
		"shouldInclude":   shouldInclude,
		"formatTimestamp": formatTimestamp,
	}
	tmpl, err := template.New("default.tmpl").Funcs(extraFunctions).ParseFiles(templatePath)
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	// Render the template with the alert data
	var renderedTemplate bytes.Buffer
	if err := tmpl.Execute(&renderedTemplate, payload); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		return
	}

	// Send the rendered template to Microsoft Teams webhook
	webhookURL := os.Getenv("WEBHOOK_URL")

	resp, err := http.Post(webhookURL, "application/json", &renderedTemplate)
	if err != nil {
		log.Printf("Error sending to Teams webhook: %v", err)
		http.Error(w, "Unable to send to Teams webhook", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Error response from Teams webhook: %s", bodyBytes)
		http.Error(w, "Unable to send to Teams webhook", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully forwarded %s (%s)\n", payload.CommonLabels["alertname"], payload.Status)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert forwarded to Teams"))
}
