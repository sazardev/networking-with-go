// 09-context-values/main.go
// Demonstrates context values and propagation.
package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.WithValue(context.Background(), "userID", 42)
	process(ctx)
}

func process(ctx context.Context) {
	if v := ctx.Value("userID"); v != nil {
		fmt.Println("UserID from context:", v)
	} else {
		fmt.Println("No userID in context")
	}
}
