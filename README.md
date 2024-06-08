# Watch-Dawg

## Description

Watch-Dawg is a little tool written in Go to observe your download folder and move known files into organized customizable folders.

E.g. when you download a file with extension '.txt' or '.pdf' it automatically gets moved to a folder 'Documents' in your download folder.

It should work fine on Windows, Linux and MacOS.

## Configuration

You can configure the tool by editing the `config.json` file. Also you can change the constant `DOWNLOAD_FOLDER` to your custom download folder in `main.go` file.

## Disclaimer

Please note: This tool is a work in progress and helps me learning to code in go. Use it at your own risk.
