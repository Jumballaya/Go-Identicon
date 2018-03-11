# Identicon Generator



### Installation:

  Must have [Mage](https://github.com/magefile/mage) to build with the magefile

  1. Clone this repo
  2. `cd` into the folder
  3. Run `mage install`

### Usage
  1. `identicon SEED_STRING`
  2. If you want to specify the output directory (./dist by default) use the `-o DIR` flag
  3. Build with `mage build`
  4. Test with `mage test`
  5. Install with `mage install`


  If I wanted to generate an identicon based on the project name and save it to my pictures folder it might look like the following

  `identicon -o ~/Documents/Pictures Identicon`

  <img src="assets/Identicon.png">
