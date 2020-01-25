package thermometer



import (
   "sync"
   "github.com/muka/go-bluetooth/bluez"
   "github.com/muka/go-bluetooth/util"
   "github.com/muka/go-bluetooth/props"
   "github.com/godbus/dbus"
)

var ThermometerWatcher1Interface = "org.bluez.ThermometerWatcher1"


// NewThermometerWatcher1 create a new instance of ThermometerWatcher1
//
// Args:
// - servicePath: unique name
// - objectPath: freely definable
func NewThermometerWatcher1(servicePath string, objectPath dbus.ObjectPath) (*ThermometerWatcher1, error) {
	a := new(ThermometerWatcher1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: ThermometerWatcher1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(ThermometerWatcher1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
ThermometerWatcher1 Health Thermometer Watcher hierarchy

*/
type ThermometerWatcher1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*ThermometerWatcher1Properties
	watchPropertiesChannel chan *dbus.Signal
}

// ThermometerWatcher1Properties contains the exposed properties of an interface
type ThermometerWatcher1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

//Lock access to properties
func (p *ThermometerWatcher1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *ThermometerWatcher1Properties) Unlock() {
	p.lock.Unlock()
}



// Close the connection
func (a *ThermometerWatcher1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return ThermometerWatcher1 object path
func (a *ThermometerWatcher1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return ThermometerWatcher1 dbus client
func (a *ThermometerWatcher1) Client() *bluez.Client {
	return a.client
}

// Interface return ThermometerWatcher1 interface
func (a *ThermometerWatcher1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *ThermometerWatcher1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

	if a.objectManagerSignal == nil {
		if a.objectManager == nil {
			om, err := bluez.GetObjectManager()
			if err != nil {
				return nil, nil, err
			}
			a.objectManager = om
		}

		s, err := a.objectManager.Register()
		if err != nil {
			return nil, nil, err
		}
		a.objectManagerSignal = s
	}

	cancel := func() {
		if a.objectManagerSignal == nil {
			return
		}
		a.objectManagerSignal <- nil
		a.objectManager.Unregister(a.objectManagerSignal)
		a.objectManagerSignal = nil
	}

	return a.objectManagerSignal, cancel, nil
}


// ToMap convert a ThermometerWatcher1Properties to map
func (a *ThermometerWatcher1Properties) ToMap() (map[string]interface{}, error) {
	return props.ToMap(a), nil
}

// FromMap convert a map to an ThermometerWatcher1Properties
func (a *ThermometerWatcher1Properties) FromMap(props map[string]interface{}) (*ThermometerWatcher1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an ThermometerWatcher1Properties
func (a *ThermometerWatcher1Properties) FromDBusMap(props map[string]dbus.Variant) (*ThermometerWatcher1Properties, error) {
	s := new(ThermometerWatcher1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// ToProps return the properties interface
func (a *ThermometerWatcher1) ToProps() bluez.Properties {
	return a.Properties
}

// GetWatchPropertiesChannel return the dbus channel to receive properties interface
func (a *ThermometerWatcher1) GetWatchPropertiesChannel() chan *dbus.Signal {
	return a.watchPropertiesChannel
}

// SetWatchPropertiesChannel set the dbus channel to receive properties interface
func (a *ThermometerWatcher1) SetWatchPropertiesChannel(c chan *dbus.Signal) {
	a.watchPropertiesChannel = c
}

// GetProperties load all available properties
func (a *ThermometerWatcher1) GetProperties() (*ThermometerWatcher1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *ThermometerWatcher1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *ThermometerWatcher1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *ThermometerWatcher1) GetPropertiesSignal() (chan *dbus.Signal, error) {

	if a.propertiesSignal == nil {
		s, err := a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
		if err != nil {
			return nil, err
		}
		a.propertiesSignal = s
	}

	return a.propertiesSignal, nil
}

// Unregister for changes signalling
func (a *ThermometerWatcher1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *ThermometerWatcher1) WatchProperties() (chan *bluez.PropertyChanged, error) {
	return bluez.WatchProperties(a)
}

func (a *ThermometerWatcher1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	return bluez.UnwatchProperties(a, ch)
}




/*
MeasurementReceived 
			This callback gets called when a measurement has been
			scanned in the thermometer.

			Measurement:

				int16 Exponent:
				int32 Mantissa:

					Exponent and Mantissa values as
					extracted from float value defined by
					IEEE-11073-20601.


*/
func (a *ThermometerWatcher1) MeasurementReceived(measurement map[string]interface{}) error {
	
	return a.client.Call("MeasurementReceived", 0, measurement).Store()
	
}

