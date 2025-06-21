# Working with IP, Ports, and Addresses: Concepts and Go Implementation 🏷️

> "Imagine the internet as a vast city. Every building (device) has an address (IP), every door (service) has a number (port), and Go is your GPS, map, and keyring—all in one!"

---

## 🧭 What Are IPs, Ports, and Addresses?

- **IP Address:** Like a street address for your device on a network. IPv4 (e.g., 192.168.1.1) and IPv6 (e.g., 2001:db8::1).
- **Port:** Like an apartment number—identifies a specific service on a device (e.g., port 80 for web servers).
- **Socket:** The combo of IP + port (e.g., 192.168.1.1:80)—the full address for network communication.

---

## 🌍 IP Address Theory

- **IPv4:** 32 bits, four numbers (0–255), e.g., 8.8.8.8 (Google DNS).
- **IPv6:** 128 bits, eight groups, e.g., 2001:4860:4860::8888 (also Google DNS).
- **Public vs Private:** Public IPs are globally unique; private IPs are for local networks (e.g., 192.168.x.x).

**Analogy:**
- IPv4 is like a city with 4 billion houses—running out of space!
- IPv6 is a city so big, every grain of sand on Earth could have its own address.

**How Go Handles IPs:**
Go abstracts away the complexity of IP addresses. When you use `net.LookupIP`, Go asks your operating system to resolve the hostname, which in turn queries DNS servers. Go then parses the response and gives you a list of IPs—no need to worry about the protocol details!

---

## 🚪 Ports and Services

- **Well-known ports:** 80 (HTTP), 443 (HTTPS), 22 (SSH), 25 (SMTP), etc.
- **Dynamic ports:** Used for temporary connections (e.g., 49152–65535).

**Analogy:**
- Ports are like doors in a building—each service (web, mail, FTP) has its own entrance.

**How Go Handles Ports:**
When you open a connection in Go (e.g., `net.Dial("tcp", "example.com:80")`), Go creates a socket, negotiates with the OS, and binds to a random local port if you don't specify one. It handles all the low-level details, so you can focus on your app logic.

---

## 🛠️ Go in Action: Lookup IPs

Let’s see how to resolve a hostname to its IP addresses in Go:

```go
package main
import (
    "fmt"
    "net"
)
func main() {
    host := "google.com"
    ips, err := net.LookupIP(host)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("IP addresses for", host, ":")
    for _, ip := range ips {
        fmt.Println(" -", ip)
    }
}
```

[Ejercicio: Lookup IP](../../exercises/part2/04-lookup-ip/main.go)

---

## 🛠️ Go in Action: Parsing and Using Ports

```go
package main
import (
    "fmt"
    "net"
)
func main() {
    addr := "192.168.1.10:8080"
    host, port, err := net.SplitHostPort(addr)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Host:", host)
    fmt.Println("Port:", port)
}
```

[Ejercicio: Split Host and Port](../../exercises/part2/04-split-host-port/main.go)

---

## 🛠️ Go in Action: Checking Local IPs

```go
package main
import (
    "fmt"
    "net"
)
func main() {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Local IP addresses:")
    for _, addr := range addrs {
        fmt.Println(" -", addr.String())
    }
}
```

[Ejercicio: List Local IPs](../../exercises/part2/04-list-local-ips/main.go)

---

## 🛠️ Go in Action: Simple TCP Port Scanner

```go
package main
import (
    "fmt"
    "net"
    "time"
)
func main() {
    host := "scanme.nmap.org"
    ports := []int{22, 80, 443, 8080}
    for _, port := range ports {
        address := fmt.Sprintf("%s:%d", host, port)
        conn, err := net.DialTimeout("tcp", address, 2*time.Second)
        if err != nil {
            fmt.Printf("Port %d closed\n", port)
            continue
        }
        fmt.Printf("Port %d open!\n", port)
        conn.Close()
    }
}
```

[Ejercicio: Simple Port Scanner](../../exercises/part2/04-port-scanner/main.go)

---

## 🛠️ Go in Action: Custom Address Parsing

```go
package main
import (
    "fmt"
    "net"
)
func main() {
    addr := "[2001:db8::1]:443"
    host, port, err := net.SplitHostPort(addr)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("IPv6 Host:", host)
    fmt.Println("Port:", port)
    ip := net.ParseIP(host)
    if ip != nil && ip.To16() != nil {
        fmt.Println("It's a valid IPv6 address!")
    }
}
```

[Ejercicio: IPv6 Address Parsing](../../exercises/part2/04-ipv6-parse/main.go)

---

## 🎨 Visual Summary

```
[Hostname] --DNS Lookup--> [IP Address]
     |
 [Port] <--- Service (HTTP, SSH, etc.)
     |
 [Socket] = [IP:Port]
```

---

## 🤩 Fun Facts & Go Memes
- The first IP address ever assigned was 0.0.0.0 (reserved for special use).
- IPv6 has enough addresses for every atom on Earth (and then some!).
- In Go, you can build a DNS resolver, a port scanner, or even your own network protocol with just a few lines of code.
- Ports below 1024 are called "privileged"—only root/admin can use them on most systems.
- Go’s networking is so fast, you can scan thousands of ports in seconds!
- Go’s `net` package is cross-platform: your code works on Windows, Linux, and Mac without changes.
- Go handles all the low-level socket magic for you—no need to mess with C or system calls!

---

[Previous: Go Networking Packages Overview](03-go-networking-packages-overview.md) | [Next: TCP in Depth: Protocol Theory and Go Implementation](05-tcp-in-depth-protocol-theory-and-go-implementation.md)
