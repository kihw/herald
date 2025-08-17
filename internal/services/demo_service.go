package services

import (
	"database/sql"
	"fmt"
	"time"
	"lol-match-exporter/internal/db"
	"math/rand"
)

// DemoService handles demo data generation
type DemoService struct {
	db *sql.DB
}

// NewDemoService creates a new demo service
func NewDemoService(database *db.Database) *DemoService {
	return &DemoService{
		db: database.DB,
	}
}

// GenerateDemoMatches génère des matches de démonstration pour un utilisateur
func (ds *DemoService) GenerateDemoMatches(userID int, username, tagline string) error {
	// Champions populaires et leurs IDs
	champions := map[int]string{
		1:   "Annie",
		2:   "Olaf", 
		3:   "Galio",
		4:   "Twisted Fate",
		5:   "Xin Zhao",
		7:   "LeBlanc",
		8:   "Vladimir",
		9:   "Fiddlesticks",
		10:  "Kayle",
		11:  "Master Yi",
		12:  "Alistar",
		13:  "Ryze",
		14:  "Sion",
		15:  "Sivir",
		16:  "Soraka",
		17:  "Teemo",
		18:  "Tristana",
		19:  "Warwick",
		20:  "Nunu",
		21:  "Miss Fortune",
		22:  "Ashe",
		23:  "Tryndamere",
		24:  "Jax",
		25:  "Morgana",
		26:  "Zilean",
		27:  "Singed",
		28:  "Evelynn",
		29:  "Twitch",
		30:  "Karthus",
		31:  "Cho'Gath",
		32:  "Amumu",
		33:  "Rammus",
		34:  "Anivia",
		35:  "Shaco",
		36:  "Dr. Mundo",
		37:  "Sona",
		38:  "Kassadin",
		39:  "Irelia",
		40:  "Janna",
		41:  "Gangplank",
		42:  "Corki",
		43:  "Karma",
		44:  "Taric",
		45:  "Veigar",
		48:  "Trundle",
		50:  "Swain",
		51:  "Caitlyn",
		53:  "Blitzcrank",
		54:  "Malphite",
		55:  "Katarina",
		56:  "Nocturne",
		57:  "Maokai",
		58:  "Renekton",
		59:  "Jarvan IV",
		60:  "Elise",
		61:  "Orianna",
		62:  "Wukong",
		63:  "Brand",
		64:  "Lee Sin",
		67:  "Vayne",
		68:  "Rumble",
		69:  "Cassiopeia",
		72:  "Skarner",
		74:  "Heimerdinger",
		75:  "Nasus",
		76:  "Nidalee",
		77:  "Udyr",
		78:  "Poppy",
		79:  "Gragas",
		80:  "Pantheon",
		81:  "Ezreal",
		82:  "Mordekaiser",
		83:  "Yorick",
		84:  "Akali",
		85:  "Kennen",
		86:  "Garen",
		89:  "Leona",
		90:  "Malzahar",
		91:  "Talon",
		92:  "Riven",
		96:  "Kog'Maw",
		98:  "Shen",
		99:  "Lux",
		101: "Xerath",
		102: "Shyvana",
		103: "Ahri",
		104: "Graves",
		105: "Fizz",
		106: "Volibear",
		107: "Rengar",
		110: "Varus",
		111: "Nautilus",
		112: "Viktor",
		113: "Sejuani",
		114: "Fiora",
		115: "Ziggs",
		117: "Lulu",
		119: "Draven",
		120: "Hecarim",
		121: "Kha'Zix",
		122: "Darius",
		126: "Jayce",
		127: "Lissandra",
		131: "Diana",
		133: "Quinn",
		134: "Syndra",
		136: "Aurelion Sol",
		141: "Kayn",
		142: "Zoe",
		143: "Zyra",
		145: "Kai'Sa",
		147: "Seraphine",
		150: "Gnar",
		154: "Zac",
		157: "Yasuo",
		161: "Vel'Koz",
		163: "Taliyah",
		164: "Camille",
		166: "Akshan",
		200: "Bel'Veth",
		201: "Braum",
		202: "Jhin",
		203: "Kindred",
		221: "Zeri",
		222: "Jinx",
		223: "Tahm Kench",
		234: "Viego",
		235: "Senna",
		236: "Lucian",
		238: "Zed",
		240: "Kled",
		245: "Ekko",
		246: "Qiyana",
		254: "Vi",
		266: "Aatrox",
		267: "Nami",
		268: "Azir",
		350: "Yuumi",
		360: "Samira",
		412: "Thresh",
		420: "Illaoi",
		421: "Rek'Sai",
		427: "Ivern",
		429: "Kalista",
		432: "Bard",
		516: "Ornn",
		517: "Sylas",
		518: "Neeko",
		523: "Aphelios",
		526: "Rell",
		555: "Pyke",
		777: "Yone",
		875: "Sett",
		876: "Lillia",
		887: "Gwen",
		888: "Renata Glasc",
		895: "Nilah",
		897: "K'Sante",
		901: "Smolder",
	}

	positions := []string{"TOP", "JUNGLE", "MIDDLE", "BOTTOM", "UTILITY"}
	queueIds := []int{420, 440} // Ranked Solo/Duo, Ranked Flex

	// Générer 15 matches de démonstration
	for i := 0; i < 15; i++ {
		// Données du match
		matchID := fmt.Sprintf("DEMO_%d_%d", userID, i+1)
		gameCreation := time.Now().Add(-time.Duration(i*24) * time.Hour) // Un match par jour les 15 derniers jours
		gameDuration := 1800 + rand.Intn(1200) // Entre 30 et 50 minutes
		queueID := queueIds[rand.Intn(len(queueIds))]
		
		// Insérer le match
		matchQuery := `
			INSERT INTO matches (
				match_id, game_creation, game_duration, game_mode, game_type, 
				game_version, map_id, platform_id, queue_id, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
			ON CONFLICT (match_id) DO NOTHING
			RETURNING id`
		
		var matchDBID int
		err := ds.db.QueryRow(matchQuery, matchID, gameCreation, gameDuration,
			"CLASSIC", "MATCHED_GAME", "14.16.1", 11, "EUW1", queueID).Scan(&matchDBID)
		
		if err != nil {
			if err == sql.ErrNoRows {
				// Le match existe déjà, passer au suivant
				continue
			}
			return fmt.Errorf("failed to insert match: %w", err)
		}

		// Générer 10 participants (2 équipes de 5)
		for participantID := 1; participantID <= 10; participantID++ {
			// Choisir un champion aléatoire
			championKeys := make([]int, 0, len(champions))
			for k := range champions {
				championKeys = append(championKeys, k)
			}
			championID := championKeys[rand.Intn(len(championKeys))]
			championName := champions[championID]
			
			// Équipe (1 ou 2)
			teamID := 1
			if participantID > 5 {
				teamID = 2
			}
			
			// Position aléatoire
			position := positions[rand.Intn(len(positions))]
			
			// Est-ce le joueur principal?
			isMainPlayer := participantID == 1 // Le joueur principal est toujours le participant 1
			
			// Générer des stats réalistes
			var kills, deaths, assists, goldEarned, minions, visionScore, damage int
			var win bool
			
			if isMainPlayer {
				// Stats légèrement meilleures pour le joueur principal
				kills = 3 + rand.Intn(12)        // 3-15 kills
				deaths = 1 + rand.Intn(8)        // 1-8 deaths  
				assists = 2 + rand.Intn(15)      // 2-17 assists
				goldEarned = 12000 + rand.Intn(8000) // 12k-20k gold
				minions = 140 + rand.Intn(100)   // 140-240 CS
				visionScore = 15 + rand.Intn(35)  // 15-50 vision
				damage = 15000 + rand.Intn(20000) // 15k-35k damage
				
				// 60% de chance de gagner
				win = rand.Float32() < 0.6
			} else {
				kills = rand.Intn(15)
				deaths = 1 + rand.Intn(10)
				assists = rand.Intn(20)
				goldEarned = 8000 + rand.Intn(12000)
				minions = 80 + rand.Intn(160)
				visionScore = 5 + rand.Intn(45)
				damage = 8000 + rand.Intn(25000)
				
				// Si c'est la même équipe que le joueur principal, même résultat
				if (participantID <= 5 && teamID == 1) || (participantID > 5 && teamID == 2) {
					// Même équipe que le participant 1
					win = rand.Float32() < 0.6
				} else {
					// Équipe adverse
					win = rand.Float32() < 0.4
				}
			}
			
			// PUUID factice pour le joueur principal, autres PUUID aléatoires
			puuid := fmt.Sprintf("DEMO-PUUID-%d-%d", userID, participantID)
			if isMainPlayer {
				puuid = fmt.Sprintf("MAIN-PUUID-%s-%s", username, tagline)
			}
			
			// Insérer le participant
			participantQuery := `
				INSERT INTO match_participants (
					match_id, participant_id, puuid, champion_id, champion_name, team_id,
					position, kills, deaths, assists, gold_earned, total_minions_killed,
					vision_score, damage_dealt_to_champions, win, created_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW())
				ON CONFLICT (match_id, participant_id) DO NOTHING`

			_, err = ds.db.Exec(participantQuery,
				matchDBID, participantID, puuid, championID, championName, teamID,
				position, kills, deaths, assists, goldEarned, minions, visionScore,
				damage, win)
			if err != nil {
				return fmt.Errorf("failed to insert participant: %w", err)
			}
		}
	}

	return nil
}
