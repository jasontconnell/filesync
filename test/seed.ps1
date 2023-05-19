New-Item -ItemType Directory -Force -Path source\infolder
For ($i = 0; $i -lt 100; $i++) {
    $path = "source\infolder\$i.txt"
    Write-Host $path
    "hello $i another" | Out-File -Filepath $path -Force
}