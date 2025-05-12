# MisterFPGA Neon68k unpacker

Take the [Neon68k zip archive](https://neon68k.com/) and unzip one of the collections into an appropriate structure so that it can be upload directly to MisterFPGA.

- Sparsely documented for now, but it should be straightforward to figure out.
- Only lightly tested. I believe it's not dangerous to use [caveat emptor].
- I've only tested the Windows version. Relying on Go's cross-platform compile to make good on the other versions!

## Usage

`.\misterfpga-neon68k-unzipper.exe --src-zip "C:\Users\James\Desktop\neon68k.zip" --collection "Neon68K-20250428\MiSTer Upscaler\English\_Main" --dest-folder .\folder-for-misterfpga`

The contents of the destination folder don't get deleted, so you can run this a few times to bring together a selection of the collections.

## Possible extensions

- FTP direct to MisterFPGA

## License and Credit

Credit to James Rutherford and link to the root repository (https://github.com/creativenucleus/misterfpga-neon68k-unzipper)

This is a liberal license - [Code Credit 1.1.0](https://codecreditlicense.com/license/1.1.0) - you are welcome to fork and modify as you like, but ensure you retain this original credit.

Please send me a mention / pop me a hello if you use it :) [Mastodon](https://mastodon.social/@jtruk)

## Thanks

The [Neon68K team](https://neon68k.com/) of course!