# Chipmunk
<p>
  <img src="https://github.com/swissymissy/chipmunk/blob/main/cmd/frontend/images/chipmunk.png" width="300" />
</p>

Chipmunk is a tool to help professors take students attendance during classtime faster, save class time, and improve record accuracy with cheating prevention layers. It is a hybrid web-application that lets professors host the webapp on their computer (localhost) and use cloudflare to create a secured tunnel (https) to let student register/check-in by scanning a QR code. 

## Features
- Rotating QR code check-in 
- GPS/Location check
- Student Login/Register system
- Device fingerprint check during check-in time
- Professor dashboard
- Attendance records Excel export: daily report, semester report

## Deployment Options
There are two ways to deploy the webapp using Cloudflare - Quick Tunnel or Named Tunnel
### Option 1: Quick Tunnel 
The program starts a local server on the professor's computer (localhost) and create a temporary public HTTPS link using Cloudflare Quick Tunnel. The URL is generated with random names each time the program starts and is only valid while the program is running, for example `https://starsmerchant-councils-sealed-impressed.trycloudflare.com`.This does not require:
- domain name
- cloudlare account
- hosting

  
_However, because the temporary URL uses the shared `trycloudflare.com` domain, some browsers may show a warning page such as "This site may be dangerous." This warning is related to the reputation of shared temporary tunnel domains, not necessarily Chipmunk itself_


Quick Tunnel is for:
- testing
- demos
- personal use
- users who wants the simplest setup
### Option 2: Named Tunnel
For real classroom use, a named Cloudflare Tunnel with a custom domain is recommended. This will avoid the browser warning caused by temporary shared tunnel URLs.
Name Tunnel needs:
- domain name
- cloudflare account to setup tunnel token


_To establish named tunnel, user should buy a domain, register with Cloudflare, create a Cloudflare account to use their dashboard. [Create a Tunnel](https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/get-started/create-remote-tunnel/)_

## Setup Instructions
### Option A: Download executable files (Windows)
If you do not have Golang installed, this is the easiest way to run the program.
The zip file should include: 
```
chipmunk.exe
setup.exe
cloudflared.exe
.env.example
```
1. First, run: `setup.exe` to generate jwt secret for server and create password for professor. Note: It is normal that the terminal does not print inputs out. It hides the inputs intentionally. After first setup, you don't need to run this file again unless you want to change password.
2. Then run `chipmunk.exe`. This is the main file to run server on localhost and create a secured tunnel with cloudflare.  

### Option B: Compile from source code:
- Requirement: [Go](https://go.dev/doc/install)
1. Clone the git repo to your computer: 
```bash
git clone https://github.com/swissymissy/chipmunk.git
```
2. Download dependencies:
```bash
go mod tidy
``` 
3. Create .env file (optional):
```bash
cp .env.example .env
```
Environment Variables:
```env
PORT="8080"
PLATFORM="prof"

# when CLOUDFLARE_TUNNEL_TOKEN is set, this must match the public hostname configured on cloudflare
# format: https://example.domain.com ( no trailing slash)
BASE_URL="http://localhost"

DB_URL="./chipmunk.db"

JWT_SECRET=""
PROFESSOR_PASSWORD_HASH=''

# this is for named tunnel token when you setup a named tunnel on Cloudflare dashboard
# if you don't, just leave it empty
CLOUDFLARE_TUNNEL_TOKEN=""
```
4. Compile setup binary file:
```bash
go build -o setup ./cmd/setup
```
5. Compile chipmunk binary file:
```bash
go build -o chipmunk ./cmd/server
```
6. Run setup tool:
```bash
./setup
```
7. Run the chipmunk server:
```bash
./chipmunk
```
## What do Setup tool and Chipmunk server do:
1. The setup tool will:
- create .env file from .env.example 
- generate a jwt secret 
- prompt for password from user
- save in .env
2. When Chipmunk server starts, it will:
- load config from .env
- open SQLite database
- run database migration
- start local web server
- start cloudflare tunnel
- allow professor to log in and start attendance
- generate public HTTPS URL for students to check-in

## Tech Stack
- Go 
- SQLite
- HTML/CSS/JavaScript
- Goose migration
- sqlc
- Cloudflare
----
### Privacy Notice: Device Signals and Fingerprint
Chipmunk may collect limited device-related signal during attendance check-in, such as browser/device information, user agent, IP address, and optional device identifiers. 
These signals are **used only for attendance integrity purposes**, such as **helping professors detect suspicious patterns like one device being used to check in for multiple students.**
Device fingerprinting is not perfect and should not be treated as absolute proof of cheating. It is only a supporting signal for review.
**Chipmunk does not use device signals for advertising, tracking across websites, or any non-attendance purpose.**
### Intended Use
Chipmunk is intended for **adult students in higher-education or similar classroom settings.**
Do not use Chipmunk with minors or underage students unless you fully understand and comply with all applicable school policies, privacy laws, parental consent requirements, and institutional rules.
The developer is not responsible for how users deploy Chipmunk or collect student data.
### Disclaimer
Chipmunk uses multiple layers to reduce attendance cheating, including rotating QR codes, student login, device signals, GPS/location checks, and check-in windows(session).
However, no browser-based attendance system can guarantee perfect proof of physical presence.
GPS can be inaccurate or spoofed. Device fingerprints can change or be unreliable. Shared networks may make many students appear to have the same public IP address.
It is designed to increase friction, reduce manual attendance work, and flag suspicious activity for professor review. It should not be treated as a fully automated disciplinary system.
