package auth

import "encoding/json"

type Ownership struct {
	Items []struct {
		Name      string
		Signature string
	}
	Signature string
	KeyId     string
}

func (acc *MinecraftAuthentication) OwnsGame() bool {
	resp, err := acc.Client.GET("https://api.minecraftservices.com/entitlements/mcstore", []string{"Authorization"}, []string{"Bearer " + acc.MinecraftToken})
	if err != nil {
		return false
	}

	var j Ownership
	if err := json.Unmarshal(resp, &j); err != nil {
		return false
	}

	prod, game := false, false
	for _, j := range j.Items {
		if j.Name == "product_minecraft" {
			prod = true
		} else if j.Name == "game_minecraft" {
			game = true
		}
	}

	return prod && game
}
