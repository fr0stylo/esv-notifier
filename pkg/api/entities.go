package api

type TimeListEntry struct {
	PractitionerID        interface{} `json:"practitionerId"`
	ProfessionCode        interface{} `json:"professionCode"`
	ProfessionName        interface{} `json:"professionName"`
	HealthcareServiceID   int         `json:"healthcareServiceId"`
	HealthcareServiceName string      `json:"healthcareServiceName"`
	FundType              struct {
		Type      string      `json:"type"`
		Name      string      `json:"name"`
		IsEnabled interface{} `json:"isEnabled"`
	} `json:"fundType"`
	ReferralNeed struct {
		Type      string      `json:"type"`
		Name      string      `json:"name"`
		IsEnabled interface{} `json:"isEnabled"`
	} `json:"referralNeed"`
	OrganizationID   int    `json:"organizationId"`
	OrganizationName string `json:"organizationName"`
	EarliestTime     int64  `json:"earliestTime"`
	PeriodTimesCount int    `json:"periodTimesCount"`
}

type Specialist struct {
	ID            int         `json:"id"`
	SpcAsmAk      interface{} `json:"spcAsmAk"`
	FullName      string      `json:"fullName"`
	SpcAsmVardas  string      `json:"spcAsmVardas"`
	SpcAsmPavarde string      `json:"spcAsmPavarde"`
	PhoneNo       interface{} `json:"phoneNo"`
	Email         interface{} `json:"email"`
	Profession    struct {
		Code interface{} `json:"code"`
		Name interface{} `json:"name"`
		ID   interface{} `json:"id"`
	} `json:"profession"`
	Professions interface{} `json:"professions"`
	Institution struct {
		IstgID           int         `json:"istgId"`
		IstgPavadinimas  interface{} `json:"istgPavadinimas"`
		NameExt          interface{} `json:"nameExt"`
		Address          interface{} `json:"address"`
		JarCode          interface{} `json:"jarCode"`
		PhoneNo          interface{} `json:"phoneNo"`
		SveidraNumber    interface{} `json:"sveidraNumber"`
		MunicipalityCode interface{} `json:"municipalityCode"`
		MunicipalityID   interface{} `json:"municipalityId"`
	} `json:"institution"`
	Availability              interface{} `json:"availability"`
	AvailableUntil            interface{} `json:"availableUntil"`
	RoomNumber                interface{} `json:"roomNumber"`
	OrganizationBranchAddress interface{} `json:"organizationBranchAddress"`
	HealthcareServices        interface{} `json:"healthcareServices"`
	StampNumber               interface{} `json:"stampNumber"`
	Reception                 interface{} `json:"reception"`
	RoleCode                  interface{} `json:"roleCode"`
	QualificationDescription  interface{} `json:"qualificationDescription"`
}

type Response[T any] struct {
	Data  []T       `json:"data"`
	Links LinksData `json:"links"`
	Meta  MetaData  `json:"meta"`
}

type LinksData struct {
	Self string      `json:"self"`
	Prev interface{} `json:"prev"`
	Next string      `json:"next"`
}

type MetaData struct {
	TotalPages int `json:"totalPages"`
}
