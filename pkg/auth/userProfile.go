package auth

import "encoding/json"

type MinecraftProfile struct {
	Id    string
	Name  string
	Skins []struct {
		Id      string
		State   string
		Url     string
		Variant string
		Alias   string
	}
	Capes []interface{}
}

func (ac *MinecraftAuthentication) GetProfile() (*MinecraftProfile, error) {
	resp, err := ac.Client.GET(
		"https://api.minecraftservices.com/minecraft/profile",
		[]string{"Authorization"},
		[]string{"Bearer " + ac.MinecraftToken},
	)
	if err != nil {
		return nil, err
	}

	var j MinecraftProfile
	if err := json.Unmarshal(resp, &j); err != nil {
		return nil, err
	}

	return &j, nil
}
