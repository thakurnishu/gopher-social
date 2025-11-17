package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/thakurnishu/gopher-social/internal/store"
)

var usernames = []string{
	"silverhawk", "nightwanderer", "techsage", "stormrider", "emberwolf",
	"neonpixel", "cloudseeker", "ironquill", "bytecrafter", "frostbyte",
	"wildsparrow", "lunarcrest", "shadowmint", "glimmertrail", "rustleaves",
	"codemancer", "silentvoyager", "skybreaker", "marblestone", "crypticorbit",
	"solarbloom", "mistymeadow", "orbitalpulse", "deepcurrent", "echoingstep",
	"pixelrune", "softwhisper", "lightforge", "crystalbranch", "windshaper",
	"ironember", "cobaltpeak", "starlitstride", "brightmesa", "novastrider",
	"mysticflare", "hazelroot", "quartzline", "flintshadow", "goldenquartz",
	"riverscribe", "thundercrest", "nightcycle", "brightbound", "voidrunner",
	"treetopmuse", "ashenpath", "silentcircuit", "meadowstride",
}


var titles = []string{
	"Shadows of the Silent Forest", "Echoes Beyond the Horizon", "The Last Ember of Dawn",
	"Whispers in the Copper Sky", "Unraveled Threads of Fate", "Journey Through the Forgotten Isles",
	"Fragments of a Broken Crown", "Secrets Beneath the Starlit Sea", "The Clockmaker’s Paradox",
	"Empire of Ash and Glass", "Footsteps Through Midnight Rain", "Rise of the Fallen Lanterns",
	"Letters From the Edge of Nowhere", "A Song for the Wounded Earth", "The Vanishing Light of Solace",
	"Labyrinth of Rust and Dreams", "Beyond the Hollow Gate", "The Art of Chasing Thunder",
	"Chronicles of the Crimson Path", "Where the Wild Winds Sleep",
}

var contents = []string{
	"Under the fading light of dusk, the world seemed to hold its breath as secrets long forgotten resurfaced.",
	"In the quiet moments between storms, the truth revealed itself in whispers carried by the wind.",
	"Every journey begins with a single uncertain step, but courage is born from walking anyway.",
	"The river carved its story through the land, shaping everything it touched across centuries.",
	"Among the scattered ruins, echoes of old civilizations shimmered with memories of greatness.",
	"A lone traveler stood against the endless horizon, searching for meaning beyond the known.",
	"The scent of rain lingered in the air, promising renewal, change, and the birth of something new.",
	"In the heart of the ancient forest, shadows danced with an eerie sort of elegance.",
	"Time moved strangely in the old city, where the past and present intertwined effortlessly.",
	"Every decision left a mark, shaping destinies in ways no one could fully understand.",
	"The storm clouds gathered with an intensity that signaled something more than mere weather.",
	"Dreams often hold truths we are not yet ready to face in the waking world.",
	"Beneath the ocean’s surface lay secrets too vast and ancient to ever be fully uncovered.",
	"The fire crackled warmly, offering comfort even as uncertainty loomed beyond the walls.",
	"With each passing season, the land transformed, telling a story of resilience and rebirth.",
	"Nightfall brought with it a quiet sense of mystery, urging wanderers to tread carefully.",
	"In forgotten libraries, dust-covered books held wisdom that shaped the world long ago.",
	"The path ahead twisted unpredictably, yet hope lit the way for those determined enough.",
	"Old legends whispered through the mountains, carried by winds older than any living soul.",
	"Even in darkness, a spark of light can guide the lost toward a new beginning.",
}

var tags = []string{
	"adventure", "mystery", "fantasy", "technology", "science", "travel",
	"philosophy", "storytelling", "nature", "history", "creativity", "culture",
	"innovation", "art", "writing", "discovery", "exploration", "ideas",
	"insight", "reflection",
}

 var comments = []string{
	"Really enjoyed reading this!", "This gave me a lot to think about.",
	"Amazing insight—thanks for sharing.", "I didn’t expect that, very cool.",
	"Great perspective, I totally agree.", "This was surprisingly relatable.",
	"Looking forward to more posts like this.", "Nicely written and easy to follow.",
	"I learned something new today.", "Wow, this hit deeper than I expected.",
	"Simple but powerful message.", "Thanks for explaining it so clearly.",
	"I love the way you express these ideas.", "This was exactly what I needed to hear.",
	"Your writing style is fantastic.", "Interesting angle, didn’t consider that before.",
	"Thought-provoking and well-structured.", "Please make a part two!",
	"Saved this—so good!", "Beautifully articulated.",
}


func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating users:", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComments(50000, posts) 
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding has completed")
}

func generateComments(num int, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := range cms {

		post := posts[rand.Intn(len(posts))]
		cms[i] = &store.Comment{
			PostID: post.ID,
			UserID: post.UserID,
			Content: comments[rand.Intn(len(comments))], 
		}
	}
	return cms
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := range posts {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID: user.ID,
			Title: titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))], 
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := range users {
		username := usernames[i%len(usernames)] + fmt.Sprintf("%d", i)
		users[i] = &store.User{
			Username: username,
			Email:    fmt.Sprintf("%s@example.com", username),
			Password: "password123",
		}
	}
	return users
}
