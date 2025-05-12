package main

import(
	"encoding/json"
  "log"
)

func (wsm *WebSocketManager) fetchNewTokens() {
  if err := wsm.subscribe("subscribeNewToken"); err != nil {
    log.Printf("subscribeNewToken failed: %v", err)
    return
  }

  message, err := wsm.readMessage()
  if err != nil {
    log.Printf("Read error: %v", err)
  } else {
      log.Printf("%s\n", message)
  }

  for {
    message, err := wsm.readMessage()
    if err != nil {
      log.Printf("Read error: %v", err)
      continue
    }

		var token map[string]interface{}
		if err := json.Unmarshal(message, &token); err != nil {
			log.Printf("JSON parse error: %v", err)
			continue
		}

    solAmount, _ := token["solAmount"].(float64)
    vSolInBondingCurve, _ := token["vSolInBondingCurve"].(float64)
    marketCapSol, _ := token["marketCapSol"].(float64)
    mint, _ := token["mint"].(string)

    if solAmount >= 3 && vSolInBondingCurve >= 15 && marketCapSol < 40 {
      log.Printf("%s\n", mint)
    }
  }
}

func (wsm *WebSocketManager) analyzeTokenTrades() {
  if err := wsm.subscribe("subscribeTokenTrade"); err != nil {
    log.Printf("Trade subscription failed: %v", err)
  }

  message, err := wsm.readMessage()
  if err != nil {
    log.Printf("Read error: %v", err)
  } else {
      log.Printf("%s\n", message)
  }
}
