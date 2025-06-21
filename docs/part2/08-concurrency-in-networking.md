# Concurrency in Networking: Theory and Go Implementation 🧵

> "Imagine a busy restaurant kitchen: chefs (goroutines) work on different dishes at the same time, waiters (channels) deliver orders and plates, and the manager (main function) keeps everything running smoothly. Welcome to the world of concurrency in Go!"

---

## 🚦 Why Concurrency Matters in Networking

Networking is all about handling many things at once: multiple clients, simultaneous requests, and real-time data. Without concurrency, your server would be like a single cashier at a supermarket—everyone waits in line, and things get slow fast!

- **Analogy:** Concurrency is like having many hands to juggle multiple balls at once.
- **In Go:** Goroutines and channels make concurrency easy, safe, and fun.

---

## 🧠 What is Concurrency? What is Parallelism?

- **Concurrency:** Doing many things at once (not necessarily at the same time, but making progress on all).
- **Parallelism:** Actually running things at the same time (on multiple CPU cores).
- **Go’s Superpower:** Goroutines are so lightweight, you can have thousands running without breaking a sweat!

**Fun Fact:**
- Go’s concurrency model is inspired by Tony Hoare’s CSP (Communicating Sequential Processes).

---

## 🧵 Goroutines: Lightweight Threads

- **What are they?** Functions that run independently, started with the `go` keyword.
- **How lightweight?** Each goroutine uses only a few KB of memory.
- **Why use them?** Handle many clients, background tasks, or timers without blocking your main program.

**Example:**

```go
package main
import (
    "fmt"
    "time"
)

func sayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    go sayHello("Alice") // Runs in a new goroutine
    go sayHello("Bob")
    time.Sleep(1 * time.Second) // Wait for goroutines to finish
    fmt.Println("Done!")
}
```

[Ejercicio: Goroutines Básicas](../../exercises/part2/08-goroutines-basic/main.go)

---

## 📬 Channels: Safe Communication Between Goroutines

- **What are they?** Typed pipes for sending data between goroutines.
- **Why use them?** Avoid race conditions and share data safely.

**Example:**

```go
package main
import "fmt"

func main() {
    ch := make(chan string) // Create a channel
    go func() {
        ch <- "Hello from goroutine!" // Send data
    }()
    msg := <-ch // Receive data
    fmt.Println(msg)
}
```

[Ejercicio: Canales Básicos](../../exercises/part2/08-channels-basic/main.go)

---

## 🛠️ Go in Action: Concurrent TCP Server

Let’s build a TCP server that handles each client in a separate goroutine.

```go
package main
import (
    "fmt"
    "net"
)

func handleConn(c net.Conn) {
    fmt.Fprintln(c, "Welcome to the concurrent server!")
    c.Close()
}

func main() {
    ln, _ := net.Listen("tcp", ":8082")
    fmt.Println("Server listening on :8082")
    for {
        conn, _ := ln.Accept()
        go handleConn(conn) // Each client handled concurrently
    }
}
```

[Ejercicio: Concurrent TCP Server](../../exercises/part2/08-tcp-concurrent-server/main.go)

---

## 🛠️ Go in Action: Channel-Based Message Broadcast

A simple example where one goroutine sends messages to many receivers using channels.

```go
package main
import (
    "fmt"
    "time"
)

func broadcaster(ch chan string) {
    for i := 1; i <= 3; i++ {
        ch <- fmt.Sprintf("Message %d", i)
        time.Sleep(500 * time.Millisecond)
    }
    close(ch)
}

func main() {
    ch := make(chan string)
    go broadcaster(ch)
    for msg := range ch {
        fmt.Println("Received:", msg)
    }
}
```

[Ejercicio: Channel Broadcast](../../exercises/part2/08-channel-broadcast/main.go)

---

## 🧩 How Go Makes Concurrency Safe and Easy

- **Scheduler:** Go’s runtime schedules goroutines efficiently across CPU cores.
- **No manual threads:** No need to manage OS threads or locks for most cases.
- **Race Detector:** Run `go run -race` to catch data races.
- **Channels:** Make communication safe and explicit.
- **Select Statement:** Wait on multiple channels at once—like a switchboard for goroutines.

