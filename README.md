# MisterFPGA Neon68k Unpacker

[Neon68k](https://neon68k.com/) is an archive of game setups for the [X68000](https://en.wikipedia.org/wiki/X68000).  
It's a curated zip-file of zip-files, and in order to use the games with MisterFPGA you need to unpack them into a specific folder structure.

The Neon68K website [outlines how to do this](https://neon68k.com/docs#unpacking-the-games) and it's fiddly - especially for multiple games.

This tool is a simple command-line utility to make unpacking easy. It can unpack a selection of games to a folder ready for you to upload to MiSTer, or - even easier - it can directly FTP them to a MiSTer on your local network.

## Caveats

- This tool has only been lightly tested. I believe it's not dangerous to use [caveat emptor], but it may fail without a clear error message.
- I've only tested on Windows. One of Go's strengths is cross-platform compilation, and I'm relying on that to build the other versions.
- I would appreciate any feedback to confirm that it works well, and welcome issue reports for suspected bugs via GitHub.

## Usage

The tool requires you to specify a path to match for the zip-files inside the archive to specify a particular collection.

The most common use-case will be to extract the English-Main collection, something like this:

`.\misterfpga-neon68k-unpacker.exe --src-zip ".\path-to-source\neon68k.zip" --src-collection "Neon68K-20250428\MiSTer Upscaler\English\_Main" --dest-type=file --dest-folder .\path-to-folder-for-output`

The tool does not delete anything in the destination folder, so you can run it multiple times with different collections to build up your own selection of the games.

## Arguments

- `--src-zip` - the path to the zip-file you downloaded from the Neon68K website. This is required.
- `--src-collection` - the path to the collection you want to extract. This is required. Take note of the example above, and you will need to manually open the zip-file if you'd like a different collection. Examples of this might be:
    - `Neon68K-20250428\MiSTer Upscaler\English\_Main` (for the main collection)
    - `Neon68K-20250428\MiSTer Upscaler\English\_Keyboard+Mouse` (for the keyboard and mouse collection)
    - `Neon68K-20250428\MiSTer Upscaler\English` (for all the English collections - this includes games with major bugs)
    - `Neon68K-20250428\MiSTer Upscaler\English\_Main\Sol-Feace.zip` (for a single game)
    - `Neon68K-20250428\External 4K Upscaler\Japanese\_Main` (for the main Japanese collection, with an external scaler)
- `--dest-type` - the type of destination you want. This can be either `file` or `ftp`.
    - `file` - the destination is a folder on the filesystem.
    - `ftp` - the destination is an FTP server. For uploading directly to a MiSTerFPGA on your local network.
- `--dest-folder` - the folder on the filesystem to unpack the files to. This is required if you are using `file` as the destination type.
- `--dest-ftp` - the IP address of your MiSTerFPGA. This is required if you are using `ftp` as the destination type. The tool will attempt to connect to the MiSTer using the default FTP port (21) and will use the default MiSTer credentials (username: `root`, password: `1`) - please let me know if you use anything different to that. I'll consider modifying the tool if there's a requirement.

## License and Credit

Credit: James Rutherford and link to the root repository (https://github.com/creativenucleus/misterfpga-neon68k-unpacker)

This is a liberal license - [Code Credit 1.1.0](https://codecreditlicense.com/license/1.1.0) - you are welcome to fork and modify as you like, but ensure you retain this original credit.

Please send me a mention / pop me a hello if you use it :) [Mastodon](https://mastodon.social/@jtruk)

## Thanks

- The [Neon68K team](https://neon68k.com/) of course!
- Authors of Go packages used in this project.