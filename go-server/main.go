package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "go-server/docs"
	"go-server/pkg/authentication"
	"go-server/pkg/games"
	"go-server/pkg/reviews"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

//	@title			Game Review API
//	@version		1.0
//	@description	This is an Api AuthService for Cool Game Review Api.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	autolarry55@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/
// @schemes	http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// get the environment variables and initialize the application
	initResponse, err := InitializationHandler()
	if err != nil {
		panic(err)
	}

	var app *fiber.App
	app = fiber.New()

	// log the requests
	app.Use(cors.New())
	//app.Use(csrf.New())
	app.Use(logger.New())
	app.Use(recover.New())

	apiGroup := app.Group("/api/")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "apiVersion", "/v1")

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/loaderio-baf5658b393ffe75ded7e5209eb81d79.txt", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("resources/loaderTest.txt")
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello Here")
	})
	// register a ping route
	apiGroup.Get("/v1/ping", ping)

	authNeeds := authentication.NewAuthNeeds()

	err = authentication.Register(initResponse.MongoDbClient, ctx, apiGroup)
	err = games.Register(initResponse.MongoDbClient, ctx, apiGroup, authNeeds)
	err = reviews.Register(initResponse.MongoDbClient, ctx, apiGroup, authNeeds)

	//_generateGames(initResponse.MongoDbClient)

	if err != nil {
		log.Fatal(err)
		return
	}

	port := os.Getenv("PORT")

	fmt.Println("Server is running on port: " + port)

	err = app.Listen(":" + port)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

var d []games.Game

type Rand struct {
	Title     string              `json:"title"`
	Developer string              `json:"developer"`
	Sum       string              `json:"sum"`
	Dev       string              `json:"dev"`
	Genres    []map[string]string `json:"genres"`
	Pub       string              `json:"pub"`
	Image     string              `json:"image"`
}

