module github.com/sanathp/statusok

go 1.12

require (
	github.com/Sirupsen/logrus v1.4.0
	github.com/codegangsta/cli v1.21.0
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/influxdata/influxdb v1.7.8
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mailgun/mailgun-go v2.0.0+incompatible
	github.com/onsi/ginkgo v1.9.0 // indirect
	github.com/onsi/gomega v1.6.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586 // indirect
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
)

replace (
	github.com/Sirupsen/logrus v1.0.5 => github.com/sirupsen/logrus v1.0.5
	github.com/Sirupsen/logrus v1.3.0 => github.com/Sirupsen/logrus v1.0.6
	github.com/Sirupsen/logrus v1.4.0 => github.com/sirupsen/logrus v1.0.6

	github.com/codegangsta/cli v1.21.0 => github.com/urfave/cli v1.19.1

	github.com/mailgun/mailgun-go v2.0.0+incompatible => github.com/mailgun/mailgun-go v1.1.0
)
