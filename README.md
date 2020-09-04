# IPA-Manager
 
### Description
IPA Manager was created to simplify the transfer and installation of mobile client builds for iOS.

### How it works
IPA Manager monitors the "\\share\ipas" directory to see if files with the ".ipa " extension appear in it. When such a file appears, it copies it to itself, parses the Info.plist of the application, and adds the collected information to its database and forms a plist for subsequent installation of the IPA. Then deletes the source file from the \\share\ipas directory.

Also this checks for the existence of an IPA in the database using the CFBundleVersion parameter. If such a file has already been added to the database, it is simply deleted from the \\share\ipas directory.

When creating a page with an IPA list, the list is ordered in reverse order by the date the IPA was added to facilitate access to the last added IPAs.

### Key features
Installing an IPA using the button from the page
Generating a QR code for installing an IPA by pointing the mobile device's camera at the PC monitor screen

## License
[MIT](LICENSE)
