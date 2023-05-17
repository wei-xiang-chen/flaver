package utils

import "flaver/globals"

func GetRunTimeEnv() string {
	return globals.GetViper().GetString("server.env")
}