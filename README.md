# TokenVault

**TokenVault** is a developer experience (DevEx) utility designed to automate local authentication workflows. It functions as a centralized, local token store that captures authentication tokens (such as JWTs) from running backend services and automatically injects them into API testing clients like Postman.

This tool eliminates the manual process of copying and pasting access tokens between browser/terminal and API clients during development, ensuring that test requests always use the most recent valid credentials.

---

## Architecture

The system operates as a local **sidecar service** consisting of three main components:

### The Vault (Server)
A lightweight Go-based REST API that runs locally (default port: **9999**) and persists data to a SQLite database.

### The Client Integration
A minimal HTTP hook implemented in your backend projects (Node.js, Go, Python, etc.) that pushes fresh tokens to the Vault upon successful login.

### The Automation Script
A Postman pre-request script that queries the Vault based on the target port or domain and injects the authorization header dynamically.

---

## Features

- **Centralized Storage**  
  Persists tokens locally using SQLite, ensuring data remains available across system restarts.

- **Port-Based Context Switching**  
  Automatically detects which project is being tested based on the target port (e.g., `5000` vs `8080`) and retrieves the corresponding token.

- **Interactive Configuration**  
  Includes a CLI utility to generate custom automation scripts based on your specific development environment.

- **Single Binary Distribution**  
  Built with Go for easy cross-platform usage without external runtime dependencies.

---

## Technical Stack

- **Language:** Go (Golang)
- **CLI Framework:** Cobra
- **Web Framework:** Gin Gonic
- **Database:** SQLite (`modernc.org/sqlite`, CGO-free implementation)

---

## Installation

### Prerequisites

- Go **1.20 or higher** (required only for building from source)

### Build from Source

Clone the repository:

```bash
git clone https://github.com/yourusername/tokenvault.git
cd tokenvault
```

Build the executable:
```bash
go build -o tokenvault .
```
(Optional) Add the binary to your system PATH for global access.


## Usage Guide
1. Start the TokenVault Server

Open a terminal window and start the local server. This process must remain running in the background.
```bash
./tokenvault serve
```

Expected output: Listening and serving HTTP on :9999

2. Configure Postman

TokenVault provides an automated setup command to generate the required JavaScript for Postman.

Run the setup command:
```bash
./tokenvault setup-postman
```
Follow the interactive prompts to map local ports to project names.

### Example mappings:

- Port 5000 → realtime-bidding
- Port 8080 → ecommerce-api
- The command generates a file named tokenvault_script.js.

## Postman Setup Steps
- Open Postman
- Go to Collection Settings
- Open the Pre-request Script tab
- Copy the contents of tokenvault_script.js
- Paste it into the editor and save


3. Integrate with Backend Projects

To store tokens, your backend application must send a POST request to TokenVault whenever a user logs in successfully.

Example: Node.js / Express

Add this immediately after generating the authentication token:
```bash
// Fire-and-forget request to TokenVault
fetch("http://localhost:9999/store", {
  method: "POST",
  headers: {
    "Content-Type": "application/json"
  },
  body: JSON.stringify({
    project: "realtime-bidding", // Must match setup-postman configuration
    token: yourGeneratedToken
  })
}).catch(() => {
  console.error("TokenVault unavailable");
});
```

4. Testing Workflow

- Perform a login action in your application (frontend or Postman).

- Verify the token is stored by checking TokenVault server logs.

- Execute any protected API request in Postman.

- The pre-request script automatically fetches and injects the latest token.


##Troubleshooting

Common Issues
1. Connection Refused / Server Not Reachable

Cause: TokenVault server is not running.
```bash
./tokenvault serve
```

2. Token Not Updating in Postman

Cause: The project name sent by the backend does not match the Postman configuration.

Solution:
Ensure the project value in backend code exactly matches the mapping defined during setup-postman.

3. Address Already in Use (Port 9999)

Cause: Another TokenVault instance or application is using port 9999.

Solution:
Stop the conflicting process or terminate the existing TokenVault instance.

4. Setup Command Error: undefined: SetupCmd

Cause: Running Go files individually instead of the full module.

Solution:

Use:
```bash 
go run . setup-postman
```

5. Database Permission Issues

Cause: TokenVault cannot create the .tokenvault directory.

Solution:
Ensure the current user has write access to their home directory.

Database location:
```bash
~/.tokenvault/token.db
```

##Project Structure
```bash
tokenvault/
├── cmd/                # (Deprecated) Legacy command structure
├── internal/
│   ├── api/            # HTTP handlers for REST API
│   └── database/       # SQLite connection and queries
├── main.go             # Application entry point and CLI registration
├── setup.go            # Postman script generation logic
├── go.mod              # Go module definitions
└── README.md           # Project documentation
```


