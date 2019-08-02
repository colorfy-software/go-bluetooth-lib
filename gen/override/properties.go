package override

func GetPropertiesOverride(iface string) (map[string]string, bool) {
	if props, ok := PropertyTypes[iface]; ok {
		return props, ok
	}
	return map[string]string{}, false
}

var PropertyTypes = map[string]map[string]string{
	"org.bluez.Device1": map[string]string{
		// "ServiceData":      "map[string]dbus.Variant",
		"ManufacturerData": "map[uint16]dbus.Variant",
	},
	"org.bluez.GattCharacteristic1": map[string]string{
		"Value":       "[]byte `dbus:\"emit\"`",
		"Descriptors": "[]dbus.ObjectPath",
	},
	"org.bluez.GattDescriptor1": map[string]string{
		"Value":          "[]byte `dbus:\"emit\"`",
		"Characteristic": "dbus.ObjectPath",
	},
	"org.bluez.GattService1": map[string]string{
		"Characteristics": "[]dbus.ObjectPath `dbus:\"emit\"`",
		"Device":          "[]dbus.ObjectPath `dbus:\"ignore=isService\"`",
		"IsService":       "bool `dbus:\"ignore\"`",
	},
	"org.bluez.LEAdvertisement1": map[string]string{
		// dbus type: (yv) dict of byte variant (array of bytes)
		"Data": "map[byte]interface{}",
		// dbus type: (qv) dict of uint16 variant (array of bytes)
		"ManufacturerData": "map[uint16]interface{}",
		// dbus type: (s[v]) dict of string variant (array of bytes)
		"ServiceData": "map[string]interface{}",
	},
}
