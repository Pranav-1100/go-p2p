# Go P2P Chat
**Production-Grade Peer-to-Peer Chat Application with Complete DevSecOps CI/CD Pipeline**

---

## 📌 Project Overview

Go P2P Chat is a lightweight, peer-to-peer chat application built with Go that enables direct communication between two users without requiring a central server. This project demonstrates production-grade DevOps practices with a complete CI/CD pipeline implementing security at every stage.

**Key Highlights:**
- Pure P2P architecture with TCP connections
- Zero external dependencies (Go stdlib only)
- Complete DevSecOps pipeline with 4 security layers
- Multi-stage Docker build (~10-15MB final image)
- Kubernetes deployment with health checks
- 7 passing unit tests with 26.9% coverage

---

## 🎯 Problem & Solution

**Problem:**
Traditional chat applications rely on centralized servers, creating:
- Single points of failure
- Scalability bottlenecks
- Privacy concerns (all messages through central authority)
- Network latency issues
- High infrastructure costs

**Solution:**
A peer-to-peer chat system with:
- Direct connections between peers
- No central server dependency
- Simple deployment (single binary)
- Production-ready CI/CD
- Comprehensive security scanning

---

## 🏗️ Architecture

```
┌──────────────────┐                           ┌──────────────────┐
│   Peer A         │                           │   Peer B         │
│   (Listener)     │    TCP Connection         │   (Dialer)       │
│                  │◄─────────────────────────►│                  │
│  ./go-p2p        │    Direct, No Server      │  ./go-p2p        │
│   -listen :8080  │                           │   -connect       │
│                  │                           │   localhost:8080 │
│                  │                           │                  │
│  ┌────────────┐  │                           │  ┌────────────┐  │
│  │ Goroutine  │──┼───── Messages ───────────┼─►│ Goroutine  │  │
│  │  (Stdin)   │  │                           │  │  (Network) │  │
│  └────────────┘  │                           │  └────────────┘  │
│                  │                           │                  │
│  ┌────────────┐  │                           │  ┌────────────┐  │
│  │ Goroutine  │◄─┼───── Messages ───────────┼──│ Goroutine  │  │
│  │  (Network) │  │                           │  │  (Stdin)   │  │
│  └────────────┘  │                           │  └────────────┘  │
└──────────────────┘                           └──────────────────┘
```

**How It Works:**
1. **Connection**: One peer listens (host), other dials (guest)
2. **Bidirectional Communication**: Two goroutines per peer
   - Goroutine 1: Reads from keyboard, sends to network
   - Goroutine 2: Reads from network, displays to user
3. **No Blocking**: Concurrent send/receive using Go channels

---

## 🛠️ Technology Stack

| Component | Technology | Why? |
|-----------|-----------|------|
| Language | Go 1.22 | Excellent concurrency, single binary, fast |
| Networking | TCP (net package) | Reliable, bidirectional, built-in |
| Concurrency | Goroutines | Lightweight threads, perfect for I/O |
| Container | Docker (Multi-stage) | Minimal image size, secure |
| Orchestration | Kubernetes | Auto-scaling, health checks |
| CI/CD | GitHub Actions | Free, powerful, native integration |
| Security | CodeQL, Trivy, OWASP ZAP | Industry-standard tools |
| Testing | Go testing package | Built-in, no dependencies |

---

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- Docker (optional)
- kubectl (for Kubernetes)

### Run Locally

**1. Clone & Build:**
```bash
git clone https://github.com/pranav-1100/go-p2p.git
cd go-p2p
go build -o go-p2p
```

**2. Test:**
```bash
go test ./... -v
```

**3. Start Chatting:**

Terminal 1 (Host):
```bash
./go-p2p -listen :8080
```

Terminal 2 (Guest):
```bash
./go-p2p -connect localhost:8080
```

Type messages and press Enter. Type `exit` to quit.

---

## 🐳 Docker Deployment

### Build & Run

```bash
# Build image
docker build -t go-p2p:latest .

# Run container (CI mode with health check)
docker run -d -p 8081:8081 --name go-p2p-ci go-p2p:latest

# Test health endpoint
curl http://localhost:8081/health
# Expected: {"status":"UP"}

# Cleanup
docker rm -f go-p2p-ci
```

**Dockerfile Features:**
- Multi-stage build (builder + runtime)
- Builder: golang:1.22-alpine
- Runtime: alpine:3.19 (minimal)
- Final image: ~10-15MB
- Static binary (CGO_ENABLED=0)
- Health check configured

---

## ☸️ Kubernetes Deployment

### Deploy to K8s

