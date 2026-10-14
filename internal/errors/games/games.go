package games

import "errors"

var ErrInvalidGameID = errors.New("id du jeu invalide")
var ErrGameNotFound = errors.New("jeu introuvable")
