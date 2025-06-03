package pushover_notificationchannel

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/DazWilkin/go-gcp-pushover-notificationchannel/pushover"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const (
	device string = "cloud-monitoring"
)

var (
	pushoverClient *pushover.Client
)

func init() {
	functions.HTTP("Webhook", Webhook)

	userKey := os.Getenv(("PUSHOVER_USERKEY"))
	if userKey == "" {
		log.Fatal("Unable to get `PUSHOVER_USERKEY` from the environment")
	}

	token := os.Getenv("PUSHOVER_TOKEN")
	if token == "" {
		log.Fatal("Unable to get `PUSHOVER_TOKEN` from the environment")
	}

	pushoverClient = pushover.New(device, userKey, token)
}

// pushover is an HTTP Cloud Function.
func Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("Expect `POST` (got: `%s`)", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Printf("Expect `application/json` (got: `%s`)", contentType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Unable to read message body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Received: %s", string(data))

	body := &Body{}
	if err := json.Unmarshal(data, body); err != nil {
		log.Println("Unable to unmarshal message body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Construct Pushover Message from Incident Body
	message, err := NewMessage(body)
	if err != nil {
		log.Println("Unable to construct Pushover message from Incident body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Message:\n%+v", message)

	// Execute template to convert Message into HTML
	var buf bytes.Buffer
	t := template.Must(template.New("message").Parse(templateMessage))
	if err := t.Execute(&buf, message); err != nil {
		log.Println("Unable to construct message template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("HTML:\n%s", buf.String())

	// Make request to Pushover
	if err := pushoverClient.SendMessage(body.Incident.ID, buf.String()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