func _generateGames(mongo *mongo.Client) {
	gs := []Rand{
		{
			Title: "Galactic Explorer",
			Sum:   "Embark on a thrilling journey through the stars in this space exploration game. Discover new planets, encounter alien species, and unravel the mysteries of the universe.",
			Dev:   "Galactic Explorer",
			Pub:   "Galactic Explorer",
			Image: "https://preview.redd.it/twz4uw5pvpj91.png?width=1280&format=png&auto=webp&v=enabled&s=30c1352154f8936c65bfc43ec767fc501f2d3deb",
			Genres: []map[string]string{
				{
					"title": "Adventure",
					"slug":  "adventure",
				},
				{
					"title": "Indie",
					"slug":  "indie",
				},
				{
					"title": "Simulation",
					"slug":  "simulation",
				},
			},
		},
		{
			Title: "Frostbite",
			Sum:   "Survive a brutal winter in this intense survival game. Scavenge for resources, build shelter, and fend off dangerous predators as you fight to stay alive in a frozen wasteland.",
			Dev:   "Frozen North Studios",
			Pub:   "Arctic Games",
			Image: "https://www.pngmart.com/files/22/Fornite-Frostbite-PNG.png",
			Genres: []map[string]string{
				{
					"title": "survival",
					"slug":  "survival",
				},
				{
					"title": "open world",
					"slug":  "open-world",
				},
				{
					"title": "Simulation",
					"slug":  "simulation",
				},
			},
		},
		{
			Title: "Neon Nights",
			Sum:   "Race through the city streets in this neon-soaked racing game. Customize your ride, compete against other racers, and become the king of the night.",
			Dev:   "Luminous Games",
			Pub:   "Midnight Racing",
			Image: "https://www.ggrecon.com/media/wlnjq5gh/rocket-league-neon-nights.png",
			Genres: []map[string]string{
				{
					"title": "action",
					"slug":  "action",
				},
				{
					"title": "racing",
					"slug":  "racing",
				},
				{
					"title": "arcade",
					"slug":  "arcade",
				},
			},
		},
		{
			Title: "Terraformers",
			Sum:   "Transform a barren planet into a thriving world in this strategy game. Build infrastructure, research new technologies, and navigate diplomatic relationships with rival factions as you compete to terraform the planet.",
			Dev:   "Galactic Enterprises",
			Pub:   "Interplanetary Games",
			Image: "https://img.gg.deals/dd/14/8beabddf6017017c5eff09cdb3edb601d384_1232xr706_Q100.jpg",
			Genres: []map[string]string{
				{
					"title": "strategy",
					"slug":  "strategy",
				},
				{
					"title": "simulation",
					"slug":  "simulation",
				},
				{
					"title": "sci-fi",
					"slug":  "sci-fi",
				},
			},
		},
		{
			Title: "Galactic Conquest",
			Sum:   "Build an empire and conquer the galaxy in this epic space strategy game. Colonize planets, build fleets, and engage in tactical battles against other factions.",
			Dev:   "Nova Studios",
			Pub:   "Universal Games",
			Image: "https://www.gannett-cdn.com/media/USATODAY/GenericImages/2013/04/29/screen-shot-2013-04-29-at-8_43_40-am-16_9.jpg?width=660&height=373&fit=crop&format=png",
			Genres: []map[string]string{
				{
					"title": "strategy",
					"slug":  "strategy",
				},
				{
					"title": "sci-fi",
					"slug":  "sci-fi",
				},
				{
					"title": "multiplayer",
					"slug":  "multiplayer",
				},
			},
		},
		{
			Title: "Dungeon Delve",
			Sum:   "Explore a dangerous dungeon filled with traps, monsters, and treasure in this classic adventure game. Use your wits and skills to survive and claim the riches that lie within.",
			Dev:   "Catacomb Games",
			Pub:   "Dragonfire Publishing",
			Image: "https://ksr-ugc.imgix.net/assets/028/055/586/bcf83900fe65e157f45fd3fbdcf0d5c7_original.jpg?ixlib=rb-4.0.2&crop=faces&w=1552&h=873&fit=crop&v=1581450028&auto=format&frame=1&q=92&s=a8afee9291d566f1654f2704071c1b16",
			Genres: []map[string]string{
				{
					"title": "adventure",
					"slug":  "adventure",
				},
				{
					"title": "fantasy",
					"slug":  "fantasy",
				},
				{
					"title": "single-player",
					"slug":  "single-player",
				},
			},
		},
		{
			Title: "Rally Racer",
			Sum:   "Speed through winding tracks and sharp turns in this intense racing game. Customize your car, compete against other drivers, and become the champion of the rally circuit.",
			Dev:   "Velocity Games",
			Pub:   "Raceway Entertainment",
			Image: "https://upload.wikimedia.org/wikipedia/commons/6/6a/Petter_Solberg_-_2006_Cyprus_Rally.jpg",
			Genres: []map[string]string{
				{
					"title": "racing",
					"slug":  "racing",
				},
				{
					"title": "sports",
					"slug":  "sports",
				},
				{
					"title": "multiplayer",
					"slug":  "multiplayer",
				},
			},
		},
		{
			Title: "Ninja Assassin",
			Sum:   "Take on the role of a skilled ninja and eliminate your targets with stealth and precision in this action-packed game. Sneak past guards, sabotage defenses, and strike at the right moment to succeed.",
			Dev:   "Shinobi Studios",
			Pub:   "Ninjitsu Games",
			Image: "https://a0.anyrgb.com/pngimg/986/1980/ninja-academy-ninja-blade-hattori-hanzo-ninja-assassin-ninja-ninjutsu-shuriken-katana-samurai-lance.png",
			Genres: []map[string]string{
				{
					"title": "action",
					"slug":  "action",
				},
				{
					"title": "stealth",
					"slug":  "stealth",
				},
				{
					"title": "single-player",
					"slug":  "single-player",
				},
			},
		},
		{
			Title: "Civilization Builder",
			Sum:   "Lead your civilization from the Stone Age to the Information Age in this grand strategy game. Develop your cities, conduct diplomacy, and wage war against rival nations to achieve global domination.",
			Dev:   "Empire Games",
			Pub:   "Worldwide Gaming",
			Image: "https://i.ytimg.com/vi/cQzQR9SKYQo/maxresdefault.jpg",
			Genres: []map[string]string{
				{
					"title": "strategy",
					"slug":  "strategy",
				},
				{
					"title": "simulation",
					"slug":  "simulation",
				},
				{
					"title": "historical",
					"slug":  "historical",
				},
			},
		},
		{
			Title: "Dragon's Den",
			Sum:   "Build a lair for your dragon and protect your treasure hoard from would-be thieves in this management game. Manage your resources, train your dragon, and fend off intruders to become the richest dragon in the land.",
			Dev:   "Mythical Games",
			Pub:   "Dragonfire Publishing",
			Image: "https://armadaboost.com/storage/product/images/buy-sunspire-carries22.png",
			Genres: []map[string]string{
				{
					"title": "management",
					"slug":  "management",
				},
				{
					"title": "fantasy",
					"slug":  "fantasy",
				},
				{
					"title": "single-player",
					"slug":  "single-player",
				},
			},
		},
		{
			Title: "Zombie Apocalypse",
			Sum:   "Survive in a world overrun by the undead in this survival game. Scavenge for supplies, fortify your shelter, and fight off hordes of zombies to stay alive.",
			Dev:   "Undead Labs",
			Pub:   "Zombie Games",
			Image: "https://www.gamersdecide.com/sites/default/files/styles/news_images/public/main_image_15_zom_.jpg",
			Genres: []map[string]string{
				{
					"title": "survival",
					"slug":  "survival",
				},
				{
					"title": "horror",
					"slug":  "horror",
				},
				{
					"title": "multiplayer",
					"slug":  "multiplayer",
				},
			},
		},
		{
			Title: "Mystic Quest",
			Sum:   "Embark on a mystical journey and explore an enchanted world filled with wonder and danger in this RPG. Choose your character class, level up your skills, and uncover the secrets of the ancient realm.",
			Dev:   "Arcane Games",
			Pub:   "Fantasy Books",
			Image: "https://pbs.twimg.com/media/EtFsfMRXAAAD2yX.jpg",
			Genres: []map[string]string{
				{
					"title": "RPG",
					"slug":  "rpg",
				},
				{
					"title": "fantasy",
					"slug":  "fantasy",
				},
				{
					"title": "single-player",
					"slug":  "single-player",
				},
			},
		},
		{
			Title: "Gladiator Arena",
			Sum:   "Step into the arena and battle against fierce opponents in this action-packed fighting game. Choose your weapons and fight your way to the top to become the champion of the gladiators.",
			Dev:   "Arena Games",
			Pub:   "Gladiator Entertainment",
			Image: "https://t4.ftcdn.net/jpg/05/58/21/53/360_F_558215338_THTSPM2aCRYvjI10vL3SKl4zPEFLR4zc.jpg",
			Genres: []map[string]string{
				{
					"title": "fantasy",
					"slug":  "fantasy",
				},
				{
					"title": "single-player",
					"slug":  "single-player",
				},
			},
		},
		{
			Title: "Chrono Crusade",
			Sum:   "Travel through time and prevent a dark force from altering the course of history in this epic adventure game. Collect artifacts, solve puzzles, and make choices that will impact the fate of humanity.",
			Dev:   "Time Warp Studios",
			Pub:   "Chrono Games",
			Image: "https://i.ytimg.com/vi/rXFBasb4xwI/maxresdefault.jpg",
			Genres: []map[string]string{
				{
					"title": "adventure",
					"slug":  "adventure",
				},
				{
					"title": "puzzle",
					"slug":  "puzzle",
				},
				{
					"title": "sci-fi",
					"slug":  "sci-fi",
				},
			},
		},
		{
			Title: "Mage's Quest",
			Sum:   "Embark on a magical journey through a mystical realm in this action-packed RPG. Battle fierce monsters, level up your character, and discover powerful spells as you seek to defeat an ancient evil.",
			Dev:   "Mythical Games",
			Pub:   "Wizard World Entertainment",
			Image: "https://assets2.rockpapershotgun.com/exodiamagedecklistguideungoro.jpg/BROK/thumbnail/1200x900/quality/100/exodiamagedecklistguideungoro.jpg",
			Genres: []map[string]string{
				{
					"title": "RPG",
					"slug":  "rpg",
				},
				{
					"title": "action",
					"slug":  "action",
				},
				{
					"title": "fantasy",
					"slug":  "fantasy",
				},
			},
		},
		{
			Title: "Cosmic Colonies",
			Sum:   "Build and manage a thriving space colony in this simulation game. Explore the galaxy, mine valuable resources, and expand your territory as you strive to become the dominant force in the universe.",
			Dev:   "Star Systems Studio",
			Pub:   "Galactic Games",
			Image: "https://image-pastemagazine-com-public-bucket.storage.googleapis.com/wp-content/uploads/2022/06/20233759/cosmic_colonies.jpg",
			Genres: []map[string]string{
				{
					"title": "simulation",
					"slug":  "simulation",
				},
				{
					"title": "strategy",
					"slug":  "strategy",
				},
				{
					"title": "sci-fi",
					"slug":  "sci-fi",
				},
			},
		},
		{
			Title: "Sword and Sorcery",
			Sum:   "Journey through a mythical land of dragons and magic in this immersive RPG. Customize your hero, engage in epic battles, and unlock powerful spells as you seek to vanquish the evil that threatens the kingdom.",
			Dev:   "Mythical Games",
			Pub:   "Wizard World Entertainment",
			Image: "https://www.grimdarkmagazine.com/wp-content/uploads/2018/02/conanageof_940-492.jpg",
			Genres: []map[string]string{
				{
					"title": "RPG",
					"slug":  "rpg",
				},
				{
					"title": "action",
					"slug":  "action",
				},
				{
					"title": "fantasy",
					"slug":  "fantasy",
				},
			},
		},
		{
			Title: "Terra Nova",
			Sum:   "Colonize a new world and build a thriving civilization in this strategy game. Manage resources, research new technologies, and fend off alien threats as you compete with rival factions for control of the planet.",
			Dev:   "New Horizons Games",
			Pub:   "Galactic Enterprises",
			Image: "https://www.boardgamequest.com/wp-content/uploads/2023/01/Terra-Nova.jpg",
			Genres: []map[string]string{
				{
					"title": "simulation",
					"slug":  "simulation",
				},
				{
					"title": "strategy",
					"slug":  "strategy",
				},
				{
					"title": "sci-fi",
					"slug":  "sci-fi",
				},
			},
		},
		{
			Title: "The Great Heist",
			Sum:   "Plan and execute the ultimate heist in this action-packed game. Assemble a crew, acquire gear, and take on the most challenging targets in the city. But beware, the authorities are always watching, and one false move could land you in jail!",
			Dev:   "Crime Games",
			Pub:   "Rogue Entertainment",
			Image: "https://sm.ign.com/t/ign_in/blogroll/g/gta-5-onli/gta-5-online-heists-hands-on-theyre-real-and-theyr_grkm.1280.jpg",
			Genres: []map[string]string{
				{
					"title": "action",
					"slug":  "action",
				},
				{
					"title": "crime",
					"slug":  "crime",
				},
				{
					"title": "adventure",
					"slug":  "adventure",
				},
			},
		},
		{
			Title: "Virtual Dreamscape",
			Sum:   "Enter a virtual world unlike any other in this cutting-edge VR game. Explore vast landscapes, encounter exotic creatures, and unravel the mysteries of the dream realm. But be warned, the line between reality and fantasy begins to blur...",
			Dev:   "Dreamscape Studios",
			Pub:   "VR Games",
			Image: "https://deadline.com/wp-content/uploads/2018/02/alien-zoo-_pod_-moving-through-landscape-21.jpg",
			Genres: []map[string]string{
				{
					"title": "adventure",
					"slug":  "adventure",
				},
				{
					"title": "fantasy",
					"slug":  "fantasy",
				},
				{
					"title": "simulation",
					"slug":  "simulation",
				},
			},
		},
		{
			Title: "Retro Racer",
			Sum:   "Take a trip back in time with this retro-inspired racing game. Choose your favorite classic car, customize it to your liking, and hit the road in a series of high-speed challenges. Can you become the ultimate Retro Racer?",
			Dev:   "Vintage Games",
			Pub:   "Retro Gaming Co.",
			Image: "https://www.bikebound.com/wp-content/uploads/2019/07/BMW-S1000RR-Retro-Track-Bike-38.jpg",
			Genres: []map[string]string{
				{
					"title": "racing",
					"slug":  "racing",
				},
				{
					"title": "retro",
					"slug":  "retro",
				},
				{
					"title": "action",
					"slug":  "action",
				},
			},
		},
		{
			Title: "Galactic Odyssey",
			Sum:   "Embark on an epic space adventure in this immersive RPG. Explore distant planets, encounter strange creatures, and uncover ancient mysteries as you travel across the galaxy. Will you save the universe from certain doom?",
			Dev:   "Starlight Studios",
			Pub:   "Interstellar Games",
			Image: "https://news.tfw2005.com/wp-content/uploads/sites/10/2020/11/124557574_4576898952380263_2581044460213279492_o.jpg",
			Genres: []map[string]string{
				{
					"title": "RPG",
					"slug":  "RPG",
				},
				{
					"title": "sci-fi",
					"slug":  "sci-fi",
				},
				{
					"title": "adventure",
					"slug":  "adventure",
				},
			},
		},
		{
			Title: "Castle Siege",
			Sum:   "Lead your army to victory in this medieval strategy game. Build your castle, train your troops, and conquer neighboring lands in fierce battles. But be careful, your enemies are not to be underestimated!",
			Dev:   "Kingdom Games",
			Pub:   "Castle Co.",
			Image: "https://cdna.artstation.com/p/assets/images/images/007/923/612/large/maxim-prodanov-castle-siege-max-small.jpg",
			Genres: []map[string]string{
				{
					"title": "strategy",
					"slug":  "strategy",
				},
				{
					"title": "medieval",
					"slug":  "medieval",
				},
				{
					"title": "simulation",
					"slug":  "simulation",
				},
			},
		},
		{
			Title: "Escape the Island",
			Sum:   "Stranded on a deserted island, your only hope is to find a way to escape. Explore the island, scavenge for resources, and build shelter to survive. But watch out for dangerous wildlife and hostile tribes!",
			Dev:   "Survival Games",
			Pub:   "Island Escape Co.",
			Image: "https://assets.reedpopcdn.com/escape-dead-island-review-1417019207828.jpg/BROK/resize/1200x1200%3E/format/jpg/quality/70/escape-dead-island-review-1417019207828.jpg",
			Genres: []map[string]string{
				{
					"title": "survival",
					"slug":  "survival",
				},
				{
					"title": "adventure",
					"slug":  "adventure",
				},
				{
					"title": "simulation",
					"slug":  "simulation",
				},
			},
		},
	}
	group := sync.WaitGroup{}

	genreChan := make(chan games.EmbeddedGameGenre, 100)
	for _, game := range gs {
		group.Add(1)
		go func(game Rand) {
			defer group.Done()

			date := getRandomDate()

			var genres []*games.EmbeddedGameGenre
			for _, genre := range game.Genres {
				newGenre := games.EmbeddedGameGenre{
					Title: genre["title"],
					Slug:  genre["slug"],
				}
				genres = append(genres, &newGenre)

				genreChan <- newGenre
			}

			g := games.Game{
				Title:       game.Title,
				Summary:     game.Sum,
				Id:          primitive.NewObjectID(),
				ReleaseDate: date,
				Developer:   game.Dev,
				Publisher:   game.Pub,
				Rating: games.RatingStats{
					Count: 0,
					Sum:   0,
				},
				CreatedAt: time.Now(),
				IsDeleted: false,
				Image:     game.Image,
				Genres:    genres,
			}
			one, err := mongo.Database("test").Collection("games").InsertOne(context.Background(), g)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(one)

		}(game)
	}
	group.Wait()

	go func() {

		for embGenre := range genreChan {
			//check if genre exists
			var g games.EmbeddedGameGenre
			err := mongo.Database("test").Collection("genres").FindOne(context.Background(), bson.M{"slug": embGenre.Slug}).Decode(&g)

			if err == nil {
				continue
			}

			genre := games.GameGenre{
				Title:     embGenre.Title,
				Slug:      embGenre.Slug,
				CreatedAt: time.Now(),
				IsDeleted: false,
				Desc:      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl nec ultricies lacinia, nunc nisl aliquam nisl, eget aliquam nisl nisl sit amet nisl. Sed euismod, nisl nec ultricies lacinia, nunc nisl aliquam nisl, eget aliquam nisl nisl sit amet nisl.",
			}

			one, err := mongo.Database("test").Collection("genres").InsertOne(context.Background(), genre)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(one)
		}
	}()

	log.Println("done")
}

func getRandomDate() time.Time {
	year := rand.Intn(20) + 2000
	month := rand.Intn(12) + 1
	day := rand.Intn(28) + 1
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// Ping godoc
//
//	@Summary		Show the status of server.
//	@Description	get the status of server.
//
//	@ID				ping
//
//	@Tags			ping
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/api/v1/ping [get]
func ping(ctx *fiber.Ctx) error {
	res := map[string]interface{}{
		"status":  "success",
		"result":  "pong",
		"message": "Server is up and running",
	}

	if err := ctx.JSON(res); err != nil {
		return err
	}

	return nil

}

type JSONResult struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JSONErrorRes struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}
