# Identicon Generator



### Installation:
  1. Clone this repo
  2. cd into the folder
  3. Install cli dependency
  4. `make install`

### Usage
  1. `identicon SEED_STRING`
  2. If you want to specify the output directory (./dist by default) use the `-o DIR` flag
  3. Build with `make`
  4. Test with `make test`
  5. Install with `make install`


  If I wanted to generate an identicon based on the project name and save it to my pictures folder it might look like the following

  `identicon -o ~/Documents/Pictures Identicon`

  <img src="assets/Identicon.png">
