# Holo Checker App

This is a **Go-based monitoring tool** designed to fetch and track Hololive karaoke streams in real-time.  
It periodically queries an API for scheduled or live streams, filters them based on Hololive criteria, and triggers notifications or focus mode when changes are detected.

---

## **Features**

- **Automatic Stream Monitoring**: Checks for new Hololive karaoke streams every 10 minutes.
- **Smart Notifications**: Notifies on the first run or when new streams are found (based on `ChangeChecker`).
- **Focus Mode Scheduling**: Automatically schedules focus mode on stream updates.
- **Prometheus Metrics**: Exposes a `/metrics` endpoint on `localhost:2112` for monitoring.
- **Logging & Retry Mechanism**: Built-in logging using `logrus` and retry mechanism for API calls.
- **System Tray Integration**: Provides a system tray interface for better user interaction.

---

## **Core Workflow**

### **Monitor Function**

The `Monitor` function:

1. Creates an API client (`controller.NewAPIClient`).
2. Fetches the list of available streams with a retry mechanism (`utility.Retry`).
3. Filters streams to include only Hololive (`FilterStreams` with `IsHololive`).
4. Delegates stream change handling to `handleStreamUpdate`.

### **handleStreamUpdate**

- On **first run**: Calls `Notify` and schedules `FocusMode` (even if no streams are available).
- On **subsequent runs**: Calls `Notify` and schedules `FocusMode` only if new streams are detected.

---

## **Prometheus Integration**

- Metrics are served at `http://localhost:2112/metrics`.
- Useful for monitoring the app’s internal performance.

---

## **Installation**

### **Prerequisites**

- Go 1.20+ installed on your system.
- Internet connection (for API calls).
- Git (if cloning from repository).

### **Clone Repository**

```bash
git clone https://github.com/<your-username>/holo-checker-app.git
cd holo-checker-app
go mod tidy
```

## Environment Variables Setup

To run this project, you need to set up a `.env` file in the root directory with the following variables:

```env
TELEGRAM_BOT_TOKEN=xxxx:xxxxx-xxx-xxxx
TELEGRAM_CHAT_ID=625020000
WHATSAPP_PHONE_NUMBER=000000
WHATSAPP_API_KEY=00000
XAPIKEY=w2w2w2w2-169f-q1q1q1-xxxx-asasad
```

**Instructions:**

1. Create a new file named `.env` in the root folder of your project.
2. Copy and paste the example above into your `.env` file.
3. Replace the placeholder values with your actual credentials if you have them.

**Note:**  
Do not share your real credentials publicly. The above values are examples only.

## Development

To run the application in development mode:

```sh
go run ./cmd/main.go
```

## Production Build (Windows GUI, no console)

To build the application for Windows as a GUI app (no console window):

```sh
rsrc -ico favicon.ico -o ./cmd/rsrc.syso
go build -ldflags="-H=windowsgui" -o holo-checker-app.exe ./cmd
```

## Normal Build (with console, for debugging)

To build the application with a console window (useful for debugging):

```sh
go build -o holo-checker-app.exe ./cmd/main.go
```

---

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements.

## Acknowledgements

Thank you for checking out this project! If you find it  useful, please star the repository and share it with others
