package utils

import (
	"testing"
)

func TestSubnetExpand(t *testing.T) {
	failMsg := "SubnetExpand fails on valid input"
	shouldFailMsg := "SubnetExpand does not fail when it should"
	wrongIpMsg := "SubnetExpand returns incorrect IP"
	wrongNetmaskMsg := "SubnetExpand returns incorrect netmask"

	subnet := "10.46.0.0/24"
	ip, netmask, err := SubnetExpand(subnet)

	if (err != nil) { t.Error(failMsg) }
	if (ip != "10.46.0.0") { t.Error(wrongIpMsg) }
	if (netmask != "255.255.255.0") { t.Error(wrongNetmaskMsg) }

	subnet = "192.168.255.255/32"
	ip, netmask, err = SubnetExpand(subnet)

	if (err != nil) { t.Error(failMsg) }
	if (ip != "192.168.255.255") { t.Error(wrongIpMsg) }
	if (netmask != "255.255.255.255") { t.Error(wrongNetmaskMsg) }

	subnet = "10.46.0.0/0"
	ip, netmask, err = SubnetExpand(subnet)

	if (err != nil) { t.Error(failMsg) }
	if (ip != "10.46.0.0") { t.Error(wrongIpMsg) }
	if (netmask != "0.0.0.0") { t.Error(wrongNetmaskMsg) }

	subnet = "10.46.1.100/18"
	ip, netmask, err = SubnetExpand(subnet)

	if (err != nil) { t.Error(failMsg) }
	if (ip != "10.46.1.100") { t.Error(wrongIpMsg) }
	if (netmask != "255.255.192.0") { t.Error(wrongNetmaskMsg) }

	subnet = "this.is.not.a/subnet"
	ip, netmask, err = SubnetExpand(subnet)

	if (err == nil) { t.Error(shouldFailMsg) }
}

