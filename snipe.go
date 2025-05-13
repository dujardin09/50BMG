package main

import(
	"encoding/json"
  "log"
  "sync"
)

type AddressBook struct {
	mu       sync.RWMutex
	addresses []string
}

func NewAddressBook() *AddressBook {
	return &AddressBook{
		addresses: make([]string, 0),
	}
}

func (ab *AddressBook) Add(address string) {
	ab.mu.Lock()
	defer ab.mu.Unlock()
	ab.addresses = append(ab.addresses, address)
}

func (ab *AddressBook) GetAll() []string {
	ab.mu.RLock()
	defer ab.mu.RUnlock()
	return append([]string(nil), ab.addresses...) // copy
}

func (wsm *WebSocketManager) fetchNewTokens(update chan struct {}, ab *AddressBook) {
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
      ab.Add(mint)
      select {
      case update <- struct{}{}:
      default:
      }
    }
  }
}

func (wsm *WebSocketManager) analyzeTokenTrades(update chan struct {}, ab *AddressBook) {
  for {
    select {
    case <- update:
      log.Printf("Nouvelle address :%v", ab.GetAll())
    default:
    }
  }
  // if err := wsm.subscribe("subscribeTokenTrade"); err != nil {
  //   log.Printf("Trade subscription failed: %v", err)
  // }
  //
  // message, err := wsm.readMessage()
  // if err != nil {
  //   log.Printf("Read error: %v", err)
  // } else {
  //     log.Printf("%s\n", message)
  // }
}
