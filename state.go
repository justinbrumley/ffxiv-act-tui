package main

import (
	"encoding/json"
	"fmt"
)

type Player struct {
	ID   uint   `json:"charID"`
	Name string `json:"charName"`
}

type Zone struct {
	ID   uint   `json:"zoneID"`
	Name string `json:"zoneName"`
}

// Combatant can hold a lot of stats, so will try to keep it simple to what's most important for TUI.
// Everything comes in as a string, so any math will need to cast the numbers first.
type Combatant struct {
	Name        string `json:"name"`
	TotalDamage string `json:"damage"`
	DPS         string `json:"dps"`
	Job         string `json:"Job"`
	Deaths      string `json:"deaths"`
	CritPercent string `json:"crithit%"`

	/* Example Combatant Data
	"absorbHeal": "0",
	"critheal%": "0%",
	"critheals": "0",
	"crithit%": "17%",
	"crithits": "5",
	"crittypes": "0.0% Legendary - 0.0% Fabled - 0.0% Mythical",
	"cures": "0",
	"damage": "751665",
	"damage%": "100%",
	"damage-*": "751.00K",
	"damage-b": "0.00",
	"damage-m": "0.75",
	"damageShield": "0",
	"damagetaken": "0",
	"damagetaken-*": "0",
	"deaths": "0",
	"dps": "25341.01",
	"dps-*": "25.00K",
	"duration": "00:29",
	"encdps": "25341.01",
	"encdps-*": "25.00K",
	"enchps": "0.00",
	"enchps-*": "0",
	"healed": "0",
	"healed%": "--",
	"heals": "0",
	"healstaken": "0",
	"healstaken-*": "0",
	"hitfailed": "0",
	"hits": "29",
	"kills": "0",
	"maxheal": "",
	"maxheal-*": "",
	"maxhealward": "",
	"maxhealward-*": "",
	"maxhit": "Perfectio-68630",
	"maxhit-*": "Perfectio-68.00K",
	"misses": "0",
	"n": "\n",
	"name": "YOU",
	"overHeal": "0",
	"powerdrain": "0",
	"powerdrain-*": "0",
	"powerheal": "0",
	"powerheal-*": "0",
	"swings": "29",
	"t": "\t",
	"threatdelta": "0",
	"threatstr": "+(0)0/-(0)0",
	"tohit": "100.00"
	*/
}

// Encounter contains party-wide stats about the current fight as a whole.
type Encounter struct {
	CurrentZoneName string `json:"CurrentZoneName"`
	Damage          string `json:"damage"`
	Deaths          string `json:"deaths"`
	DPS             string `json:"dps"`
	Duration        string `json:"duration"`
	Title           string `json:"title"`
	MaxHit          string `json:"maxhit"`

	/* Example Encounter Data
	   "CurrentZoneName": "The Lavender Beds",
	   "DAMAGE-*": "751K",
	   "DAMAGE-b": "0",
	   "DAMAGE-k": "752",
	   "DAMAGE-m": "1",
	   "DPS": "25341",
	   "DPS-*": "25K",
	   "DPS-k": "25",
	   "DPS-m": "DPS-m",
	   "DURATION": "30",
	   "ENCDPS": "25341",
	   "ENCDPS-*": "25K",
	   "ENCDPS-k": "25",
	   "ENCDPS-m": "0",
	   "ENCHPS": "0",
	   "ENCHPS-*": "0",
	   "ENCHPS-k": "0",
	   "ENCHPS-m": "0",
	   "Last10DPS": "0",
	   "Last30DPS": "0",
	   "Last60DPS": "0",
	   "MAXHEAL": "",
	   "MAXHEAL-*": "",
	   "MAXHEALWARD": "",
	   "MAXHEALWARD-*": "",
	   "MAXHIT": "YOU-68630",
	   "MAXHIT-*": "YOU-68.00K",
	   "TOHIT": "100",
	   "critheal%": "NaN",
	   "critheals": "5",
	   "crithit%": "0%",
	   "crithits": "5",
	   "cures": "0",
	   "damage": "751665",
	   "damage-*": "751.00K",
	   "damage-m": "0.75",
	   "damagetaken": "0",
	   "damagetaken-*": "0",
	   "deaths": "0",
	   "dps": "25341.01",
	   "dps-*": "dps-*",
	   "duration": "00:29",
	   "encdps": "25341.01",
	   "encdps-*": "25.00K",
	   "enchps": "0.00",
	   "enchps-*": "0",
	   "healed": "0",
	   "heals": "0",
	   "healstaken": "0",
	   "healstaken-*": "0",
	   "hitfailed": "0",
	   "hits": "29",
	   "kills": "0",
	   "maxheal": "",
	   "maxheal-*": "",
	   "maxhealward": "",
	   "maxhealward-*": "",
	   "maxhit": "YOU-Perfectio-68630",
	   "maxhit-*": "YOU-Perfectio-68K",
	   "misses": "0",
	   "n": "\n",
	   "powerdrain": "0",
	   "powerdrain-*": "0",
	   "powerheal": "0",
	   "powerheal-*": "0",
	   "swings": "29",
	   "t": "\t",
	   "title": "Striking Dummy",
	   "tohit": "100.00"
	*/
}

