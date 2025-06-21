# TCP in Depth: Protocol Theory and Go Implementation 🔗

> "Imagine sending a letter, but you want to make sure it arrives, in order, and without any missing pages. TCP is the postal service that guarantees delivery, order, and reliability—Go gives you the tools to become a master mail carrier!"

---

## 📚 What is TCP?

**TCP (Transmission Control Protocol)** is the backbone of reliable communication on the internet. It’s like a polite, organized courier: every message is delivered, in order, and checked for errors. If something goes wrong, TCP tries again—no lost mail!

- **Connection-oriented:** Like a phone call—both sides say "hello" before talking.
- **Reliable:** Every packet is acknowledged. If lost, it’s resent.
- **Ordered:** Data arrives in the same order it was sent.
- **Stream-based:** Data flows like a river, not in fixed chunks.

**Analogy:**
- TCP is a certified mail service: you get a receipt, tracking, and confirmation of delivery.

---

## 🧬 How TCP Works (Theory)

1. **Three-Way Handshake:**
   - SYN: "Can we talk?" (Client asks to start a conversation)
   - SYN-ACK: "Yes, let’s talk!" (Server agrees)
   - ACK: "Great, I’m ready!" (Client confirms)
2. **Data Transfer:**
   - Data is sent in segments. Each is acknowledged (ACK).
   - Lost segments are retransmitted automatically.
3. **Connection Teardown:**
   - Both sides say goodbye (FIN/ACK exchange) to close the connection cleanly.

**Diagram:**

```
[Client] --SYN--> [Server]
[Client] <--SYN-ACK-- [Server]
[Client] --ACK--> [Server]
(Data flows)
[Client] --FIN--> [Server]
[Client] <--ACK-- [Server]
```

---

## 🛠️ Go in Action: Simple TCP Client (Paso a Paso)

Este ejemplo muestra cómo Go crea un cliente TCP, se conecta a un servidor, envía una petición HTTP y recibe la respuesta.

```go
package main
import (
    "fmt"      // Para imprimir en consola
    "net"      // Provee funciones de red
    "os"       // Para salir del programa en caso de error
)

func main() {
    // 1. net.Dial abre una conexión TCP a example.com en el puerto 80 (HTTP)
    conn, err := net.Dial("tcp", "example.com:80")
    if err != nil {
        fmt.Println("Error al conectar:", err)
        os.Exit(1)
    }
    // 2. Enviamos una petición HTTP simple
    fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
    // 3. Creamos un buffer para leer la respuesta
    buf := make([]byte, 4096)
    n, _ := conn.Read(buf) // Leemos la respuesta del servidor
    fmt.Println(string(buf[:n])) // Imprimimos la respuesta
    conn.Close() // Cerramos la conexión
}
```

- **¿Qué hace cada parte?**
  - `net.Dial`: Abre el canal TCP y realiza el handshake.
  - `fmt.Fprintf`: Envía datos al servidor.
  - `conn.Read`: Recibe la respuesta.
  - `conn.Close`: Cierra la conexión limpiamente.

[Ejercicio: Simple TCP Client](../../exercises/part2/05-tcp-client/main.go)

---

## 🛠️ Go in Action: Simple TCP Server (Explicado)

Este servidor TCP escucha conexiones en el puerto 8080 y responde con un mensaje a cada cliente.

```go
package main
import (
    "fmt"      // Para imprimir mensajes
    "net"      // Para funciones de red
)

func main() {
    // 1. net.Listen crea un listener TCP en el puerto 8080
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Servidor escuchando en :8080")
    for {
        // 2. Espera conexiones entrantes
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error al aceptar conexión:", err)
            continue
        }
        // 3. Atiende cada conexión en una goroutine (concurrencia)
        go func(c net.Conn) {
            fmt.Fprintln(c, "¡Hola desde el servidor TCP en Go!") // Envía mensaje
            c.Close() // Cierra la conexión
        }(conn)
    }
}
```

- **¿Qué hace cada parte?**
  - `net.Listen`: Abre el puerto y espera conexiones.
  - `ln.Accept`: Acepta una nueva conexión entrante.
  - `go func`: Atiende cada cliente en paralelo (¡puedes tener miles!).
  - `fmt.Fprintln`: Envía datos al cliente.
  - `c.Close`: Cierra la conexión.

[Ejercicio: Simple TCP Server](../../exercises/part2/05-tcp-server/main.go)

---

## 📝 Real-World Example: Echo Server (Comentado)

Un echo server devuelve exactamente lo que recibe. Es ideal para probar conexiones y entender cómo fluye la información.

```go
package main
import (
    "fmt"
    "net"
)

func main() {
    ln, _ := net.Listen("tcp", ":9000")
    fmt.Println("Echo server en :9000")
    for {
        conn, _ := ln.Accept()
        go func(c net.Conn) {
            buf := make([]byte, 1024)
            for {
                n, err := c.Read(buf) // Lee datos del cliente
                if err != nil {
                    break // Si hay error (cliente cerró), salimos
                }
                c.Write(buf[:n]) // Reenviamos lo recibido
            }
            c.Close() // Cerramos la conexión
        }(conn)
    }
}
```

- **¿Qué hace cada parte?**
  - `c.Read`: Lee datos del cliente.
  - `c.Write`: Devuelve los mismos datos (eco).
  - El ciclo permite múltiples mensajes por conexión.

[Ejercicio: TCP Echo Server](../../exercises/part2/05-tcp-echo-server/main.go)

---

## 🧠 ¿Qué hace Go detrás de cámaras?

- `net.Dial` y `net.Listen` usan llamadas al sistema operativo para crear sockets TCP.
- Go gestiona el handshake, la retransmisión y el cierre de conexiones automáticamente.
- Las goroutines permiten manejar miles de conexiones concurrentes sin esfuerzo.
- El paquete `net` es multiplataforma: ¡tu código funciona igual en Windows, Linux y Mac!

---

## 🎨 Visual Summary

```
[Client] <---TCP---> [Server]
   |                   |
[Send] <---ACK---> [Receive]
   |                   |
[Data] <---Order---> [Data]
```

---

## 🤩 Fun Facts & Go Memes
- TCP fue inventado en los 70s y aún es el rey de la red.
- El stack TCP de Go es tan eficiente que lo usan Docker, Kubernetes y muchos servicios en la nube.
- Puedes construir un chat, una herramienta de transferencia de archivos o tu propio protocolo con el paquete `net` de Go.
- TCP es como un amigo confiable: nunca olvida, siempre entrega.
- Las goroutines de Go hacen que manejar miles de conexiones TCP sea pan comido.

---

[Previous: Working with IP, Ports, and Addresses](04-working-with-ip-ports-and-addresses.md) | [Next: UDP in Depth: Protocol Theory and Go Implementation](06-udp-in-depth-protocol-theory-and-go-implementation.md)
