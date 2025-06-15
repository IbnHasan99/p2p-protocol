package internal

import (
    "encoding/json"
    "io"
    "log"
    "time"
    "github.com/libp2p/go-libp2p/core/network"
)

func (n *Node) HandleTask(s network.Stream) {
    defer s.Close()
    data, err := io.ReadAll(s)
    if err != nil { return }
    var msg TaskMessage
    err = json.Unmarshal(data, &msg)
    if err != nil { return }
    log.Printf("Received task: %s", msg.Instruction)
    time.Sleep(1 * time.Second)
    log.Printf("Completed task: %s", msg.Instruction)
}