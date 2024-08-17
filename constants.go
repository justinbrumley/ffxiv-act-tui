package main

var TankJobs = []string{"Pld", "War", "Drk", "Gnb"}
var HealerJobs = []string{"Whm", "Sch", "Ast", "Sge"}
var DpsJobs = []string{"Mnk", "Drg", "Nin", "Sam", "Rpr", "Vpr", "Brd", "Mch", "Dnc", "Blm", "Smn", "Rdm", "Pct", "Blu"}

const TankColor = "81"   // Blueish
const HealerColor = "41" // Greenish
const DpsColor = "203"   // Reddish

func GetRoleColor(job string) string {
	for _, j := range HealerJobs {
		if j == job {
			return HealerColor
		}
	}

	for _, j := range TankJobs {
		if j == job {
			return TankColor
		}
	}

	return DpsColor
}
