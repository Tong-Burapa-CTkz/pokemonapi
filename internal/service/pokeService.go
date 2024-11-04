package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	DTOs "sms2pro/internal/DTOs"
)

type Pokemon struct {
	Name      string           `json:"name"`
	Base_Exp  int              `json:"base_experience"`
	Weight    int              `json:"weight`
	Height    int              `json:"height`
	Abilities []PokemonAbility `json:"abilities"`
}

type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonAbility struct {
	Ability  Ability `json:"ability"`
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
}

var client = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_ADDRESS") + ":" + os.Getenv("REDIS_PORT"),
	Password: "",
	DB:       0,
})

func GetPokemon(c *gin.Context) {
	name := c.Param("name")
	ctx := context.Background()

	// Call redis
	pokemonJSON, err := client.Get(ctx, name).Result()
	fmt.Println(err)

	if err == nil {
		var redisPokemon Pokemon
		err := json.Unmarshal([]byte(pokemonJSON), &redisPokemon)
		response := DTOs.Pokemon{
			Name:     redisPokemon.Name,
			Base_Exp: redisPokemon.Base_Exp,
			Weight:   redisPokemon.Weight,
			Height:   redisPokemon.Height,
		}
		if err == nil {
			fmt.Println("From cache")
			c.JSON(http.StatusOK, response)
			return
		}
	}

	// No redis cache
	apiURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from the Pokémon API"})
		return
	}
	defer resp.Body.Close()

	// Read and parse the API response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Pokémon data"})
		return
	}

	// Cache Pokemon info
	jsonData, err := json.Marshal(pokemon)
	fmt.Println(jsonData)
	if err == nil {
		client.Set(ctx, name, string(jsonData), 600*time.Second)
	}

	response := DTOs.Pokemon{
		Name:     pokemon.Name,
		Base_Exp: pokemon.Base_Exp,
		Weight:   pokemon.Weight,
		Height:   pokemon.Height,
	}

	c.JSON(http.StatusOK, response)
}

func GetPokemonAbility(c *gin.Context) {
	name := c.Param("name")
	ctx := context.Background()

	pokemonJSON, err := client.Get(ctx, name).Result()
	fmt.Println(pokemonJSON)

	if err == nil {
		var cachedPokemonAbility Pokemon
		err := json.Unmarshal([]byte(pokemonJSON), &cachedPokemonAbility.Abilities)
		if err == nil {
			c.JSON(http.StatusOK, cachedPokemonAbility.Abilities)
			return
		}
	}

	// No redis cache
	apiURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	resp, err := http.Get(apiURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from the Pokémon API"})
		return
	}
	defer resp.Body.Close()

	// Read and parse the API response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Pokémon data"})
		return
	}

	jsonData, err := json.Marshal(pokemon)
	if err == nil {
		client.Set(ctx, name, string(jsonData), 600*time.Second)
	}

	c.JSON(http.StatusOK, pokemon.Abilities)

}
