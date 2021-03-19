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

Unzip "$((Get-Item -Path ".\").FullName)\wix310-binaries.zip" "C:\\Users\ContainerAdministrator\AppData\Local\Microsoft\WindowsApps"

Write-Output 'Setting Release Version Variable'

$release_version=$(Get-Content .\dist\release_version.txt)

mkdir dist\installer

cd packaging\windows

Write-Output 'GENERATING WIX MSI TEMPLATE'

& 'C:\Program Files\go-msi\go-msi.exe' generate-templates --path wix.json --version $release_version --src ritchie-admin-wix-templates --out $release_version\admin

& 'C:\Program Files\go-msi\go-msi.exe' generate-templates --path wix.json --version $release_version --src ritchie-user-wix-templates --out $release_version\user

Write-Output 'GENERATING MSI INSTALLER'

& 'C:\Program Files\go-msi\go-msi.exe' make --msi ritchiecli.msi --version $release_version --path wix.json --src $release_version\admin

& 'C:\Program Files\go-msi\go-msi.exe' make --msi ritchiecli-user.msi --version $release_version --path wix.json --src $release_version\user

Write-Output 'GENERATING CHOCO INSTALLER'

& 'C:\Program Files\go-msi\go-msi.exe' choco --version $release_version"-ritchie" --input ritchiecli.msi --path wix.json --src $release_version\admin

& 'C:\Program Files\go-msi\go-msi.exe' choco --version $release_version"-ritchie" --input ritchiecli-user.msi --path wix.json --src $release_version\user


Write-Output 'COPYING FILES TO THE RIGHT PLACE'

copy *.msi ..\..\dist\installer

copy *.nupkg ..\..\dist\installer