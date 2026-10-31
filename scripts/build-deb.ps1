$ErrorActionPreference = 'Stop'
$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$Version = if ($env:VERSION) { $env:VERSION } else { '1.0.0' }
$Arch = if ($env:ARCH) { $env:ARCH } else { 'amd64' }
$Output = Join-Path $Root "dist\liveryapi_${Version}_linux_${Arch}"

Remove-Item $Output -Recurse -Force -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force "$Output\DEBIAN", "$Output\usr\bin", "$Output\lib\systemd\system", "$Output\var\lib\liveryapi" | Out-Null

$env:GOOS = 'linux'
$env:GOARCH = $Arch
$env:CGO_ENABLED = '0'
& 'D:\environment\go\bin\go.exe' build -trimpath -ldflags '-s -w' -o "$Output\usr\bin\liveryapi" $Root
if ($LASTEXITCODE -ne 0) { throw 'Go Linux build failed' }

Copy-Item "$Root\debian\control" "$Output\DEBIAN\control"
Copy-Item "$Root\debian\postinst" "$Output\DEBIAN\postinst"
Copy-Item "$Root\debian\prerm" "$Output\DEBIAN\prerm"
Copy-Item "$Root\debian\liveryapi.service" "$Output\lib\systemd\system\liveryapi.service"
(Get-Content "$Output\DEBIAN\control") -replace '^Version: .*', "Version: $Version" | Set-Content "$Output\DEBIAN\control" -NoNewline

Write-Host "Linux binary created at $Output\usr\bin\liveryapi"
Write-Host 'To create the .deb on Debian/Ubuntu, copy this project and run:'
Write-Host '  chmod +x scripts/build-deb.sh'
Write-Host '  ./scripts/build-deb.sh'
