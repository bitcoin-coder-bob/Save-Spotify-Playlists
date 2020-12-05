package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	// "ioutil"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

func main() {

	// Follow this link and after you authenticate, it will redirect, and the url will have a code query parameter. Paste the code in codeQueryParam variable below.
	// https://accounts.spotify.com/authorize?client_id=&response_type=code&redirect_uri=https%3A%2F%2Fexample.com%2Fcallback&scope=playlist-read-private%20playlist-modify-private%20playlist-modify-public%20user-read-email&state=34fFs29kd09

	codeQueryParam := ""
	base64Auth := ""
	// a := NewApp()
	// go a.Run()
	// defer a.Shutdown()
	// params := url.Values{}
	// fmt.Println(params.Encode())
	// addItmesToPlaylistURL := "https://api.spotify.com/v1/playlists/"
	myPlaylistsURL := "https://api.spotify.com/v1/me/playlists"
	tokenURL := "https://accounts.spotify.com/api/token"
	playlistItemsURL := "https://api.spotify.com/v1/playlists/"
	userURL := "https://api.spotify.com/v1/me"

	// req1, err := http.NewRequest("GET", getCodeURL, nil)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// client1 := &http.Client{}
	// response1, err := client1.Do(req1)
	// if err != nil {
	// 	fmt.Printf("error getting response from client when streaming from location:\n%v\n%s", response, err.Error())
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// content, err := ioutil.ReadAll(response1.H)
	// if err != nil {
	// 	fmt.Printf("error reading data from response body:\n%s", err.Error())
	// 	return
	// }
	// fmt.Println(response1.Header)

	/////////////////////

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", codeQueryParam)
	data.Set("redirect_uri", "https://example.com/callback")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("Error making post to token url: %v\n", err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Authorization", "Basic "+base64Auth)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("error getting response from client when streaming from location:\n%v\n%s", response, err.Error())
		fmt.Println(err.Error())
		return
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Printf("failed getting token: status code %d\n", response.StatusCode)
		return
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error reading data from response body:\n%s", err.Error())
		return
	}
	strContent := string(content)
	accessToken := gjson.Get(strContent, "access_token").String()
	fmt.Printf("Access token: %v\n", accessToken)
	scope := gjson.Get(strContent, "scope").String()
	fmt.Printf("Scope : %v\n", scope)
	expires := gjson.Get(strContent, "expires_in").String()
	fmt.Printf("expires: %v\n", expires)
	refreshToken := gjson.Get(strContent, "refresh_token").String()
	fmt.Printf("reresh token: %v\n", refreshToken)

	//now getting the playlist info
	DELETELATERTOKEN := accessToken
	req, err = http.NewRequest("GET", myPlaylistsURL, nil)
	if err != nil {
		fmt.Printf("Error making post to playlist url: %v\n", err.Error())
		return
	}
	req.Header.Add("Authorization", "Bearer "+DELETELATERTOKEN)
	client = &http.Client{}
	response, err = client.Do(req)
	if err != nil {
		fmt.Printf("error getting response from client when streaming from location:\n%v\n%s", response, err.Error())
		fmt.Println(err.Error())
		return
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Printf("failed getting user's playlists: status code %d\n", response.StatusCode)
		return
	}

	content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error reading data from response body:\n%s", err.Error())
		return
	}
	playlists := gjson.GetBytes(content, "items").Array()
	releaseRadarID := ""
	for _, playlist := range playlists {
		if playlist.Get("name").String() == "Release Radar" {
			releaseRadarID = playlist.Get("id").String()
		}
	}

	// Now get all track IDs from playlist
	req, err = http.NewRequest("GET", playlistItemsURL+releaseRadarID, nil)
	if err != nil {
		fmt.Printf("Error making post to playlist url: %v\n", err.Error())
		return
	}
	req.Header.Add("Authorization", "Bearer "+DELETELATERTOKEN)
	client = &http.Client{}
	response, err = client.Do(req)
	if err != nil {
		fmt.Printf("error getting response from client when streaming from location:\n%v\n%s", response, err.Error())
		fmt.Println(err.Error())
		return
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Printf("failed getting tracks from playlist: status code %d\n", response.StatusCode)
		return
	}

	content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error reading data from response body:\n%s", err.Error())
		return
	}
	tracks := gjson.GetBytes(content, "tracks.items").Array()
	trackURIsToAdd := []string{}
	for _, track := range tracks {
		trackURIsToAdd = append(trackURIsToAdd, track.Get("track.uri").String())
	}
	//get user's id
	req, err = http.NewRequest("GET", userURL, nil)
	if err != nil {
		fmt.Printf("Error making post to playlist url: %v\n", err.Error())
		return
	}
	req.Header.Add("Authorization", "Bearer "+DELETELATERTOKEN)
	client = &http.Client{}
	response, err = client.Do(req)
	if err != nil {
		fmt.Printf("error getting response from client when streaming from location:\n%v\n%s", response, err.Error())
		fmt.Println(err.Error())
		return
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Printf("failed getting user info: status code %d\n", response.StatusCode)
		return
	}

	content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error reading data from response body:\n%s", err.Error())
		return
	}
	userID := gjson.GetBytes(content, "id").String()
	fmt.Printf("user id: %v\n", userID)

	// make a new playlist
	createPlaylistURL := "https://api.spotify.com/v1/users/" + userID + "/playlists"
	bodyData := "{\"name\":\"RR 11/30/2020\", \"public\":false}"

	req, err = http.NewRequest("POST", createPlaylistURL, strings.NewReader(bodyData))
	if err != nil {
		fmt.Printf("Error making post to playlist url: %v\n", err.Error())
		return
	}
	req.Header.Add("Authorization", "Bearer "+DELETELATERTOKEN)
	req.Header.Add("Content-Type", "application/json")
	client = &http.Client{}
	response, err = client.Do(req)
	if err != nil {
		fmt.Printf("error getting response from client when streaming from location:\n%v\n%s", response, err.Error())
		fmt.Println(err.Error())
		return
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Printf("failed response making playlist: status code %d\n", response.StatusCode)
		return
	}

	content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error reading data from response body:\n%s", err.Error())
		return
	}
	newPlaylistID := gjson.GetBytes(content, "id").String()

	// add tracks to new playlist
	stringifiedTrackURIs := ""
	for num, track := range trackURIsToAdd {
		stringifiedTrackURIs += "\"" + track + "\""
		if len(trackURIsToAdd)-1 != num {
			stringifiedTrackURIs += ","
		}
	}
	urisData := "{\"uris\":[" + stringifiedTrackURIs + "]}"

	addTrackURL := "https://api.spotify.com/v1/playlists/" + newPlaylistID + "/tracks"
	req, err = http.NewRequest("POST", addTrackURL, strings.NewReader(urisData))
	if err != nil {
		fmt.Printf("Error adding tracks to playlist url: %v\n", err.Error())
		return
	}
	req.Header.Add("Authorization", "Bearer "+DELETELATERTOKEN)
	req.Header.Add("Content-Type", "application/json")
	client = &http.Client{}
	response, err = client.Do(req)
	if err != nil {
		fmt.Printf("error getting response from client when streaming from location:\n%v\n%s", response, err.Error())
		fmt.Println(err.Error())
		return
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Printf("failed response adding tracks to playlist: status code %d\n", response.StatusCode)
		return
	}

	content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error reading data from response body:\n%s", err.Error())
		return
	}
	fmt.Println("added tracks to new playlist!")
}
