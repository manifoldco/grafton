module github.com/manifoldco/grafton

require (
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/client9/misspell v0.3.4
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/runtime v0.0.0-20170303002511-e66a4c440602
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.0.0-20170921144055-dc8a684882cf
	github.com/go-swagger/go-swagger v0.0.0-20170414161553-fbc64c262a83
	github.com/go-swagger/scan-repo-boundary v0.0.0-20180623220736-973b3573c013 // indirect
	github.com/go-zoo/bone v0.0.0-20160911183509-fd0aebc74e90
	github.com/gobuffalo/packr v1.30.1
	github.com/golangci/golangci-lint v1.21.0
	github.com/gordonklaus/ineffassign v0.0.0-20180909121442-1003c8bd00dc
	github.com/manifoldco/go-base32 v1.0.3
	github.com/manifoldco/go-base64 v1.0.2
	github.com/manifoldco/go-connector v0.0.3
	github.com/manifoldco/go-jwt v0.1.2
	github.com/manifoldco/go-manifold v0.14.0
	github.com/manifoldco/go-signature v1.0.2
	github.com/manifoldco/promptui v0.3.3-0.20190411181407-35bab80e16a4
	github.com/onsi/gomega v1.7.1
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/tsenart/deadcode v0.0.0-20160724212837-210d2dc333e9
	github.com/urfave/cli/v2 v2.0.0
	golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392
	golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3
	golang.org/x/net v0.0.0-20190923162816-aa69164e4478
)

replace sourcegraph.com/sourcegraph/go-diff v0.5.1 => github.com/sourcegraph/go-diff v0.5.1

go 1.13
