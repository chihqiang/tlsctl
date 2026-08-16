package register

func GetRegister(kid, hmacEncoded string) (register IRegister) {
	if kid != "" && hmacEncoded != "" {
		register = &EABRegister{
			TermsOfServiceAgreed: true,
			Kid:                  kid,
			HmacEncoded:          hmacEncoded,
		}
	} else {
		register = &Register{}
	}
	return register
}
