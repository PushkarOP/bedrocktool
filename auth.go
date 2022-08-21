package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sandertv/gophertunnel/minecraft/auth"
	"golang.org/x/oauth2"
)

var G_token_src oauth2.TokenSource

func GetTokenSource() oauth2.TokenSource {
	if G_token_src != nil {
		return G_token_src
	}
	token := get_token()
	G_token_src = auth.RefreshTokenSource(&token)
	new_token, err := G_token_src.Token()
	if err != nil {
		panic(err)
	}
	if !token.Valid() {
		fmt.Println("Refreshed token")
		write_token(new_token)
	}

	return G_token_src
}

var _G_xbl_token *auth.XBLToken

func GetXBLToken(ctx context.Context) *auth.XBLToken {
	if _G_xbl_token != nil {
		return _G_xbl_token
	}
	_token, err := GetTokenSource().Token()
	if err != nil {
		panic(err)
	}
	_G_xbl_token, err = auth.RequestXBLToken(ctx, _token, "https://pocket.realms.minecraft.net/")
	if err != nil {
		panic(err)
	}
	return _G_xbl_token
}

func write_token(token *oauth2.Token) {
	buf, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}
	os.WriteFile(TOKEN_FILE, buf, 0755)
}

func get_token() oauth2.Token {
	var token oauth2.Token
	if _, err := os.Stat(TOKEN_FILE); err == nil {
		f, err := os.Open(TOKEN_FILE)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := json.NewDecoder(f).Decode(&token); err != nil {
			panic(err)
		}
	} else {
		_token, err := auth.RequestLiveToken()
		if err != nil {
			panic(err)
		}
		write_token(_token)
		token = *_token
	}
	return token
}