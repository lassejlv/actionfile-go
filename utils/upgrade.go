package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Masterminds/semver"
	"github.com/rs/zerolog/log"
)

func Upgrade() {

	currentVersion := LoadVersion()

	if currentVersion == " " {
		log.Warn().Msg("Could not detect version")
		os.Exit(1)
	}

	// Get the latest github release and match it to the current version

	log.Info().Msg("Checking for updates...")

	url := "https://api.github.com/repos/lassejlv/actionfile-go/releases/latest"
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		panic(err)
	}

	log.Info().Msgf("Latest version: %s", release.TagName)
	log.Info().Msgf("Current version: %s", currentVersion)

	v, err := semver.NewVersion(release.TagName)

	if err != nil {
		panic(err)
	}

	isOutdated, err := semver.NewConstraint("> " + currentVersion)

	if err != nil {
		panic(err)
	}

	if isOutdated.Check(v) {
		log.Info().Msg("Upgrading to " + release.TagName)
		RunCmd("go install github.com/lassejlv/action@" + release.TagName)
		log.Info().Msg("Upgrade complete")
		log.Info().Msg(fmt.Sprintf("Read more about this release at https://github.com/lassejlv/action/releases/tag/%s", release.TagName))
	} else {
		log.Info().Msg("You are already on the latest version")
	}

}
