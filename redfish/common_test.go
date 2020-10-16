package redfish_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
Redfish example data retreived from a Dell EMC PowerEdge R7415
*/
var (
	//Root example response from a Redfish API v1
	rootRedfishJSON = "{\"@odata.context\":\"/redfish/v1/$metadata#ServiceRoot.ServiceRoot\",\"@odata.id\":\"/redfish/v1\",\"@odata.type\":\"#ServiceRoot.v1_3_0.ServiceRoot\",\"AccountService\":{\"@odata.id\":\"/redfish/v1/Managers/iDRAC.Embedded.1/AccountService\"},\"Chassis\":{\"@odata.id\":\"/redfish/v1/Chassis\"},\"Description\":\"Root Service\",\"EventService\":{\"@odata.id\":\"/redfish/v1/EventService\"},\"Fabrics\":{\"@odata.id\":\"/redfish/v1/Fabrics\"},\"Id\":\"RootService\",\"JsonSchemas\":{\"@odata.id\":\"/redfish/v1/JSONSchemas\"},\"Links\":{\"Sessions\":{\"@odata.id\":\"/redfish/v1/Sessions\"}},\"Managers\":{\"@odata.id\":\"/redfish/v1/Managers\"},\"Name\":\"Root Service\",\"Oem\":{\"Dell\":{\"@odata.context\":\"/redfish/v1/$metadata#DellServiceRoot.DellServiceRoot\",\"@odata.type\":\"#DellServiceRoot.v1_0_0.ServiceRootSummary\",\"IsBranded\":0,\"ManagerMACAddress\":\"d0:94:66:10:04:b3\",\"ServiceTag\":\"90YSGL2\"}},\"Product\":\"Integrated Dell Remote Access Controller\",\"ProtocolFeaturesSupported\":{\"ExpandQuery\":{\"ExpandAll\":true,\"Levels\":true,\"Links\":true,\"MaxLevels\":1,\"NoLinks\":true},\"FilterQuery\":true,\"SelectQuery\":true},\"RedfishVersion\":\"1.4.0\",\"Registries\":{\"@odata.id\":\"/redfish/v1/Registries\"},\"SessionService\":{\"@odata.id\":\"/redfish/v1/SessionService\"},\"Systems\":{\"@odata.id\":\"/redfish/v1/Systems\"},\"Tasks\":{\"@odata.id\":\"/redfish/v1/TaskService\"},\"UpdateService\":{\"@odata.id\":\"/redfish/v1/UpdateService\"}}"

	//Systems redfish collection example
	systemsRedfishJSON = "{\"@odata.context\":\"/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection\",\"@odata.id\":\"/redfish/v1/Systems\",\"@odata.type\":\"#ComputerSystemCollection.ComputerSystemCollection\",\"Description\":\"Collection of Computer Systems\",\"Members\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1\"}],\"Members@odata.count\":1,\"Name\":\"Computer System Collection\"}"

	//System embedded collection example
	systemEmbeddedRedfishJSON = "{\"@odata.context\":\"/redfish/v1/$metadata#ComputerSystem.ComputerSystem\",\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1\",\"@odata.type\":\"#ComputerSystem.v1_5_0.ComputerSystem\",\"Actions\":{\"#ComputerSystem.Reset\":{\"ResetType@Redfish.AllowableValues\":[\"On\",\"ForceOff\",\"ForceRestart\",\"GracefulShutdown\",\"PushPowerButton\",\"Nmi\"],\"target\":\"/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset\"}},\"AssetTag\":\"\",\"Bios\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Bios\"},\"BiosVersion\":\"1.8.7\",\"Boot\":{\"BootOptions\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/BootOptions\"},\"BootOrder\":[\"Boot0003\",\"Boot0000\",\"Boot0001\",\"Boot0004\"],\"BootOrder@odata.count\":4,\"BootSourceOverrideEnabled\":\"Once\",\"BootSourceOverrideMode\":\"UEFI\",\"BootSourceOverrideTarget\":\"None\",\"BootSourceOverrideTarget@Redfish.AllowableValues\":[\"None\",\"Pxe\",\"Floppy\",\"Cd\",\"Hdd\",\"BiosSetup\",\"Utilities\",\"UefiTarget\",\"SDCard\",\"UefiHttp\"],\"UefiTargetBootSourceOverride\":\"\"},\"Description\":\"Computer System which represents a machine (physical or virtual) and the local resources such as memory, cpu and other devices that can be accessed from that machine.\",\"EthernetInterfaces\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/EthernetInterfaces\"},\"HostName\":\"CentOS-Host\",\"HostWatchdogTimer\":{\"FunctionEnabled\":false,\"Status\":{\"State\":\"Disabled\"},\"TimeoutAction\":\"None\"},\"HostingRoles\":[],\"HostingRoles@odata.count\":0,\"Id\":\"System.Embedded.1\",\"IndicatorLED\":\"Off\",\"Links\":{\"Chassis\":[{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1\"}],\"Chassis@odata.count\":1,\"CooledBy\":[{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17%7C%7CFan.Embedded.1\"},{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17%7C%7CFan.Embedded.2\"},{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17%7C%7CFan.Embedded.3\"},{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17%7C%7CFan.Embedded.4\"},{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17%7C%7CFan.Embedded.5\"},{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Sensors/Fans/0x17%7C%7CFan.Embedded.6\"}],\"CooledBy@odata.count\":6,\"ManagedBy\":[{\"@odata.id\":\"/redfish/v1/Managers/iDRAC.Embedded.1\"}],\"ManagedBy@odata.count\":1,\"Oem\":{\"Dell\":{\"BootOrder\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/BootSources\"},\"DellNumericSensorCollection\":{\"@odata.id\":\"/redfish/v1/Dell/Systems/System.Embedded.1/DellNumericSensorCollection\"},\"DellOSDeploymentService\":{\"@odata.id\":\"/redfish/v1/Dell/Systems/System.Embedded.1/DellOSDeploymentService\"},\"DellPresenceAndStatusSensorCollection\":{\"@odata.id\":\"/redfish/v1/Dell/Systems/System.Embedded.1/DellPresenceAndStatusSensorCollection\"},\"DellRaidService\":{\"@odata.id\":\"/redfish/v1/Dell/Systems/System.Embedded.1/DellRaidService\"},\"DellSensorCollection\":{\"@odata.id\":\"/redfish/v1/Dell/Systems/System.Embedded.1/DellSensorCollection\"},\"DellSoftwareInstallationService\":{\"@odata.id\":\"/redfish/v1/Dell/Systems/System.Embedded.1/DellSoftwareInstallationService\"}}},\"PoweredBy\":[{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.1\"},{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Power/PowerSupplies/PSU.Slot.2\"}],\"PoweredBy@odata.count\":2},\"Manufacturer\":\"Dell Inc.\",\"Memory\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Memory\"},\"MemorySummary\":{\"MemoryMirroring\":\"System\",\"Status\":{\"Health\":\"OK\",\"HealthRollup\":\"OK\",\"State\":\"Enabled\"},\"TotalSystemMemoryGiB\":119.2093440},\"Model\":\"PowerEdge R7415\",\"Name\":\"System\",\"NetworkInterfaces\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/NetworkInterfaces\"},\"Oem\":{\"Dell\":{\"DellSystem\":{\"@odata.context\":\"/redfish/v1/$metadata#DellSystem.DellSystem\",\"@odata.id\":\"/redfish/v1/Dell/Systems/System.Embedded.1/DellSystem/System.Embedded.1\",\"@odata.type\":\"#DellSystem.v1_0_0.DellSystem\",\"BIOSReleaseDate\":\"04/02/2019\",\"BaseBoardChassisSlot\":\"NA\",\"BatteryRollupStatus\":\"OK\",\"BladeGeometry\":\"NotApplicable\",\"CMCIP\":null,\"CPURollupStatus\":\"OK\",\"ChassisServiceTag\":\"90YSGL2\",\"ExpressServiceCode\":\"19649475830\",\"FanRollupStatus\":\"OK\",\"IntrusionRollupStatus\":\"OK\",\"LicensingRollupStatus\":\"OK\",\"MaxDIMMSlots\":16,\"MaxPCIeSlots\":5,\"NodeID\":\"90YSGL2\",\"PSRollupStatus\":\"OK\",\"PowerCapEnabledState\":\"Disabled\",\"StorageRollupStatus\":\"OK\",\"SysMemPrimaryStatus\":\"OK\",\"SystemGeneration\":\"14G Monolithic\",\"SystemID\":2039,\"TempRollupStatus\":\"OK\",\"UUID\":\"4c4c4544-0030-5910-8053-b9c04f474c32\",\"VoltRollupStatus\":\"OK\"}}},\"PCIeDevices\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/132-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/64-8\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/64-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/64-7\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/69-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/129-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/64-2\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/64-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-27\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-24\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/64-4\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-25\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-7\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-26\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-4\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-8\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/5-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/131-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-20\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-2\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/0-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/193-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/128-8\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/128-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/128-7\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/134-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/128-2\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/128-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/128-4\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/192-8\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/192-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/192-7\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/192-2\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/192-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevice/192-4\"}],\"PCIeDevices@odata.count\":37,\"PCIeFunctions\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/132-0-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/132-0-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/64-8-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/64-1-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/64-7-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/64-8-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/69-0-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/64-7-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/129-0-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/64-2-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/64-3-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-27-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-24-4\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/64-4-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-25-7\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-7-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-0-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-26-6\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-26-2\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-24-6\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-27-6\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-4-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-8-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-25-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-27-2\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-27-4\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-25-5\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-26-5\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/5-0-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-27-7\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/131-0-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-25-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-7-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-24-2\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-20-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-25-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-25-2\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/129-0-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-27-5\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-8-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-24-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-24-7\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-26-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-1-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-25-6\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-27-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-24-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-24-5\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-24-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-2-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-3-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-26-7\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-27-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-26-4\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-26-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-20-3\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-25-4\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/0-26-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/193-0-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/128-8-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/128-1-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/128-7-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/134-0-2\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/128-8-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/128-7-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/128-2-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/128-3-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/128-4-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/192-8-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/192-1-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/192-7-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/192-7-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/192-8-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/192-2-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/192-3-0\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeFunction/192-4-0\"}],\"PCIeFunctions@odata.count\":76,\"PartNumber\":\"065PKDX30\",\"PowerState\":\"On\",\"ProcessorSummary\":{\"Count\":1,\"LogicalProcessorCount\":32,\"Model\":\"AMD EPYC 7551P 32-Core Processor\",\"Status\":{\"Health\":\"OK\",\"HealthRollup\":\"OK\",\"State\":\"Enabled\"}},\"Processors\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Processors\"},\"SKU\":\"90YSGL2\",\"SecureBoot\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/SecureBoot\"},\"SerialNumber\":\"CNFCP0078T003O\",\"SimpleStorage\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/SimpleStorage/Controllers\"},\"Status\":{\"Health\":\"OK\",\"HealthRollup\":\"OK\",\"State\":\"Enabled\"},\"Storage\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage\"},\"SystemType\":\"Physical\",\"TrustedModules\":[{\"InterfaceType\":\"TPM2_0\",\"Status\":{\"State\":\"Disabled\"}}],\"UUID\":\"4c4c4544-0030-5910-8053-b9c04f474c32\"}"

	//Storage collection example
	storageRedfishJSON = "{\"@odata.context\":\"/redfish/v1/$metadata#StorageCollection.StorageCollection\",\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage\",\"@odata.type\":\"#StorageCollection.StorageCollection\",\"Description\":\"Collection Of Storage entities\",\"Members\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/RAID.Integrated.1-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/AHCI.Embedded.3-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/CPU.1\"}],\"Members@odata.count\":3,\"Name\":\"Storage Collection\"}"

	//Storage collections (3 in total) (Storage here means disk controllers)
	//Storage1
	storage1RedfishJSON = "{\"@odata.context\":\"/redfish/v1/$metadata#Storage.Storage\",\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/RAID.Integrated.1-1\",\"@odata.type\":\"#Storage.v1_7_1.Storage\",\"Description\":\"PERC H740P Mini \",\"Drives\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/RAID.Integrated.1-1/Drives/Disk.Bay.0:Enclosure.Internal.0-1:RAID.Integrated.1-1\"},{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/RAID.Integrated.1-1/Drives/Disk.Bay.1:Enclosure.Internal.0-1:RAID.Integrated.1-1\"}],\"Drives@odata.count\":2,\"Id\":\"RAID.Integrated.1-1\",\"Links\":{\"Enclosures\":[{\"@odata.id\":\"/redfish/v1/Chassis/Enclosure.Internal.0-1:RAID.Integrated.1-1\"},{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1\"}],\"Enclosures@odata.count\":2},\"Name\":\"PERC H740P Mini \",\"Oem\":{\"Dell\":{\"DellController\":{\"@odata.context\":\"/redfish/v1/$metadata#DellController.DellController\",\"@odata.id\":\"/redfish/v1/Dell/Systems/System.Embedded.1/Storage/DellController/RAID.Integrated.1-1\",\"@odata.type\":\"#DellController.v1_1_0.DellController\",\"AlarmState\":\"AlarmNotPresent\",\"BootVirtualDiskFQDD\":null,\"CacheSizeInMB\":8192,\"CachecadeCapability\":\"NotSupported\",\"ConnectorCount\":2,\"ControllerFirmwareVersion\":\"50.5.0-1750\",\"Device\":\"0\",\"DeviceCardDataBusWidth\":\"Unknown\",\"DeviceCardSlotLength\":\"Unknown\",\"DeviceCardSlotType\":\"Unknown\",\"DriverVersion\":\"--NA--\",\"EncryptionCapability\":\"LocalKeyManagementCapable\",\"EncryptionMode\":\"None\",\"KeyID\":null,\"LastSystemInventoryTime\":\"2020-10-15T20:41:24+00:00\",\"LastUpdateTime\":\"2020-10-15T20:41:24+00:00\",\"MaxAvailablePCILinkSpeed\":null,\"MaxPossiblePCILinkSpeed\":null,\"PCISlot\":null,\"PatrolReadState\":\"Stopped\",\"PersistentHotspare\":\"Disabled\",\"RealtimeCapability\":\"Capable\",\"RollupStatus\":\"OK\",\"SASAddress\":\"5D0946600D339900\",\"SecurityStatus\":\"EncryptionCapable\",\"SharedSlotAssignmentAllowed\":\"NotApplicable\",\"SlicedVDCapability\":\"Supported\",\"SupportControllerBootMode\":\"Supported\",\"SupportEnhancedAutoForeignImport\":\"Supported\",\"SupportRAID10UnevenSpans\":\"Supported\",\"T10PICapability\":\"NotSupported\"},\"DellControllerBattery\":{\"@odata.context\":\"/redfish/v1/$metadata#DellControllerBattery.DellControllerBattery\",\"@odata.id\":\"/redfish/v1/Dell/Chassis/System.Embedded.1/DellControllerBattery/Battery.Integrated.1:RAID.Integrated.1-1\",\"@odata.type\":\"#DellControllerBattery.v1_0_0.DellControllerBattery\",\"FQDD\":\"Battery.Integrated.1:RAID.Integrated.1-1\",\"Name\":\"Battery on Integrated RAID Controller 1\",\"PrimaryStatus\":\"OK\",\"RAIDState\":\"Ready\"}}},\"Status\":{\"Health\":\"OK\",\"HealthRollup\":\"OK\",\"State\":\"Enabled\"},\"StorageControllers\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/RAID.Integrated.1-1#/StorageControllers/0\",\"@odata.type\":\"#Storage.v1_7_0.StorageController\",\"Assembly\":{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Assembly\"},\"CacheSummary\":{\"TotalCacheSizeMiB\":7812},\"ControllerRates\":{\"ConsistencyCheckRatePercent\":30,\"RebuildRatePercent\":30},\"FirmwareVersion\":\"50.5.0-1750\",\"Identifiers\":[{\"DurableName\":\"5d0946600d339900\",\"DurableNameFormat\":\"NAA\"}],\"Identifiers@odata.count\":1,\"Links\":{\"PCIeFunctions\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevices/193-0/PCIeFunctions/193-0-0\"}],\"PCIeFunctions@odata.count\":1},\"Manufacturer\":\"DELL\",\"MemberId\":\"0\",\"Model\":\"PERC H740P Mini \",\"Name\":\"PERC H740P Mini \",\"SpeedGbps\":12,\"Status\":{\"Health\":\"OK\",\"HealthRollup\":\"OK\",\"State\":\"Enabled\"},\"SupportedControllerProtocols\":[\"PCIe\"],\"SupportedControllerProtocols@odata.count\":1,\"SupportedDeviceProtocols\":[\"SAS\",\"SATA\"],\"SupportedDeviceProtocols@odata.count\":2,\"SupportedRAIDTypes\":[\"RAID0\",\"RAID1\",\"RAID5\",\"RAID6\",\"RAID10\",\"RAID50\",\"RAID60\"],\"SupportedRAIDTypes@odata.count\":7}],\"StorageControllers@odata.count\":1,\"Volumes\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/RAID.Integrated.1-1/Volumes\"}}"

	//Storage2
	storage2RedfishJSON = "{\"@odata.context\":\"/redfish/v1/$metadata#Storage.Storage\",\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/AHCI.Embedded.3-1\",\"@odata.type\":\"#Storage.v1_7_1.Storage\",\"Description\":\"FCH SATA Controller [AHCI mode]\",\"Drives\":[],\"Drives@odata.count\":0,\"Id\":\"AHCI.Embedded.3-1\",\"Links\":{\"Enclosures\":[{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1\"}],\"Enclosures@odata.count\":1},\"Name\":\"FCH SATA Controller [AHCI mode]\",\"Oem\":{\"Dell\":{\"DellController\":{\"@odata.context\":\"/redfish/v1/$metadata#DellController.DellController\",\"@odata.id\":\"/redfish/v1/Dell/Systems/System.Embedded.1/Storage/DellController/AHCI.Embedded.3-1\",\"@odata.type\":\"#DellController.v1_1_0.DellController\",\"AlarmState\":\"AlarmNotPresent\",\"BootVirtualDiskFQDD\":null,\"CacheSizeInMB\":0,\"CachecadeCapability\":\"NotSupported\",\"ConnectorCount\":0,\"ControllerFirmwareVersion\":null,\"Device\":\"0\",\"DeviceCardDataBusWidth\":\"Unknown\",\"DeviceCardSlotLength\":\"Unknown\",\"DeviceCardSlotType\":\"Unknown\",\"DriverVersion\":null,\"EncryptionCapability\":\"None\",\"EncryptionMode\":\"None\",\"KeyID\":null,\"LastSystemInventoryTime\":\"2020-10-15T13:44:07+00:00\",\"LastUpdateTime\":\"2020-03-17T17:00:35+00:00\",\"MaxAvailablePCILinkSpeed\":null,\"MaxPossiblePCILinkSpeed\":null,\"PCISlot\":null,\"PatrolReadState\":\"Unknown\",\"PersistentHotspare\":\"NotApplicable\",\"RealtimeCapability\":\"Incapable\",\"RollupStatus\":\"Unknown\",\"SASAddress\":\"0\",\"SecurityStatus\":\"EncryptionNotCapable\",\"SharedSlotAssignmentAllowed\":\"NotApplicable\",\"SlicedVDCapability\":\"NotSupported\",\"SupportControllerBootMode\":\"NotSupported\",\"SupportEnhancedAutoForeignImport\":\"NotSupported\",\"SupportRAID10UnevenSpans\":\"NotSupported\",\"T10PICapability\":\"NotSupported\"}}},\"Status\":{\"Health\":null,\"HealthRollup\":null,\"State\":\"Enabled\"},\"StorageControllers\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/AHCI.Embedded.3-1#/StorageControllers/0\",\"@odata.type\":\"#Storage.v1_7_0.StorageController\",\"Assembly\":{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Assembly\"},\"CacheSummary\":{\"TotalCacheSizeMiB\":0},\"ControllerRates\":{\"ConsistencyCheckRatePercent\":null,\"RebuildRatePercent\":null},\"FirmwareVersion\":\"\",\"Identifiers\":[{\"DurableName\":null,\"DurableNameFormat\":null}],\"Identifiers@odata.count\":1,\"Links\":{\"PCIeFunctions\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevices/134-0/PCIeFunctions/134-0-2\"}],\"PCIeFunctions@odata.count\":1},\"Manufacturer\":\"DELL\",\"MemberId\":\"0\",\"Model\":\"FCH SATA Controller [AHCI mode]\",\"Name\":\"FCH SATA Controller [AHCI mode]\",\"Status\":{\"Health\":null,\"HealthRollup\":null,\"State\":\"Enabled\"},\"SupportedControllerProtocols\":[\"PCIe\"],\"SupportedControllerProtocols@odata.count\":1,\"SupportedDeviceProtocols\":[],\"SupportedDeviceProtocols@odata.count\":0,\"SupportedRAIDTypes\":[],\"SupportedRAIDTypes@odata.count\":0}],\"StorageControllers@odata.count\":1,\"Volumes\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/AHCI.Embedded.3-1/Volumes\"}}"

	//Storage3
	storage3RedfishJSON = "{\"@odata.context\":\"/redfish/v1/$metadata#Storage.Storage\",\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/AHCI.Embedded.3-1\",\"@odata.type\":\"#Storage.v1_7_1.Storage\",\"Description\":\"FCH SATA Controller [AHCI mode]\",\"Drives\":[],\"Drives@odata.count\":0,\"Id\":\"AHCI.Embedded.3-1\",\"Links\":{\"Enclosures\":[{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1\"}],\"Enclosures@odata.count\":1},\"Name\":\"FCH SATA Controller [AHCI mode]\",\"Oem\":{\"Dell\":{\"DellController\":{\"@odata.context\":\"/redfish/v1/$metadata#DellController.DellController\",\"@odata.id\":\"/redfish/v1/Dell/Systems/System.Embedded.1/Storage/DellController/AHCI.Embedded.3-1\",\"@odata.type\":\"#DellController.v1_1_0.DellController\",\"AlarmState\":\"AlarmNotPresent\",\"BootVirtualDiskFQDD\":null,\"CacheSizeInMB\":0,\"CachecadeCapability\":\"NotSupported\",\"ConnectorCount\":0,\"ControllerFirmwareVersion\":null,\"Device\":\"0\",\"DeviceCardDataBusWidth\":\"Unknown\",\"DeviceCardSlotLength\":\"Unknown\",\"DeviceCardSlotType\":\"Unknown\",\"DriverVersion\":null,\"EncryptionCapability\":\"None\",\"EncryptionMode\":\"None\",\"KeyID\":null,\"LastSystemInventoryTime\":\"2020-10-15T13:44:07+00:00\",\"LastUpdateTime\":\"2020-03-17T17:00:35+00:00\",\"MaxAvailablePCILinkSpeed\":null,\"MaxPossiblePCILinkSpeed\":null,\"PCISlot\":null,\"PatrolReadState\":\"Unknown\",\"PersistentHotspare\":\"NotApplicable\",\"RealtimeCapability\":\"Incapable\",\"RollupStatus\":\"Unknown\",\"SASAddress\":\"0\",\"SecurityStatus\":\"EncryptionNotCapable\",\"SharedSlotAssignmentAllowed\":\"NotApplicable\",\"SlicedVDCapability\":\"NotSupported\",\"SupportControllerBootMode\":\"NotSupported\",\"SupportEnhancedAutoForeignImport\":\"NotSupported\",\"SupportRAID10UnevenSpans\":\"NotSupported\",\"T10PICapability\":\"NotSupported\"}}},\"Status\":{\"Health\":null,\"HealthRollup\":null,\"State\":\"Enabled\"},\"StorageControllers\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/AHCI.Embedded.3-1#/StorageControllers/0\",\"@odata.type\":\"#Storage.v1_7_0.StorageController\",\"Assembly\":{\"@odata.id\":\"/redfish/v1/Chassis/System.Embedded.1/Assembly\"},\"CacheSummary\":{\"TotalCacheSizeMiB\":0},\"ControllerRates\":{\"ConsistencyCheckRatePercent\":null,\"RebuildRatePercent\":null},\"FirmwareVersion\":\"\",\"Identifiers\":[{\"DurableName\":null,\"DurableNameFormat\":null}],\"Identifiers@odata.count\":1,\"Links\":{\"PCIeFunctions\":[{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/PCIeDevices/134-0/PCIeFunctions/134-0-2\"}],\"PCIeFunctions@odata.count\":1},\"Manufacturer\":\"DELL\",\"MemberId\":\"0\",\"Model\":\"FCH SATA Controller [AHCI mode]\",\"Name\":\"FCH SATA Controller [AHCI mode]\",\"Status\":{\"Health\":null,\"HealthRollup\":null,\"State\":\"Enabled\"},\"SupportedControllerProtocols\":[\"PCIe\"],\"SupportedControllerProtocols@odata.count\":1,\"SupportedDeviceProtocols\":[],\"SupportedDeviceProtocols@odata.count\":0,\"SupportedRAIDTypes\":[],\"SupportedRAIDTypes@odata.count\":0}],\"StorageControllers@odata.count\":1,\"Volumes\":{\"@odata.id\":\"/redfish/v1/Systems/System.Embedded.1/Storage/AHCI.Embedded.3-1/Volumes\"}}"
)

func getReader(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

type responseBuilder struct {
	response http.Response
}

func (r *responseBuilder) Build() http.Response {
	r.response.Proto = "HTTP/1.0"
	r.response.ProtoMajor = 1
	r.response.ProtoMinor = 0
	return r.response
}

func (r *responseBuilder) Status(status string) *responseBuilder {
	r.response.Status = status
	return r
}

func (r *responseBuilder) StatusCode(status int) *responseBuilder {
	r.response.StatusCode = status
	return r
}

func (r *responseBuilder) Body(body string) *responseBuilder {
	r.response.Body = getReader(body)
	return r
}

//TODO IMPLEMENT HEADERS IN THE BUILDER
