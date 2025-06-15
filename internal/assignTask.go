package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
)

// one node run as dispatcher
func (n *Node) RunDispatcher() {
	log.Println("Running as dispatcher")
	file, err := os.Open("instructions.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		for _, instr := range tokens {
			taskType := strings.ToUpper(string(instr[0]))
			if taskType == n.Role {
				log.Printf("Executing locally: %s", instr)
				time.Sleep(1 * time.Second)
			} else {
				n.SendTask(taskType, instr)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

}

// send instruction to other nodes
func (n *Node) SendTask(role, instr string) {
	n.Lock.Lock()
	peerID, ok := n.Peers[role]
	n.Lock.Unlock()
	if !ok {
		log.Printf("No peer found for role %s", role)
		return
	}
	s, err := n.Host.NewStream(context.Background(), peerID, TaskProtocolID)
	if err != nil {
		log.Printf("Error opening stream to %s: %v", peerID, err)
		return
	}
	defer s.Close()
	msg := TaskMessage{Instruction: instr}
	data, _ := json.Marshal(msg)
	s.Write(data) // send instruction over stream
	log.Printf("Sent %s to role %s", instr, role)
}
