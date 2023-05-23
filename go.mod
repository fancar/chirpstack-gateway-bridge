module github.com/brocaar/chirpstack-gateway-bridge

go 1.20

// replace github.com/brocaar/chirpstack-api/go/v3 => github.com/fancar/chirpstack-api/go/v3 v3.8.2-0.20210706144249-37fbff59495d

//v3.8.2-0.20210630184849-72e4ff07581a

replace github.com/brocaar/chirpstack-api/go/v3 => /home/fancar/dev/iot/api/go

require (
	github.com/brocaar/chirpstack-api/go/v3 v3.11.0
	github.com/brocaar/lorawan v0.0.0-20210809075358-95fc1667572e
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/go-zeromq/zmq4 v0.7.0
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/golang/protobuf v1.5.3
	github.com/goreleaser/goreleaser v0.106.0
	github.com/goreleaser/nfpm v0.11.0
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.5.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0
	github.com/shirou/gopsutil v3.21.3+incompatible
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.8.3
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
)

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/apex/log v1.1.0 // indirect
	github.com/aws/aws-sdk-go v1.27.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blakesmith/ar v0.0.0-20150311145944-8bd4349a67f2 // indirect
	github.com/caarlos0/ctrlc v1.0.0 // indirect
	github.com/campoy/unique v0.0.0-20180121183637-88950e537e7e // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-zeromq/goczmq/v4 v4.2.2 // indirect
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jacobsa/crypto v0.0.0-20190317225127-9f44e2d11115 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/kamilsk/retry/v4 v4.0.0 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mattn/go-colorable v0.0.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/mattn/go-zglob v0.0.0-20180803001819-2ea3427bfa53 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.15.0 // indirect
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/afero v1.1.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	golang.org/x/sync v0.2.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.9.1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/protobuf v1.29.1 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
