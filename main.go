package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

func BundleWindowsApp(libDir, name string) {
	os.Mkdir(name, 0644)
	exec.Command("cp", libDir+"/*.dll", name).Run()
	exec.Command("go", "build", "-ldflags", "-H windowsgui", "-o", name + "/" + name + ".exe").Run()
}

func BundleDarwinApp(name string) {
	shell := `
#!/bin/sh

cd src
go clean
go build -tags appbundle -o "$1"
cd ..

echo "
import re
from setuptools import setup

NAME = '$1'
ICON = 'bundle/app.icns'
VERSION = '1.0.0'
INFO = '$2'

PLIST = {'CFBundleName': NAME,
		 'CFBundleIconFile': ICON,
		 'CFBundleVersion': VERSION,
		 'CFBundleExecutable': NAME,
		 'CFBundleGetInfoString': INFO,
		 'CFBundleShortVersionString': re.search(r'\d+\.\d+', VERSION).group(0)}

APP = ['dummy.py']

FRAMEWORKS = ['/usr/local/opt/csfml/lib/libcsfml-window.2.3.dylib',
			  '/usr/local/opt/csfml/lib/libcsfml-graphics.2.3.dylib',
			  '/usr/local/opt/csfml/lib/libcsfml-audio.2.3.dylib',
			  '/usr/local/opt/csfml/lib/libcsfml-system.2.3.dylib']

OPTIONS = {'frameworks': FRAMEWORKS,
		   'iconfile': ICON,
		   'plist': PLIST}

setup(
	app=APP,
	options={'py2app': OPTIONS},
	setup_requires=['py2app']
)
" > setup.py

touch dummy.py
python3 setup.py py2app --semi-standalone
rm -f dummy.py
rm -f setup.py

rm -rf dist/"$1".app/Contents/MacOS/**
rm -rf dist/"$1".app/Contents/Resources/**

rm -rf dist/"$1".app/Contents/Frameworks/Python.framework
rm -rf dist/"$1".app/Contents/Frameworks/libcrypto.1.0.0.dylib
rm -rf dist/"$1".app/Contents/Frameworks/libssl.1.0.0.dylib
rm -rf dist/"$1".app/Contents/Frameworks/liblzma.5.dylib

mv -v src/"$1" dist/"$1".app/Contents/MacOS/"$1"
cp -v bundle_res/Info.plist dist/"$1".app/Contents/
cp -v bundle_res/App.icns dist/"$1".app/Contents/Resources/

install_name_tool -change \
  /usr/local/opt/csfml/lib/libcsfml-window.2.3.dylib \
  @executable_path/../Frameworks/libcsfml-window.2.3.dylib \
  dist/"$1".app/Contents/MacOS/"$1"

install_name_tool -change \
  /usr/local/opt/csfml/lib/libcsfml-graphics.2.3.dylib \
  @executable_path/../Frameworks/libcsfml-graphics.2.3.dylib \
  dist/"$1".app/Contents/MacOS/"$1"

install_name_tool -change \
  /usr/local/opt/csfml/lib/libcsfml-audio.2.3.dylib \
  @executable_path/../Frameworks/libcsfml-audio.2.3.dylib \
  dist/"$1".app/Contents/MacOS/"$1"

install_name_tool -change \
  /usr/local/opt/csfml/lib/libcsfml-system.2.3.dylib \
  @executable_path/../Frameworks/libcsfml-system.2.3.dylib \
  dist/"$1".app/Contents/MacOS/"$1"

mv dist/"$1".app "$1".app
rm -rf dist
rm -rf build
	`

	ioutil.WriteFile("bundle.sh", []byte(shell), 0644)
	exec.Command("sh", "bundle.sh", name, "An SFML game").Run()
}

func main() {
	if len(os.Args) == 0 {
		return
	}

	systemOS := runtime.GOOS

	name := os.Args[1]
	if systemOS == "windows" {
		libDir := os.Args[2]
		BundleWindowsApp(libDir, name)
	} else if systemOS == "darwin" {
		BundleDarwinApp(name)
	}
	fmt.Println("Done")
}
