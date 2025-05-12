package main

import(
  "os/signal"
  "os"
  "log"
)

func main()  {
  done := make(chan os.Signal, 1)
  signal.Notify(done, os.Interrupt)

  wsManager := NewWebSocketManager()
  defer wsManager.Close()

  go wsManager.fetchNewTokens()
  go wsManager.analyzeTokenTrades()

  <-done
  log.Println("Closing WebSocket connections")
}
