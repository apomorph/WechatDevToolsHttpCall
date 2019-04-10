module github.com/bluecatlee/WechatDevToolsHttpCall

go 1.12

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190320223903-b7391e95e576
	golang.org/x/net => github.com/golang/net v0.0.0-20190320064053-1272bf9dcd53
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190321052220-f7bb7a8bee54
	golang.org/x/text => github.com/golang/text v0.3.0
	google.golang.org/appengine => github.com/golang/appengine v1.5.0
)

require (
	github.com/Unknwon/com v0.0.0-20190321035513-0fed4efef755 // indirect
	github.com/go-macaron/binding v0.0.0-20170611065819-ac54ee249c27
	github.com/go-macaron/cache v0.0.0-20151013081102-561735312776
	github.com/go-macaron/inject v0.0.0-20160627170012-d8a0b8677191 // indirect
	github.com/go-macaron/session v0.0.0-20190131233854-0a0a789bf193
	golang.org/x/crypto v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/ini.v1 v1.42.0 // indirect
	gopkg.in/macaron.v1 v1.3.2
	gopkg.in/yaml.v2 v2.2.2
)
