param variables_firewallRules_copyIndex_name ? /* TODO: fill in correct type */
param variables_firewallRules_copyIndex_startIPAddress ? /* TODO: fill in correct type */
param variables_firewallRules_copyIndex_endIPAddress ? /* TODO: fill in correct type */
param serverGroupName string

resource serverGroupName_variables_firewallRules_copyIndex_name_name 'Microsoft.DBforPostgreSQL/serverGroupsv2/firewallRules@2020-10-05-privatepreview' = {
  name: '${serverGroupName}/${variables_firewallRules_copyIndex_name[copyIndex()].name}'
  properties: {
    startIpAddress: variables_firewallRules_copyIndex_startIPAddress[copyIndex()].startIPAddress
    endIpAddress: variables_firewallRules_copyIndex_endIPAddress[copyIndex()].endIPAddress
  }
}