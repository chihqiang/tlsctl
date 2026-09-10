package register

func GetRegister(kid, hmacEncoded string) IRegister {
	return &Register{
		TermsOfServiceAgreed: true,
		Kid:                  kid,
		HmacEncoded:          hmacEncoded,
	}
}
