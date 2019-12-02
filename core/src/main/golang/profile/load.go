package profile

import (
	"net"

	adapters "github.com/Dreamacro/clash/adapters/outbound"
	"github.com/Dreamacro/clash/component/auth"
	trie "github.com/Dreamacro/clash/component/domain-trie"
	"github.com/Dreamacro/clash/config"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/dns"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/log"
	"github.com/Dreamacro/clash/tunnel"

	"github.com/kr328/cfa/tun"
)

// LoadDefault - load default configure
func LoadDefault() {
	defaultC := &config.Config{
		General: &config.General{
			Port:               0,
			SocksPort:          0,
			RedirPort:          0,
			Authentication:     []string{},
			AllowLan:           false,
			BindAddress:        "*",
			Mode:               tunnel.Direct,
			LogLevel:           log.SILENT,
			ExternalController: "",
			ExternalUI:         "",
			Secret:             "",
		},
		DNS: &config.DNS{
			Enable:     false,
			IPv6:       false,
			NameServer: []dns.NameServer{},
			Fallback:   []dns.NameServer{},
			FallbackFilter: config.FallbackFilter{
				GeoIP:  false,
				IPCIDR: []*net.IPNet{},
			},
			Listen:       "",
			EnhancedMode: dns.NORMAL,
			FakeIPRange:  nil,
		},
		Experimental: &config.Experimental{
			IgnoreResolveFail: false,
		},
		Hosts:   trie.New(),
		Rules:   []constant.Rule{},
		Users:   []auth.AuthUser{},
		Proxies: map[string]constant.Proxy{},
	}

	reject := adapters.NewProxy(adapters.NewReject())
	direct := adapters.NewProxy(adapters.NewDirect())
	global, _ := adapters.NewSelector("GLOBAL", []constant.Proxy{direct})

	defaultC.Proxies["DIRECT"] = direct
	defaultC.Proxies["REJECT"] = reject
	defaultC.Proxies["GLOBAL"] = adapters.NewProxy(global)

	tun.ResetDnsRedirect()

	executor.ApplyConfig(defaultC, true)
}

// LoadFromFile - load file
func LoadFromFile(path string) error {
	cfg, err := executor.ParseWithPath(path)
	if err != nil {
		return err
	}

	executor.ApplyConfig(cfg, true)

	tun.ResetDnsRedirect()

	return nil
}