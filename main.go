package main

import "github.com/rs/zerolog/log"

func main() {
	ReadEnvVars()

	CloneRepository()

	err := Promote()
	if err != nil {
		log.Panic().Err(err).Msgf("there was problem when promoting, stop execution")
	}

  PushEnv()
}
