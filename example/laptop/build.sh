#!/bin/bash
rm -f files.zip
zip -r files.zip files
go build
cat files.zip >> laptop
zip -A laptop
