package main

import "github.com/rs/zerolog/log"

func main() {
	ReadEnvVars()

	err := CloneRepostories()
	if err != nil {
		log.Panic().Err(err).Msgf("there was problem when cloning, stop execution")
	}

	err = Promote()
	if err != nil {
		log.Panic().Err(err).Msgf("there was problem when promoting, stop execution")
	}

	err = PushEnvs()
	if err != nil {
		log.Panic().Err(err).Msgf("there was problem when pushing, stop execution")
	}

	log.Info().Msgf("successfully promoted the helm chart")
}
