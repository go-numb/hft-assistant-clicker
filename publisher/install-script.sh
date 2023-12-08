# Goのインストール
Invoke-WebRequest https://golang.org/dl/go1.20.2.windows-amd64.msi -OutFile go.msi
Start-Process msiexec.exe -Wait -ArgumentList '/i go.msi /quiet'
$env:GOPATH = "$env:USERPROFILE\go"
$env:PATH += ";$env:GOPATH\bin"

# Gitのインストール
Invoke-WebRequest https://github.com/git-for-windows/git/releases/download/v2.34.0.windows.1/Git-2.34.0-64-bit.exe -OutFile git.exe
Start-Process git.exe -Wait -ArgumentList '/SP- /VERYSILENT /NORESTART /NOICONS'

# WebView2のインストール
Invoke-WebRequest https://go.microsoft.com/fwlink/p/?LinkId=2124703 -OutFile webview2.exe
Start-Process webview2.exe -Wait -ArgumentList '/install /quiet /norestart'