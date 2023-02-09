param prefix string = 'devbox'
param size string = 'Standard_D2s_v4'
param username string = ''
param sshPublicKey string = ''
param location string = resourceGroup().location
param subnet string = '172.22.254.248/29'
param allowedIP string = ''

var subnetName = prefix
var uniquePrefix = '${prefix}-${uniqueString(prefix, resourceGroup().id)}'

resource pip 'Microsoft.Network/publicIpAddresses@2019-02-01' = {
  name: '${uniquePrefix}-pip'
  location: location
  sku: {
    name: 'Basic'
  }
  properties: {
    publicIPAllocationMethod: 'Static'
  }
}

resource nsg 'Microsoft.Network/networkSecurityGroups@2019-02-01' = {
  name: '${uniquePrefix}-nsg'
  location: location
  properties: {
    securityRules: [
      {
        name: 'SSH'
        properties: {
          priority: 320
          protocol: 'Tcp'
          access: 'Allow'
          direction: 'Inbound'
          sourceAddressPrefix: allowedIP
          sourcePortRange: '*'
          destinationAddressPrefix: '*'
          destinationPortRange: '22'
        }
      }
    ]
  }
}

resource vnet 'Microsoft.Network/virtualNetworks@2019-09-01' = {
  name: '${uniquePrefix}-vnet'
  location: location
  properties: {
    addressSpace: {
      addressPrefixes: [
        subnet
      ]
    }
    subnets: [
      {
        name: subnetName
        properties: {
          addressPrefix: subnet
        }
      }
    ]
  }
}

resource nic 'Microsoft.Network/networkInterfaces@2018-10-01' = {
  name: '${uniquePrefix}-nic'
  location: location
  properties: {
    ipConfigurations: [
      {
        name: 'ipconfig1'
        properties: {
          subnet: {
            id: '${vnet.id}/subnets/${subnetName}'
          }
          privateIPAllocationMethod: 'Dynamic'
          publicIPAddress: {
            id: pip.id
          }
        }
      }
    ]
    networkSecurityGroup: {
      id: nsg.id
    }
  }
}

resource vm 'Microsoft.Compute/virtualMachines@2020-06-01' = {
  name: '${uniquePrefix}-vm'
  location: location
  properties: {
    hardwareProfile: {
      vmSize: size
    }
    storageProfile: {
      osDisk: {
        createOption: 'FromImage'
        managedDisk: {
          storageAccountType: 'Premium_LRS'
        }
      }
      imageReference: {
        publisher: 'Canonical'
        offer: 'UbuntuServer'
        sku: '18.04-LTS'
        version: 'latest'
      }
    }
    networkProfile: {
      networkInterfaces: [
        {
          id: nic.id
        }
      ]
    }
    securityProfile: {}
    osProfile: {
      computerName: prefix
      adminUsername: username
      linuxConfiguration: {
        disablePasswordAuthentication: true
        ssh: {
          publicKeys: [
            {
              path: '/home/${username}/.ssh/authorized_keys'
              keyData: sshPublicKey
            }
          ]
        }
      }
    }
    priority: 'Spot'
    evictionPolicy: 'Deallocate'
    billingProfile: {
      maxPrice: -1
    }
  }
}
