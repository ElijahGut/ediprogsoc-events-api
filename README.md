# EUPS Events API

The EUPS Events API helps EUPS manage their events. The API can be accessed <a href="https://ediprogsoc.co.uk/app">here</a>. Only EUPS committee members are authorised to use the API. 

The API is built with Go Fiber and connects to Cloud Firestore.

</br>

## Prerequisites

</br>

* Go >= v.1.18
* firebase CLI >= 11.9.0

</br>

## Local set up

</br>

For local development, we interact with the firestore emulator instead of the live version.

**NB: All commands should be run at the top level of the repo**

1. Clone the repo: `git clone REPO_LINK`
2. Generate a private key from firebase. Follow the docs <a href="https://firebase.google.com/docs/admin/setup">here</a>. The relevant section is "Initialize the SDK". Once the key has been generated, **place it at the top level of the repo** (the private key file name should begin with `prog-soc-map`)
3. Init env vars: `source .bashrc FIREBASE_PRIVATE_KEY`. Substitute your private key file name into FIREBASE_PRIVATE_KEY. **You will need to do this every time you start a new shell session**
4. Download Go deps: `go get ./...`
5. Enable firestore emulator for local dev: `firebase init emulators`. This command will start up a wizard that walks you through the emulator set-up. You only need the firestore emulator. Set the firestore emulator port to 8081, and the UI port to 4000. If you need more information explore the docs <a href="https://firebase.google.com/docs/emulator-suite/install_and_configure">here<a>. Once complete, the wizard should have generated a `firebase.json` file. You can consult this to double check/edit your config
6. Start the emulator on one terminal session: `firebase emulators:start`. You can visit the UI at localhost:4000 and explore the firestore emulator
7. Start the Go application on a second terminal session (remember to init env vars first, otherwise this won't work): `go run .`
8. Test out the API using Postman or a similar tool. Requests should be made to `localhost:8080` You should be able to see your changes in the emulator UI
9. If you need more info on the API, consult the Swagger docs at `localhost:8080/swagger`

</br>

## Testing 

</br>

To test the application, make sure that the firestore emulator and Go application are running. Then, run `./run_tests.sh`. You may need to `chmod 777` if you get a permission error.

If you get a google creds error, run `source .bashrc FIREBASE_PRIVATE_KEY` again.

Testing is done with the testify package. Test results show which tests passed/failed, and includes coverage info.


## Contact

Please contact ediprogsoc@gmail.com for any issues or questions.

