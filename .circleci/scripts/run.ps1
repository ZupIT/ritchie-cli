$url = "https://github.com/mh-cbon/go-msi/releases/download/1.0.2/go-msi-amd64.msi"
$output = "$PSScriptRoot\go-msi-amd64.msi"
$start_time = Get-Date

Import-Module BitsTransfer
Start-BitsTransfer -Source $url -Destination $output
#OR
Start-BitsTransfer -Source $url -Destination $output -Asynchronous

Write-Output "Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"

$url = "http://wixtoolset.org/downloads/v3.10.3.3007/wix310-binaries.zip"
$output = "$PSScriptRoot\wix310-binaries.zip"
$start_time = Get-Date

Import-Module BitsTransfer
Start-BitsTransfer -Source $url -Destination $output
#OR
Start-BitsTransfer -Source $url -Destination $output -Asynchronous

Write-Output "Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"