type CombatData struct {
	IsActive interface{} `json:"isActive"`

	// Encounter includes party-wide stats for the current fight.
	Encounter Encounter `json:"Encounter"`

	// Combatant is a map of players and their individual combat stats.
	Combatants map[string]Combatant `json:"Combatant"`
}

type State struct {
	PrimaryPlayer *Player
	Players       map[uint]Player
	CurrentZone   *Zone
}

var state *State

const (
	MessageType_ChangeZone   = "ChangeZone"
	MessageType_SendCharName = "SendCharName"
	MessageType_CombatData   = "CombatData"
)

func init() {
	state = &State{
		Players: make(map[uint]Player),
	}
}

// parseZone parses zone data from message.
// Example message: map[type:ChangeZone zoneID:340 zoneName:The Lavender Beds]
func parseZone(msg map[string]interface{}) (*Zone, error) {
	str, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	zone := &Zone{}
	if err := json.Unmarshal(str, &zone); err != nil {
		return nil, err
	}

	return zone, nil
}

// parsePlayer parses the player name and ID from SendCharName message.
func parsePlayer(msg map[string]interface{}) (*Player, error) {
	str, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	player := &Player{}
	if err := json.Unmarshal(str, &player); err != nil {
		return nil, err
	}

	return player, nil
}

func parseCombatData(msg map[string]interface{}) (*CombatData, error) {
	str, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	data := &CombatData{}
	if err := json.Unmarshal(str, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func handleMessage(payload *Payload) error {
	switch payload.MessageType {
	// Parse Zone and set current zone in state
	case MessageType_ChangeZone:
		zone, err := parseZone(payload.Message)
		if err != nil {
			return err
		}

		state.CurrentZone = zone
		fmt.Printf("Parsed Zone: %v\n", zone)

	// Parse player information
	case MessageType_SendCharName:
		player, err := parsePlayer(payload.Message)
		if err != nil {
			return err
		}

		state.Players[player.ID] = *player
		if t, ok := payload.Message["type"]; ok && t.(string) == "ChangePrimaryPlayer" {
			state.PrimaryPlayer = player
		}

		fmt.Printf("Parsed Player: %v\n", player)

	case MessageType_CombatData:
		data, err := parseCombatData(payload.Message)
		if err != nil {
			fmt.Printf("Error parsing combat data: %v\n", err)
			return err
		}

		fmt.Printf("Parsed combat data: %v\n", data)

	default:
		fmt.Printf("Unrecognized message received: %v\n", payload.MessageType)
		fmt.Printf("[")
		for key, _ := range payload.Message {
			fmt.Printf(" %v ", key)
		}
		fmt.Println("]")
	}

	return nil
}
