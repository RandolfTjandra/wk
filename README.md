# wk - WaniKani CLI Tool [![Build](https://github.com/RandolfTjandra/wk/actions/workflows/build.yml/badge.svg)](https://github.com/RandolfTjandra/wk/actions/workflows/build.yml)

WaniKani CLI Tool is a command line interface (CLI) tool designed to connect to the WaniKani API to retrieve and display personal learning statistics for Japanese learners. It is written in Go and uses a TUI library to create an interactive and visually pleasing interface. 

With WaniKani CLI Tool, you can quickly and easily check your current level, review statistics, and progress charts. It's perfect for those who want to keep track of their WaniKani learning journey without having to navigate through the WaniKani website.

## Features

- Interactive and user-friendly TUI interface
- Display of current level, review statistics, and progress charts
- Quick and easy access to personal learning statistics
- Secure and reliable connection to the WaniKani API

## Installation

## Dependencies
### Redis
[Install Guide](https://redis.io/docs/getting-started/installation/)

To install WaniKani CLI Tool, simply clone the repository and build the binary file using the following command:

```
$ git clone git@github.com:RandolfTjandra/wk.git
$ cd wk 
$ make build-wk
$ ./dist/wk/bootstrap
```

Once you have built the binary, you can run the tool using the following command:


## Usage

When you run WaniKani CLI Tool, you will be prompted to enter your API key. You can find your API key on the WaniKani website by going to the API Tokens section of your account settings.

Once you have entered your API key, you will be taken to the main menu. From there, you can navigate through the various options using the arrow keys and enter key. 

## Contributing

We welcome contributions to WaniKani CLI Tool! If you find a bug or have an idea for a new feature, please open an issue on the repository. If you would like to contribute code, please fork the repository and submit a pull request.

## License

WaniKani CLI Tool is released under the MIT License. See [LICENSE](https://github.com/randolftjandra/wk/blob/main/LICENSE) for more information.
