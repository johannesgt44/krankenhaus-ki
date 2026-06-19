package service

import "errors"

var (
	ErrNichtGefunden         = errors.New("krankenhaus nicht gefunden")
	ErrVersionVeraltet       = errors.New("die version ist veraltet")
	ErrEmailBereitsVorhanden = errors.New("die emailadresse ist bereits vorhanden")
)
