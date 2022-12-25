rule smtp_phishing_detection{
    meta:
        author = "İsmet Güngör"
        description = "Phishing detection"
    strings:
        $changePassword = "change password" nocase
        $password = "password" nocase
        $tckn = "tckn" nocase
        $clickhere = /click here/ nocase wide
        $sale = "sale" nocase
        
        $test = "test" nocase
    condition:
        any of them

}

rule content : mail {
	meta:
		author = "A.Sanchez <asanchez@koodous.com>"
		description = "Detects scam emails with phishing attachment."
		test1 = "email/eml/transferencia1.eml"
		test2 = "email/eml/transferencia2.eml"

	strings:
		$subject = "Asunto: Justificante de transferencia" nocase
		$body = "Adjunto justificante de transferencia"
	condition:
		all of them
}

rule attachment : mail {
	meta:
		author = "A.Sanchez <asanchez@koodous.com>"
		description = "Detects scam emails with phishing attachment."
		test1 = "email/eml/transferencia1.eml"
		test2 = "email/eml/transferencia2.eml"

	strings:
		$filename = "filename=\"scan001.pdf.html\""
		$pleaseEnter = "NTAlNkMlNjUlNjElNzMlNjUlMjAlNjUlNkUlNzQlNjUlNzIlMjAlN" // Please enter 
		$emailReq = "NkQlNjUlNkUlNzQlMkUlNjklNkUlNjQlNjUlNzglMzIlMkUlNDUlNkQlNjElNjklNkMlM0I" // ment.index2.Email;
		$pAssign = "NzAlMjAlM0QlMjAlNjQlNkYlNjMlNzUlNkQlNjUlNkUlNzQlMkUlNjklNkUlNjQlNjUl" // p = document.inde
		
	condition:
		all of them
}