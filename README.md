# Chipmunk
![chipmunk logo](https://github.com/swissymissy/chipmunk/blob/main/cmd/frontend/images/chipmunk.png)
Chipmunk is a hybrid student attendance web application designed to help professors save time taking in-class attendance while adding multiple layers to reduce cheating and improve attendance record accuracy.
Chipmunk is designed as a local-first application. The professor runs the server on their own computer, and students check in through a QR code during class.

## Features
- Rotating QR code check-in 
- GPS/Location check
- Student Login/Register system
- Device fingerprint check during check-in time
- Professor dashboard
- Attendance records Excel export: daily report, semester report

## Deployment Options
### Option 1: Quick Tunnel 
This is the easiest way to run and try Chipmunk.
In this mode, Chipmunk starts a local server on the professor's computer and create a temporary public HTTPS link using Cloudflare Quick Tunnel.
This mode does not require:
- domain name
- cloudlare account
- hosting
_However, because the temporary URL uses the shared `trycloudflare.com` domain, some browsers may show a warning page such as "This site may be dangerous." This warning is related to the reputation of shared temporary tunnel domains, not necessarily Chipmunk itself_
Quick Tunnel Mode is best for:
- testing
- demos
- personal use
- users who wants the simplest setup
### Option 2: Named Tunnel
For real classroom use, a named Cloudflare Tunnel with a custom domain is recommended. This will avoid the browser warning caused by temporary shared tunnel URLs.
Named Tunnel Mode requires:
- domain name
- cloudflare account to setup tunnel token

## Setup Instructions
### Option A: Download executable files (Windows)
If you do not have Golang installed, this is the easiest way to run the program.
The package should include: 
```
chipmunk.exe
setup.exe
cloudflared.exe
.env.example
```
1. First, run: `setup.exe` to generate jwt secret for server and create password for professor. After first setup, you don't need to run this file again unless you want to change password.
2. Then run `chipmunk.exe`. This is the main file to run server on localhost and create a security tunnel with cloudflare.  

### Option B: build from source code (required Golang installed) 
1. Clone the git repo to your computer: 
```bash
git clone https://github.com/swissymissy/chipmunk.git
```
2. Download dependencies: `go mod tidy` 
3. Create .env file (optional): `cp .env.example .env`
4. Compile setup binary file: `go build -o setup ./cmd/setup`
5. Compile chipmunk binary file: `go build -o chipmunk ./cmd/server`
6. Run setup to generate jwt secret and password for professor: `./setup`
7. Run the chipmunk server: `./chipmunk`
## How Setup tool and Chipmunk server work:
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
- allow prfessor to log in and start attendance
- generate public HTTPS URL for students to check-in

## Tech STack
- Go 
- SQLite
- HTML/CSS/JavaScrip

### Privacy Notice: Device Signals and Fingerprint
Chipmunk may collect limited device-related signal during attendance check-in, such as browser/device information, user agent, IP address, and optional device identifiers. 
These signals are used only for attendance intergrity purposes, such as helping professors detect suspicious patterns like one device being used to check in for multiple students.
Device fingerprinting is not perfect and should not be treated as absolute proof of cheating. It is only a supporting signal for review.
Chipmunk does not use device signals for advertising, tracking across websites, or any non-attendance purpose.
### Intended Use
Chipmunk is intended for adult students in higher-education or similar classroom settings.
Do not use Chipmunk with minors or underage students unless you fully understand and comply with all applicable school policies, privacy laws, parental consent requirements, and institutional rules.
The developer is not responsible for how users deploy Chipmunk or collect student data.
### Disclaimer
Chipmunk uses multiple layers to reduce attendance cheating, including rotating QR codes, student login, device signals, GPS/location checks, and check-in windows(session).
However, no browser-based attendance system can guarantee perfect proof of physical presence.
GPS can be inaccurate or spoofed. Device fingerprints can change or be unreliable. Shared networks may make many students appear to have the same public IP address.
Chipmunk is designed to increase friction, reduce manual attendance work, and flag suspicious activity for professor review. It should not be treated as a fully automated disciplinary system.