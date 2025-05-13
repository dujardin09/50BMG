package main

import(
  "os/signal"
  "os"
  "log"
)

func main()  {
  done := make(chan os.Signal, 1)
  signal.Notify(done, os.Interrupt)

  wsManagerFetch := NewWebSocketManager()
  wsManagerAnalyze := NewWebSocketManager()
  defer wsManagerFetch.Close()
  defer wsManagerAnalyze.Close()

  update := make(chan struct{}, 1)
  addressBook := NewAddressBook()

  go wsManagerFetch.fetchNewTokens(update, addressBook)
  go wsManagerAnalyze.analyzeTokenTrades(update, addressBook)

  <-done
  log.Println("Closing WebSocket connections")
}
