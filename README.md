<h1 align="center">Welcome to MyNesEmulator 👋</h1>
<p>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> My first Golang project ! A basic Nintendo NES/FAMICOM emulator

## Notes

  The work is in progress the PPU/CPU are finished. I need to add the APU and Controllers.
  All the documentations that I am using will be provided as soon as the project is finished ;) 

## Dependencies

    github.com/go-gl/gl/v2.1/gl
    github.com/go-gl/glfw/v3.3/glfw
    github.com/hadi-ilies/MyNesEmulator/src/constant
    github.com/hadi-ilies/MyNesEmulator/src/nes
    github.com/hadi-ilies/MyNesEmulator/src/nes/nescomponents

## Installation

```sh
$>go get github.com/go-gl/gl/v2.1/gl
$>go get github.com/go-gl/glfw/v3.3/glfw
$>go get github.com/hadi-ilies/MyNesEmulator/src/constant
$>go get github.com/hadi-ilies/MyNesEmulator/src/nes
$>go get github.com/hadi-ilies/MyNesEmulator/src/nes/nescomponents
```

## Usage

```sh
$>go run src/main.go assets/your_rom.nes
```
### OR

```sh
$>go build -o MyNesEmulator src/main.go

$>./MyNesEmulator assets/your_rom.nes
```

## Author

👤 **hadi-ilies.bereksi-reguig**

* Github: [@hadi-ilies](https://github.com/hadi-ilies)

## Show your support

Give a ⭐️ if this project helped you!

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_
