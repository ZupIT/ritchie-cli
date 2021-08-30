Write-Output 'DOWLOADING GO-MSI INSTALLER'

choco install go-msi

Write-Output 'DOWLOADING WIX FILES'

$url = "http://wixtoolset.org/downloads/v3.10.3.3007/wix310-binaries.zip"
$output = "$((Get-Item -Path ".\").FullName)\wix310-binaries.zip"

$wc = New-Object System.Net.WebClient
$wc.DownloadFile($url, $output)

Add-Type -AssemblyName System.IO.Compression.FileSystem
function Unzip
{
    param([string]$zipfile, [string]$outpath)

    [System.IO.Compression.ZipFile]::ExtractToDirectory($zipfile, $outpath)
}

Write-Output 'EXTRACTING WIX FILES TO PATH'

Unzip "$((Get-Item -Path ".\").FullName)\wix310-binaries.zip" "C:\Users\runneradmin\AppData\Local\Microsoft\WindowsApps"

Write-Output 'Setting Release Version Variable'
$release_version = "${Env:RELEASE_VERSION}"

Write-Output 'Create folders'

mkdir dist\installer
mkdir dist\installer\admin
mkdir dist\installer\user
mkdir dist\out
mkdir dist\out\admin
mkdir dist\out\user
cd .github\scripts\windows

Write-Output 'GENERATING WIX MSI TEMPLATE'

$path = "D:\\a\\ritchie-cli\\ritchie-cli\\.github\\scripts\\windows\\wix.json"
$src_admin = "D:\\a\\ritchie-cli\\ritchie-cli\\.github\\scripts\\windows\\ritchie-admin-wix-templates"
$src_user = "D:\\a\\ritchie-cli\\ritchie-cli\\.github\\scripts\\windows\\ritchie-user-wix-templates"
$out_admin = "D:\\a\\ritchie-cli\\ritchie-cli\\dist\\out\\admin"
$out_user = "D:\\a\\ritchie-cli\\ritchie-cli\\dist\\out\\user"
$installer_admin = "D:\\a\\ritchie-cli\\ritchie-cli\\dist\\installer\\admin"
$installer_user = "D:\\a\\ritchie-cli\\ritchie-cli\\dist\\installer\\user"
$installer_admin_name = "ritchie_$($release_version)_windows_x86_64.msi"
$installer_user_name = "ritchie-user_$($release_version)_windows_x86_64.msi"

& 'C:\Program Files\go-msi\go-msi.exe' generate-templates --path $path --version $release_version --src $src_admin --out $out_admin
& 'C:\Program Files\go-msi\go-msi.exe' generate-templates --path $path --version $release_version --src $src_user --out $out_user

Write-Output 'GENERATING MSI INSTALLER'

& 'C:\Program Files\go-msi\go-msi.exe' make --msi $installer_admin\$installer_admin_name --version $release_version --path $path --src $src_admin --out $out_admin
& 'C:\Program Files\go-msi\go-msi.exe' make --msi $installer_user\$installer_user_name --version $release_version --path $path --src $src_user --out $out_user
