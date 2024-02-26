package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "strings"
)

var gotifyURL = "https://your-gotify-server.com"
var gotifyToken = "your-gotify-access-token"

func sendGotifyNotification(message string) {
 
  req, err := http.NewRequest(http.MethodPost, gotifyURL+"/message", strings.NewReader(message))
  if err != nil {
    fmt.Println("Error creating Gotify request:", err)
    return
  }


  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Authorization", "Bearer "+gotifyToken)


  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("Error sending Gotify notification:", err)
    return
  }

  // Handle response (optional)
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    fmt.Println("Gotify notification failed with status:", resp.StatusCode)
  } else {
    fmt.Println("Gotify notification sent successfully")
  }
}
