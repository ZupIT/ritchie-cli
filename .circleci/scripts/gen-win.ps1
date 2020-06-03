
Write-Output 'DOWLOADING GO-MSI INSTALLER'

$url = "https://github.com/mh-cbon/go-msi/releases/download/1.0.2/go-msi-amd64.msi"
$output = "$((Get-Item -Path ".\").FullName)\go-msi-amd64.msi"
$start_time = Get-Date

$wc = New-Object System.Net.WebClient
$wc.DownloadFile($url, $output)

Write-Output "Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"

Write-Output 'INSTALLING GO-MSI'

Start-Process msiexec.exe -Wait -ArgumentList "/I $output /quiet"

Write-Output 'DOWLOADING WIX FILES'

$url = "http://wixtoolset.org/downloads/v3.10.3.3007/wix310-binaries.zip"
$output = "$((Get-Item -Path ".\").FullName)\wix310-binaries.zip"
$start_time = Get-Date

$wc = New-Object System.Net.WebClient
$wc.DownloadFile($url, $output)

Write-Output "Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"

Add-Type -AssemblyName System.IO.Compression.FileSystem
function Unzip
{
    param([string]$zipfile, [string]$outpath)

    [System.IO.Compression.ZipFile]::ExtractToDirectory($zipfile, $outpath)
}

Write-Output 'EXTRACTING WIX FILES TO PATH'

Unzip "$((Get-Item -Path ".\").FullName)\wix310-binaries.zip" "C:\\Users\circleci\AppData\Local\Microsoft\WindowsApps"

Write-Output 'Setting Release Version Variable'

$release_version=$(Get-Content .\release_version.txt)

mkdir dist\installer

copy LICENSE packaging/windows

cd packaging\windows

Write-Output 'GENERATING MSI TEAM INSTALLER'

& 'C:\Program Files\go-msi\go-msi.exe' make --msi ritchiecliteam.msi --version $release_version --path wix-team.json

Write-Output 'GENERATING CHOCO TEAM INSTALLER'

& 'C:\Program Files\go-msi\go-msi.exe' choco --version $release_version"-team" --input ritchiecliteam.msi --path wix-team.json

Write-Output 'GENERATING MSI SINGLE INSTALLER'

& 'C:\Program Files\go-msi\go-msi.exe' make --msi ritchieclisingle.msi --version $release_version --path wix-single.json

Write-Output 'GENERATING CHOCO SINGLE INSTALLER'

& 'C:\Program Files\go-msi\go-msi.exe' choco --version $release_version"-single" --input ritchieclisingle.msi --path wix-single.json

Write-Output 'COPYING FILES TO THE RIGHT PLACE'

copy ritchie* ..\..\dist\installer

copy *.nupkg ..\..\dist\installer