**Example:**

```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    fmt.Println("Received from ch2:", msg)
default:
    fmt.Println("No messages yet!")
}
```

[Ejercicio: Select Statement](../../exercises/part2/08-select-statement/main.go)

---

## 🧪 Go in Action: Goroutine Stress Test

Can you really launch 100,000 goroutines in Go? Yes! Here’s an exercise to try it and understand why it works:

```go
package main
import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    count := 100000
    wg.Add(count)

    for i := 0; i < count; i++ {
        go func(n int) {
            defer wg.Done()
            if n == 0 {
                fmt.Println("First goroutine running!")
            }
        }(i)
    }

    wg.Wait()
    fmt.Println("All goroutines finished!")
}
```

[Exercise: Goroutine Stress Test](../../exercises/part2/08-goroutine-stress/main.go)

**Why can Go handle so many goroutines?**
- **Lightweight:** Each goroutine uses only a few KB of memory (not megabytes like OS threads).
- **Dynamic stack:** Each goroutine’s stack grows and shrinks as needed, starting very small.
- **Own scheduler:** Go has a runtime scheduler that multiplexes goroutines over a few OS threads (M:N model).
- **Non-blocking:** If a goroutine waits (e.g., for I/O), the scheduler moves others to free threads.
- **Efficiency:** This allows Go to handle tens or hundreds of thousands of concurrent tasks without exhausting RAM or CPU.

**Visualization (Mermaid):**

```mermaid
flowchart TD
    subgraph Goroutines
        G1["Goroutine 1"]
        G2["Goroutine 2"]
        G3["..."]
        G100k["Goroutine 100,000"]
    end
    Goroutines --> Scheduler["Go Scheduler"]
    Scheduler --> Threads["Few OS Threads\n(2, 4, 8, ...)"]
    Threads --> CPU["CPU Cores"]
```

**Fun fact:**
- Launching 100,000 OS threads in C/C++ would crash your machine. In Go, it’s just another day for the gopher!

---

## 🎨 Visual Summary

```mermaid
flowchart TD
    G1["Goroutine 1"] --> CH["Channel"]
    G2["Goroutine 2"] --> CH
    G3["Goroutine 3"] --> CH
    CH --> Main["Main Goroutine"]
```

**Explanation:**
- Each goroutine can handle a different client or task.
- Channels act as safe, synchronized pipelines for data between goroutines.
- The main routine can coordinate, collect, or broadcast messages.

---

## 🖼️ TCP Server Concurrency Flow

```mermaid
flowchart TD
    Main["Main Goroutine\nListens on :8082"]
    Main -->|"Accepts connection"| C1["Goroutine for Client 1"]
    Main -->|"..."| Cn["Goroutine for Client n"]
    C1 -->|"Handles Client 1"| End1["..."]
    Cn -->|"Handles Client n"| Endn["..."]
```

**Explanation:**
- The main goroutine listens for new connections.
- For each new client, it spawns a new goroutine to handle communication.
- All clients are served concurrently, so no one waits in line!

---

## 🖼️ Channel Broadcast Flow

```mermaid
flowchart TD
    Broadcaster["Broadcast Goroutine\nSends messages"] --> CH["Channel"]
    CH --> R1["Receiver 1"]
    CH --> R2["Receiver 2"]
    CH --> R3["Receiver 3"]
```

**Explanation:**
- The broadcaster goroutine sends messages into a channel.
- Multiple receivers can read from the channel, each in their own goroutine.
- This pattern is great for chat servers, notifications, or event systems.

---

## 🤩 Fun Facts & Go Memes
- Go’s mascot, the gopher, is often shown juggling or running—just like goroutines!
- The phrase "Don’t communicate by sharing memory; share memory by communicating" is a Go proverb.
- You can run 100,000+ goroutines on a modern laptop—try it!
- Channels are so central to Go, they have their own operator: `<-`.
- Go’s concurrency model is so admired, other languages have copied it (Rust, Elixir, etc.).

---

[Previous: Error Handling and Debugging](07-error-handling-and-debugging.md) | [Next: Context and Cancellation](09-context-and-cancellation.md)
