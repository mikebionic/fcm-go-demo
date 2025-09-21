# FCM Go Demo

This project is a simple demonstration of how to integrate and use **Firebase Cloud Messaging (FCM)** with a Go backend and a web client.  
It was originally prepared as an example for **Google I/O Ashgabat** by Google Developer Group.

> See the slides from [presentation](https://docs.google.com/presentation/d/1jSfSuZX1c6LY6yHQa0_5e2mQhEn3M1xaKuyKUYnFFeQ/edit?slide=id.g38137ad3bcf_0_344#slide=id.g38137ad3bcf_0_344)

The backend (written in Go) exposes an API endpoint to send push notifications to devices using FCM.  
The frontend (HTML + JS) registers for notifications, obtains the FCM device token, and listens for messages.

## Why use this project?
- Learn how to set up FCM with a real backend.
- Understand how to generate and use device tokens.
- Explore how notifications can be sent from Go to web clients.
- Use it as a starting point for adding push notifications in your own apps.

## Run locally
1. Place your Firebase **serviceAccountKey.json** inside the `secrets/` folder.  
2. Start the Go server: `go run main.go`
3. Open [http://localhost:8000](http://localhost:8000) in your browser.
4. Allow notifications and test sending messages.

## Contribute

If you like this project:

* ⭐ Star this repository
* 🍴 Fork it and use it in your next app
* Pull requests and improvements are always welcome

---

A small but practical example to get started with **Firebase Cloud Messaging + Go** 🚀