```bash
# Setup local cluster (Kind)
kind create cluster --name go-p2p-cluster

# Build & push image
docker build -t your-username/go-p2p:latest .
docker push your-username/go-p2p:latest

# Update k8s/deployment.yaml with your image name

# Deploy
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Verify
kubectl get pods
kubectl get svc

# Test
kubectl port-forward svc/p2p-service 8081:80
curl http://localhost:8081/health

# Cleanup
kubectl delete -f k8s/
kind delete cluster --name go-p2p-cluster
```

---

## 🔄 CI/CD Pipeline

### CI Pipeline (15 Stages)

**File:** `.github/workflows/ci.yml`

**Triggers:** Push to main, Pull requests, Manual

| Stage | Purpose | Fail Criteria |
|-------|---------|---------------|
| Checkout | Clone code | - |
| Setup Go | Install Go 1.22 | - |
| Linting | Code quality check | Any linting error |
| Unit Tests | Run 7 tests | Any test failure |
| CodeQL (SAST) | Security scan | HIGH/CRITICAL findings |
| Trivy FS (SCA) | Dependency scan | HIGH/CRITICAL CVEs |
| Docker Build | Create image | Build failure |
| Trivy Image | Container scan | HIGH/CRITICAL CVEs |
| Smoke Test | Test container | Health check fails |
| Push DockerHub | Publish image | Push failure |

**Key Features:**
- Fail-fast pipeline
- Job dependencies (sequential safety)
- Dual tagging (latest + SHA)
- Security gates enforced

### CD Pipeline (7 Stages)

**File:** `.github/workflows/cd.yml`

**Trigger:** Only runs after successful CI

| Stage | Purpose |
|-------|---------|
| Checkout | Clone code |
| Create Kind Cluster | Local K8s |
| Deploy to K8s | Apply manifests |
| Wait for Ready | Pod health check |
| DAST Scan | OWASP ZAP |
| Cleanup | Delete cluster |

---

## 🔒 Security (4 Layers)

**1. SAST (Static Analysis) - CodeQL**
- Scans: Source code
- Detects: SQL injection, XSS, buffer overflows, race conditions
- When: Every push/PR

**2. SCA (Dependency Analysis) - Trivy FS**
- Scans: go.mod dependencies
- Detects: Known CVEs in packages
- When: Before Docker build

**3. Container Scan - Trivy Image**
- Scans: Docker image layers
- Detects: OS vulnerabilities, library CVEs
- When: After Docker build

**4. DAST (Runtime Analysis) - OWASP ZAP**
- Scans: Running application
- Detects: Runtime vulnerabilities, port issues
- When: After K8s deployment

**Security Best Practices:**
- No hardcoded secrets
- Multi-stage builds (no build tools in final image)
- Minimal base image (Alpine)
- Static binary (no dynamic libraries)
- Fail on HIGH/CRITICAL findings

---

## 🔑 GitHub Secrets Setup

**REQUIRED:** Configure these secrets in your GitHub repo

**Step 1: Get DockerHub Token**
1. Login to hub.docker.com
2. Account Settings → Security → New Access Token
3. Copy token (only shown once!)

**Step 2: Add to GitHub**
1. Go to your repo → Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Add:
   - Name: `DOCKERHUB_USERNAME`, Value: your-dockerhub-username
   - Name: `DOCKERHUB_TOKEN`, Value: paste-token-here

---

## 🧪 Testing

### Test Coverage: 26.9%

| Test File | Test Cases | Coverage |
|-----------|-----------|----------|
| connection_test.go | 5 tests | 100% |
| main_test.go | 2 tests | 50% |

**Run Tests:**
```bash
# All tests
go test ./... -v

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Specific test
go test -v -run TestStartListener
```

**7 Test Cases:**
1. TestStartListener - Listener accepts connection
2. TestStartListenerInvalidAddress - Error handling
3. TestDialPeer - Dialer connects successfully
4. TestDialPeerConnectionRefused - Connection refused
5. TestDialPeerInvalidAddress - Invalid address
6. TestHealthEndpoint - Health check returns 200
7. TestExitCommand - Exit/quit detection

---

## 📖 API Documentation

### Health Check Endpoint

**Endpoint:** `GET /health`
**Port:** 8081 (CI mode)

**Request:**
```bash
curl http://localhost:8081/health
```

**Response:** 200 OK
```json
{"status":"UP"}
```

**Purpose:**
- Kubernetes liveness probes
- Container health monitoring
- Auto-restart unhealthy pods

---

## ⚙️ Configuration

### Command-Line Flags

