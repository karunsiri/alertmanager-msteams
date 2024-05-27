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
)

var excludedLabels []string
var webhookURL string

func main() {
	excludedLabelsArg := flag.String("exclude-labels", "", "Comma-separated list of labels to exclude")
	flag.Parse()

	if *excludedLabelsArg != "" {
		trimmedLabels := strings.Trim(*excludedLabelsArg, `"'`)
		excludedLabels = strings.Split(trimmedLabels, ",")
		log.Println("Exclude Labels:", excludedLabels)
	}
	webhookURL = os.Getenv("WEBHOOK_URL")

	http.HandleFunc("/alert", alertHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
		"formatTimestamp": formatTimestamp,
		"shouldInclude":   shouldInclude,
		"titleize":        titleize,
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
