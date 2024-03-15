# TurboBuild Prune Go

TurboBuild Prune Go is an npm package that provides a convenient way to prune Golang apps in a [Turbo](https://turbo.build) monorepo.

## Installation

Add TurboBuild Prune Go to your turbo repo:
```bash
npm install turbobuild-prune-go --global
```

## Usage
```bash
turbobuild-prune-go -p backend
```

## Contribution

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

### Prerequisites

- Go 1.16 or higher
- Node.js 14 or higher
- pnpm

### Installation
1. Fork the Project
2. Install Go dependencies:
```bash
go mod tidy
```
3. Install node modules:
```bash
pnpm install
```

### Development
To run the Go code, use the following command:

```bash
go run turbobuild-prune-go.go
```

### Sharing
Don't forget to share changes:

1. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
2. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
3. Push to the Branch (`git push origin feature/AmazingFeature`)
4. Open a Pull Request