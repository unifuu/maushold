package service

import (
	"fmt"
	"math/rand"

	"maushold/battle-service/model"
)

type BattleEngine struct{}

func NewBattleEngine() *BattleEngine {
	return &BattleEngine{}
}

func (e *BattleEngine) SimulateBattle(p1, p2 *model.PlayerMonster) (int, string) {
	log := ""
	hp1 := p1.HP
	hp2 := p2.HP

	log += fmt.Sprintf("⚔️ Battle Start!\n%s (HP: %d) vs %s (HP: %d)\n\n",
		p1.Nickname, hp1, p2.Nickname, hp2)

	round := 1
	for hp1 > 0 && hp2 > 0 && round <= 20 {
		log += fmt.Sprintf("=== Round %d ===\n", round)

		if p1.Speed >= p2.Speed {
			damage := calculateDamage(p1.Attack, p2.Defense)
			hp2 -= damage
			log += fmt.Sprintf("%s attacks for %d damage! %s HP: %d\n",
				p1.Nickname, damage, p2.Nickname, maxInt(hp2, 0))

			if hp2 > 0 {
				damage = calculateDamage(p2.Attack, p1.Defense)
				hp1 -= damage
				log += fmt.Sprintf("%s attacks for %d damage! %s HP: %d\n",
					p2.Nickname, damage, p1.Nickname, maxInt(hp1, 0))
			}
		} else {
			damage := calculateDamage(p2.Attack, p1.Defense)
			hp1 -= damage
			log += fmt.Sprintf("%s attacks for %d damage! %s HP: %d\n",
				p2.Nickname, damage, p1.Nickname, maxInt(hp1, 0))

			if hp1 > 0 {
				damage = calculateDamage(p1.Attack, p2.Defense)
				hp2 -= damage
				log += fmt.Sprintf("%s attacks for %d damage! %s HP: %d\n",
					p1.Nickname, damage, p2.Nickname, maxInt(hp2, 0))
			}
		}

		log += "\n"
		round++
	}

	if hp1 > hp2 {
		log += fmt.Sprintf("🏆 %s wins!\n", p1.Nickname)
		return 1, log
	}
	log += fmt.Sprintf("🏆 %s wins!\n", p2.Nickname)
	return 2, log
}

func calculateDamage(attack, defense int) int {
	baseDamage := attack - (defense / 2)
	if baseDamage < 1 {
		baseDamage = 1
	}

	variance := rand.Intn(10) - 5
	damage := baseDamage + variance

	if damage < 1 {
		damage = 1
	}

	return damage
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
