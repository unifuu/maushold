-- Seed data for monsters
INSERT INTO monsters (name, type1, type2, base_hp, base_attack, base_defense, base_speed, description, image_url, created_at)
VALUES 
('Pikachu', 'Electric', '', 35, 55, 40, 90, 'Mouse Pokemon', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png', NOW()),
('Charmander', 'Fire', '', 39, 52, 43, 65, 'Lizard Pokemon', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/4.png', NOW()),
('Squirtle', 'Water', '', 44, 48, 65, 43, 'Tiny Turtle Pokemon', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/7.png', NOW()),
('Bulbasaur', 'Grass', 'Poison', 45, 49, 49, 45, 'Seed Pokemon', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/1.png', NOW()),
('Mewtwo', 'Psychic', '', 106, 110, 90, 130, 'Genetic Pokemon', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/150.png', NOW()),
('Dragonite', 'Dragon', 'Flying', 91, 134, 95, 80, 'Dragon Pokemon', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/149.png', NOW()),
('Gengar', 'Ghost', 'Poison', 60, 65, 60, 110, 'Shadow Pokemon', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/94.png', NOW()),
('Lucario', 'Fighting', 'Steel', 70, 110, 70, 90, 'Aura Pokemon', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/448.png', NOW()),
('Eevee', 'Normal', '', 55, 55, 50, 55, 'Evolution Pokemon', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/133.png', NOW()),
('Snorlax', 'Normal', '', 160, 110, 65, 30, 'Sleeping Pokemon', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/143.png', NOW());
