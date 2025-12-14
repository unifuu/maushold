package service

import (
	"log"

	"maushold/pokemon-service/model"
	"maushold/pokemon-service/repository"
)

func SeedPokemon(repo repository.PokemonRepository) {
	// Check if already seeded
	all, err := repo.FindAll()
	if err == nil && len(all) > 0 {
		log.Println("Pokemon data already seeded")
		return
	}

	starterPokemon := []model.Pokemon{
		{ID: 1, Name: "Bulbasaur", Type1: "Grass", Type2: "Poison", BaseHP: 45, BaseAttack: 49, BaseDefense: 49, BaseSpeed: 45, Description: "A strange seed was planted on its back at birth."},
		{ID: 4, Name: "Charmander", Type1: "Fire", Type2: "", BaseHP: 39, BaseAttack: 52, BaseDefense: 43, BaseSpeed: 65, Description: "Obviously prefers hot places."},
		{ID: 7, Name: "Squirtle", Type1: "Water", Type2: "", BaseHP: 44, BaseAttack: 48, BaseDefense: 65, BaseSpeed: 43, Description: "After birth, its back swells and hardens into a shell."},
		{ID: 25, Name: "Pikachu", Type1: "Electric", Type2: "", BaseHP: 35, BaseAttack: 55, BaseDefense: 40, BaseSpeed: 90, Description: "When several of these Pokémon gather, their electricity could build."},
		{ID: 39, Name: "Jigglypuff", Type1: "Normal", Type2: "Fairy", BaseHP: 115, BaseAttack: 45, BaseDefense: 20, BaseSpeed: 20, Description: "When its huge eyes light up, it sings a mysteriously soothing melody."},
		{ID: 133, Name: "Eevee", Type1: "Normal", Type2: "", BaseHP: 55, BaseAttack: 55, BaseDefense: 50, BaseSpeed: 55, Description: "Its genetic code is irregular."},
		{ID: 143, Name: "Snorlax", Type1: "Normal", Type2: "", BaseHP: 160, BaseAttack: 110, BaseDefense: 65, BaseSpeed: 30, Description: "Very lazy. Just eats and sleeps."},
		{ID: 150, Name: "Mewtwo", Type1: "Psychic", Type2: "", BaseHP: 106, BaseAttack: 110, BaseDefense: 90, BaseSpeed: 130, Description: "It was created by a scientist after years of horrific gene splicing."},
		{ID: 94, Name: "Gengar", Type1: "Ghost", Type2: "Poison", BaseHP: 60, BaseAttack: 65, BaseDefense: 60, BaseSpeed: 110, Description: "Under a full moon, this Pokémon likes to mimic the shadows of people."},
		{ID: 6, Name: "Charizard", Type1: "Fire", Type2: "Flying", BaseHP: 78, BaseAttack: 84, BaseDefense: 78, BaseSpeed: 100, Description: "Spits fire that is hot enough to melt boulders."},
	}

	for _, p := range starterPokemon {
		repo.Create(&p)
	}

	log.Println("Seeded initial Pokemon data")
}
