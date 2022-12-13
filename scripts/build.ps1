$ErrorActionPreference = "Stop"
$packages = ("sazdump", "sazserve")
$platforms = ("windows", "darwin")
$arches = "amd64"
$env:GO111Modules = 1

foreach ($platform in $platforms) {
	$env:GOOS = $platform
	foreach ($arch in $arches) {
		$env:GOARCH = $arch
		foreach ($package in $packages) {
			$name = "$package-$platform-$arch"
			if ($platform -eq "windows") {
				$name = "$package-$platform-$arch.exe"
			}
			go build -ldflags="-s -w" -o bin/$name "./cmd/$package"
		}
	}
}

