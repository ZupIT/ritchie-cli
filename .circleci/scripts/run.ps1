$url = "https://github.com/mh-cbon/go-msi/releases/download/1.0.2/go-msi-amd64.msi"
$output = "$((Get-Item -Path ".\").FullName)\go-msi-amd64.msi"
$start_time = Get-Date

$wc = New-Object System.Net.WebClient
$wc.DownloadFile($url, $output)

Write-Output "Time taken: $((Get-Date).Subtract($start_time).Seconds) second(s)"

"$((Get-Item -Path ".\").FullName)\go-msi-amd64.msi" /quiet

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

Unzip "$((Get-Item -Path ".\").FullName)\wix310-binaries.zip" "$((Get-Item -Path ".\").FullName)\packaging\windows\wix310"

dir "$((Get-Item -Path ".\").FullName)"