| Flag | Description | Example |
|------|-------------|---------|
| `-listen` | Start as listener | `./go-p2p -listen :8080` |
| `-connect` | Start as dialer | `./go-p2p -connect localhost:8080` |
| `-ci` | CI mode (non-interactive) | `./go-p2p -ci` |

### Linting Config (.golangci.yml)

```yaml
run:
  timeout: 3m
linters:
  enable:
    - govet        # Correctness
    - errcheck     # Error handling
    - staticcheck  # Bug detection
    - ineffassign  # Dead code
```

---

## 📁 Project Structure

```
go-p2p/
├── .github/workflows/
│   ├── ci.yml              # CI pipeline (15 stages)
│   └── cd.yml              # CD pipeline (7 stages)
├── k8s/
│   ├── deployment.yaml     # K8s deployment
│   └── service.yaml        # K8s service
├── main.go                 # Entry point (49 lines)
├── chat.go                 # Chat handler (65 lines)
├── connection.go           # Network layer (33 lines)
├── health.go               # Health endpoint (17 lines)
├── main_test.go            # Tests (33 lines)
├── connection_test.go      # Connection tests (91 lines)
├── Dockerfile              # Multi-stage build
├── .dockerignore           # Build optimization
├── .golangci.yml           # Linting config
├── go.mod                  # Module definition
├── go.sum                  # Dependencies
└── README.md               # This file
```

**Total Lines of Code:**
- Application: 160 lines
- Tests: 121 lines
- Total: 281 lines

---

## 🔧 Troubleshooting

### Issue: Tests Fail with "connection refused"

**Solution:**
```bash
# Find process using port
lsof -i :9999
kill -9 <PID>
```

### Issue: Docker build fails "go.sum not found"

**Solution:**
```bash
go mod tidy
docker build -t go-p2p:latest .
```

### Issue: Container health check fails

**Solution:**
```bash
docker logs go-p2p-ci
docker exec go-p2p-ci wget -qO- http://localhost:8081/health
```

### Issue: CI pipeline fails on Trivy scan

**Solution:**
```bash
# Update dependencies
go get -u ./...
go mod tidy

# Rebuild
docker build -t test .
trivy image test
```

---

## 📈 Results & Observations

### CI Pipeline Results

✅ **All Security Gates Passing**
- CodeQL: No vulnerabilities
- Trivy FS: No HIGH/CRITICAL
- Trivy Image: No HIGH/CRITICAL
- OWASP ZAP: Baseline scan clean

✅ **Build Quality**
- All 7 tests passing
- Coverage: 26.9%
- No linting issues
- Build time: ~2 minutes

✅ **Image Metrics**
- Final size: ~10-15 MB
- Base: Alpine 3.19
- Build stages: 2

### CD Pipeline Results

✅ **Kubernetes Deployment**
- Pod running and healthy
- Health checks passing
- Service accessible

✅ **DAST Scan**
- OWASP ZAP scan complete
- No HIGH findings

### Key Observations

**What Worked Well:**
1. Go's simplicity - no external dependencies
2. Goroutines made concurrent I/O trivial
3. Multi-stage build reduced image 95%
4. Security gates caught issues early
5. Kind enabled local K8s testing

**Challenges Overcome:**
1. Testing goroutines required careful sync
2. Added CI mode for container orchestration
3. Multi-stage build solved image bloat

---

## ⚠️ Limitations & Future Work

### Current Limitations

1. Network Scope: localhost/LAN only (no NAT traversal)
2. One-to-One: No group chat
3. No Encryption: Plaintext TCP
4. No Persistence: Messages not saved
5. No Authentication: No user verification
6. Test Coverage: 26.9% (target: 70%+)

### Future Enhancements

**Short-term:**
- Add TLS encryption
- Implement STUN/TURN for NAT
- Add Prometheus metrics
- Increase test coverage

**Long-term:**
- Group chat support
- Message persistence (database)
- User authentication (JWT)
- End-to-end encryption
- WebRTC for browser support
- Helm chart for production

---

## 👨‍💻 Author

**Pranav Aggarwal**

**Project Type:** Individual DevOps CI/CD Implementation
**Focus:** Production-grade Go application with complete DevSecOps pipeline

---

## 🎓 What This Project Demonstrates

### Go Programming
- Concurrency with goroutines and channels
- Network programming with TCP
- Standard library mastery
- Clean code architecture

### DevOps Practices
- CI/CD pipeline design
- Infrastructure as Code
- Container orchestration
- GitOps workflows

### DevSecOps
- Shift-left security
- Multi-layer scanning (SAST, SCA, Container, DAST)
- Secure supply chain
- Vulnerability management

### Cloud-Native Development
- Docker containerization
- Multi-stage builds
- Kubernetes patterns
- Health checks and probes

