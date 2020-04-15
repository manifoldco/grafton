module github.com/manifoldco/grafton

go 1.13

require (
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200108200545-475eaeb16496
	github.com/client9/misspell v0.3.4
	github.com/go-openapi/errors v0.19.3
	github.com/go-openapi/runtime v0.19.9
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.6
	github.com/go-openapi/validate v0.0.0-20170921144055-dc8a684882cf
	github.com/go-swagger/go-swagger v0.2.1-0.20170911151148-68ea41dcf206
	github.com/go-zoo/bone v0.0.0-20160911183509-fd0aebc74e90
	github.com/gobuffalo/packr v1.30.1
	github.com/gobuffalo/packr/v2 v2.8.0 // indirect
	github.com/golangci/golangci-lint v1.23.8
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/gordonklaus/ineffassign v0.0.0-20190601041439-ed7b1b5ee0f8
	github.com/gorilla/context v1.1.1 // indirect
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/magefile/mage v1.9.0
	github.com/manifoldco/go-base32 v1.0.4
	github.com/manifoldco/go-base64 v1.0.3
	github.com/manifoldco/go-connector v0.1.0
	github.com/manifoldco/go-jwt v0.2.0
	github.com/manifoldco/go-manifold v0.15.0
	github.com/manifoldco/go-signature v1.0.3
	github.com/manifoldco/logo v0.11.507
	github.com/manifoldco/promptui v0.3.3-0.20190411181407-35bab80e16a4
	github.com/onsi/gomega v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.5.0
	github.com/toqueteos/webbrowser v1.1.0 // indirect
	github.com/tsenart/deadcode v0.0.0-20160724212837-210d2dc333e9
	github.com/tylerb/graceful v1.2.15 // indirect
	github.com/urfave/cli/v2 v2.2.0
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
)

// We're stuck on an old go-openapi version, lock it.
replace (
	github.com/go-openapi/errors => github.com/go-openapi/errors v0.0.0-20170426151106-03cfca65330d
	github.com/go-openapi/loads => github.com/go-openapi/loads v0.0.0-20170520182102-a80dea3052f0
	github.com/go-openapi/runtime => github.com/go-openapi/runtime v0.0.0-20171117174610-4ab0c2c86df5
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.0.0-20171105025112-eef1d9a2e601
	github.com/go-openapi/strfmt => github.com/go-openapi/strfmt v0.0.0-20170822153411-610b6cacdcde
	github.com/go-openapi/swag => github.com/go-openapi/swag v0.0.0-20170606142751-f3f9494671f9
	github.com/go-openapi/validate => github.com/go-openapi/validate v0.0.0-20170921144055-dc8a684882cf
	github.com/go-swagger/go-swagger => github.com/go-swagger/go-swagger v0.2.1-0.20170911151148-68ea41dcf206
)
