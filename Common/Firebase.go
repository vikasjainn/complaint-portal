// Common/Firebase.go
package Common

import (
	"context"
	"encoding/json"
	"log"
	"os" // Use os package to read the file

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirestoreClient is a global client that can be used by other packages.
var FirestoreClient *firestore.Client

// InitFirebase initializes the Firebase app and the Firestore client.
func InitFirebase() {
	ctx := context.Background()

	// Read the credentials file to get the project ID
	credBytes, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Error reading credentials.json file: %v\nMake sure the file is in the project root directory.", err)
	}

	// A simple struct to unmarshal just the project_id from the credentials
	var creds struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(credBytes, &creds); err != nil {
		log.Fatalf("Error unmarshalling credentials: %v\n", err)
	}

	if creds.ProjectID == "" {
		log.Fatalf("ProjectID not found in credentials.json. The credentials file seems to be invalid.")
	}

	// Now, initialize the app with the explicit ProjectID
	config := &firebase.Config{
		ProjectID: creds.ProjectID,
	}
	sa := option.WithCredentialsFile("credentials.json")
	app, err := firebase.NewApp(ctx, config, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Firestore client: %v\n", err)
	}

	FirestoreClient = client
	log.Println("Successfully connected to Firebase Firestore.")
}
