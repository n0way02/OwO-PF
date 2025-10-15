package ffxiv

import (
	"reflect"
)

type Role int

const (
	DPS Role = iota
	Healer
	Tank
	Empty
)

type Roles struct {
	Roles []Role
}

func (rs Roles) Emoji() string {
	if reflect.DeepEqual(rs.Roles, []Role{DPS}) {
		return "<:dps:1358896075263705300>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{Healer}) {
		return "<:healer:1358896124315959376>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{Tank}) {
		return "<:tank:1358896312380035343>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{DPS, Healer}) {
		return "<:healerdps:1358896132360765510>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{DPS, Tank}) {
		return "<:tankdps:1358896322349891604>"
	}
	if reflect.DeepEqual(rs.Roles, []Role{Healer, Tank}) {
		return "<:tankhealer:1358896332026155159>"
	}

	if reflect.DeepEqual(rs.Roles, []Role{Healer, Tank, DPS}) {
		return "<:tankhealerdps:1358896342050803762>"
	}

	return "<:ANY:1358895963510411466>"
}
