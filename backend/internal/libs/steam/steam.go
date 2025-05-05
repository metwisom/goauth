package steam

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"goauth/internal/config"
	"regexp"
	"strings"
)

type User struct {
	SteamID                  string `json:"steamid"`
	CommunityVisibilityState int    `json:"communityvisibilitystate"`
	ProfileState             int    `json:"profilestate"`
	PersonaName              string `json:"personaname"`
	CommentPermission        int    `json:"commentpermission"`
	ProfileURL               string `json:"profileurl"`
	Avatar                   string `json:"avatar"`
	AvatarMedium             string `json:"avatarmedium"`
	AvatarFull               string `json:"avatarfull"`
	AvatarHash               string `json:"avatarhash"`
	LastLogoff               int64  `json:"lastlogoff"`
	PersonaState             int    `json:"personastate"`
	RealName                 string `json:"realname"`
	PrimaryClanID            string `json:"primaryclanid"`
	TimeCreated              int64  `json:"timecreated"`
	PersonaStateFlags        int    `json:"personastateflags"`
	LocCountryCode           string `json:"loccountrycode"`
	LocStateCode             string `json:"locstatecode"`
	LocCityID                int    `json:"loccityid"`
}

type Response struct {
	Players []User `json:"players"`
}

type ApiResponse struct {
	Response Response `json:"response"`
}

// ValidateSteamResponse Валидация ответа Steam
func ValidateSteamResponse(query *fasthttp.Args) (bool, error) {
	query.Set("openid.mode", "check_authentication")
	_, resp, err := fasthttp.Post(nil, "https://steamcommunity.com/openid/login", query)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(resp), "is_valid:true"), nil
}

// ExtractSteamID Извлечение Steam ID из claimed_id
func ExtractSteamID(claimedID string) string {
	re := regexp.MustCompile(`https?://steamcommunity\.com/openid/id/(\d+)`)
	matches := re.FindStringSubmatch(claimedID)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// GetSteamUser Получение данных пользователя из Steam API
func GetSteamUser(steamID string) (User, error) {
	url := fmt.Sprintf(
		"https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s",
		config.Config.SteamKey, steamID,
	)
	_, resp, err := fasthttp.Get(nil, url)
	if err != nil {
		return User{}, err
	}

	var steamResponse ApiResponse
	if err := json.Unmarshal(resp, &steamResponse); err != nil {
		return User{}, err
	}

	if len(steamResponse.Response.Players) == 0 {
		return User{}, fmt.Errorf("no player data found")
	}

	return steamResponse.Response.Players[0], nil
}
