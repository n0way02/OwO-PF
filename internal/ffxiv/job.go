package ffxiv

type Job int

const (
	GNB Job = iota
	PLD
	GLD
	DRK
	WAR
	MRD
	SCH
	ACN // Arcanist
	SGE
	AST
	WHM
	CNJ
	SAM
	DRG
	NIN
	MNK
	RPR
	VPR
	BRD
	MCH
	DNC
	BLM
	BLU
	SMN
	PCT
	RDM
	LNC
	PUG
	ROG
	THM
	ARC // Archer
	Unknown
)

func JobFromAbbreviation(abbreviation string) Job {
	switch abbreviation {
	case "GNB":
		return GNB
	case "PLD":
		return PLD
	case "GLD":
		return GLD
	case "DRK":
		return DRK
	case "WAR":
		return WAR
	case "MRD":
		return MRD
	case "SCH":
		return SCH
	case "ACN":
		return ACN
	case "SGE":
		return SGE
	case "AST":
		return AST
	case "WHM":
		return WHM
	case "CNJ":
		return CNJ
	case "SAM":
		return SAM
	case "DRG":
		return DRG
	case "NIN":
		return NIN
	case "MNK":
		return MNK
	case "RPR":
		return RPR
	case "VPR":
		return VPR
	case "BRD":
		return BRD
	case "MCH":
		return MCH
	case "DNC":
		return DNC
	case "BLM":
		return BLM
	case "BLU":
		return BLU
	case "SMN":
		return SMN
	case "PCT":
		return PCT
	case "RDM":
		return RDM
	case "LNC":
		return LNC
	case "PUG":
		return PUG
	case "ROG":
		return ROG
	case "THM":
		return THM
	case "ARC":
		return ARC
	}
	return Unknown
}

func (j Job) Emoji() string {
	switch j {
	case GNB:
		return "<:gunbreaker:1358896099376763184>"
	case PLD:
		return "<:paladin:1358896204733481040>"
	case GLD:
		return "<:gladiator:1358896091197603890>"
	case DRK:
		return "<:darkknight:1358896052295438529>"
	case WAR:
		return "<:warrior:1358896427689840861>"
	case MRD:
		return "<:marauder:1358896170130342132>"
	case SCH:
		return "<:scholar:1358896292432187412>"
	case ACN:
		return "<:arcanist:1358895974990483638>"
	case SGE:
		return "<:sage:1358896272718958834>"
	case AST:
		return "<:astrologian:1358895997022896148>"
	case WHM:
		return "<:whitemage:1358896436615581748>"
	case CNJ:
		return "<:conjurer:1358896035459760148>"
	case SAM:
		return "<:samurai:1358896282495615169>"
	case DRG:
		return "<:dragoon:1358896083442467056>"
	case NIN:
		return "<:ninja:1358896189675929792>"
	case MNK:
		return "<:monk:1358896178489594089>"
	case RPR:
		return "<:reaper:1358896239596409169>"
	case VPR:
		return "<:VPR:1358896415929012294>"
	case BRD:
		return "<:bard:1358896009975173380>"
	case MCH:
		return "<:machinist:1358896158616846416>"
	case DNC:
		return "<:dancer:1358896043294462243>"
	case BLM:
		return "<:blackmage:1358896019672400103>"
	case BLU:
		return "<:bluemage:1358896028496953384>"
	case SMN:
		return "<:summoner:1358896302129418321>"
	case RDM:
		return "<:redmage:1358896249448825013>"
	case PCT:
		return "<:PCT:1358896217182175268>"
	case LNC:
		return "<:lancer:1358896139855859933>"
	case PUG:
		return "<:pugilist:1358896227508424734>"
	case ROG:
		return "<:rogue:1358896262211961115>"
	case THM:
		return "<:thaumaturge:1358896352431443968>"
	case ARC:
		return "<:archer:1358895984205107294>"
	}
	return "<:DOL:1358896066845610216>"
}
