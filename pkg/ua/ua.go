package ua

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type AgentType string

const (
	Chrome  AgentType = "chrome"
	Firefox AgentType = "firefox"
	Safari  AgentType = "safari"
	IE      AgentType = "ie"

	Android  AgentType = "android"
	Ios      AgentType = "ios"
	Linux    AgentType = "linux"
	MacOS    AgentType = "macos"
	Mobile   AgentType = "mobile"
	Ipad     AgentType = "ipad"
	Iphone   AgentType = "iphone"
	Computer AgentType = "computer"
)

var (
	uaMap map[string][]string
)

// go-bindata -o=pkg/ua/asset.go -pkg=ua pkg/ua/ua.json
func GetAllUaMap() map[string][]string {
	if len(uaMap) > 0 {
		return uaMap
	} else {
		fmt.Println("pares json file")
		uaMap = make(map[string][]string, 12)
		jsonBytes, err := Asset("pkg/ua/ua.json")
		if err != nil {
			log.Fatalf("Asset() Read Fail,%s", err.Error())
		}
		err = json.Unmarshal(jsonBytes, &uaMap)
		if err != nil {
			log.Fatalf("json.Unmarshal Error,%s", err.Error())
		}
		return uaMap
	}
	return nil
}

type UserAgent struct {
	Agents []AgentType
	mux    sync.Mutex
}

func New() *UserAgent {
	return &UserAgent{
		Agents: make([]AgentType, 0),
	}
}

func (ua *UserAgent) Use(agents ...AgentType) *UserAgent {
	if len(agents) > 0 {
		ua.mux.Lock()
		defer ua.mux.Unlock()
		for _, a := range agents {
			if ua.IssetAgent(a) == false {
				ua.Agents = append(ua.Agents, a)
			}
		}
	}
	return ua
}

func (ua *UserAgent) Random() (string, error) {
	agentStr := string(ua.RandomAgent())
	um := GetAllUaMap()
	uaSlice, ok := um[agentStr]
	if !ok {
		return "", errors.New("Unable to get UA")
	}
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(uaSlice))
	return uaSlice[r], nil
}

func (ua *UserAgent) RandomAgent() AgentType {
	var prepareAgent []AgentType
	count := len(ua.Agents)
	if count > 0 {
		if len(ua.Agents) == 1 {
			return ua.Agents[0]
		} else {
			prepareAgent = ua.Agents
		}
	} else {
		prepareAgent = []AgentType{
			Chrome, Firefox, Safari, IE, Android, Ios, Linux, MacOS, Mobile, Ipad, Iphone, Computer,
		}
		count = len(prepareAgent)
	}
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(count)
	return prepareAgent[r]
}

func (ua *UserAgent) IssetAgent(agent AgentType) bool {
	isset := false
	for _, a := range ua.Agents {
		if agent == a {
			isset = true
			break
		}
	}
	return isset
}
