/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package githubClient

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/google/go-github/v74/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/oauth2"
)

var Client *github.Client

/* Based off of https://github.com/google/go-github/blob/f137c94931a722223df8cc7581a2a3e953ad8d63/README.md */
func init() {
	log.Print("Initializing GitHub client...")

	privateKey := []byte(os.Getenv("TICKETUNE_GITHUB_BOT_PKEY"))
	if len(privateKey) == 0 {
		log.Fatal("TICKETUNE_GITHUB_BOT_PKEY environment variable not set")
	}

	clientId := os.Getenv("TICKETUNE_GITHUB_CLIENT_ID")
	if clientId == "" {
		log.Fatal("TICKETUNE_GITHUB_CLIENT_ID environment variable not set")
	}

	appTokenSource, err := githubauth.NewApplicationTokenSource(clientId, privateKey)
	if err != nil {
		log.Fatal("failed to create GitHub App token source: " + err.Error())
	}

	installationIdStr, err := strconv.ParseInt(os.Getenv("TICKETUNE_INSTALL_ID"), 10, 64)
	if err != nil {
		log.Fatal("TICKETUNE_INSTALL_ID environment variable not set or not an integer")
	}

	installationTokenSource := githubauth.NewInstallationTokenSource(installationIdStr, appTokenSource)

	httpClient := oauth2.NewClient(context.Background(), installationTokenSource)

	Client = github.NewClient(httpClient)

	log.Print("GitHub client initialized.")
}

// Get the installation ID for the github app installation on the org
// This only needs to be done once per install of the app, so once the github app has been
// created and installed, it never needs to be re-run (unless the app is uninstalled and reinstalled).
// In such a case, the code should be manually invoked to print the installation ID, which should then be
// set in the environment variable `TICKETUNE_INSTALL_ID`
/*
func getInstallId(appTokenSource oauth2.TokenSource) int64 {
	// idResponse field, see
	type idResponse struct {
		ID int64 `json:"id"`
	}

	jwtToken, err := appTokenSource.Token()
	if err != nil {
		log.Fatal("failed to get JWT token from app token source: " + err.Error())
	}

	const organization = "pagefaultgames"
	req, err := http.NewRequest("GET", "https://api.github.com/orgs/"+organization+"/installation", nil)
	if err != nil {
		log.Fatal("failed to create request to get installation ID: " + err.Error())
	}

	// Request the installation ID
	// See https://docs.github.com/en/rest/apps/apps?apiVersion=2022-11-28#get-an-organization-installation-for-the-authenticated-app
	client := &http.Client{}
	var bearerToken = "Bearer " + jwtToken.AccessToken
	log.Print("Using bearer token: " + bearerToken)
	req.Header.Set("Authorization", bearerToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("failed to get installation ID: " + err.Error())
	}

	// Read the response and extract the installation ID
	defer resp.Body.Close()
	body, error := io.ReadAll(resp.Body)

	if error != nil {
		log.Fatal("failed to read installation ID response body: " + error.Error())
	}
	var responseJson idResponse
	err = json.Unmarshal(body, &responseJson)
	if err != nil {
		log.Fatal("failed to unmarshal installation ID response body: " + err.Error())
	}

	log.Printf("Installation ID for org %s is %d\n", organization, responseJson.ID)

	return responseJson.ID
}
*/
