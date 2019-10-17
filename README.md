# drop-box-it

This project emulates a file synchronization system similar to what Dropbox does, where you have a `source directory` 
and a `destination directory`. So whatever you write in the first one gets synchronized to the latter via HTTP.

To get everything up and running in a hassle-free way, I used Docker compose which will set up everything enabling you to get it up and running
with just one command which I'll describe later, but to understand how this project runs look at the `docker-compose.yml` file.

The project is divided into two services the first one which is more like CLI daemon type thing `drop` and the second one which is an HTTP server `box`.
Each service is self contain in his domain within `pkg` and everything gets wired up under `cmd`.
Both services are already preconfigured but feel free to play with the configuration, keep in mind that the default-src directory is `./srcDir` and for the destination one
is `./destDir`. 

To see it in action run the command `./script/start` wait a few seconds, you'll see every file that is contained 
within the src directory got synchronized to the destination one, and if you add, change or remove any of the files you'll see how they get synchronized too. 
Feel free of course to use any `src` directory is a CLI after all, you can update that value from the `docker-compose.yml` and then in the `drop` service update the `command`.

### Run

In order to run tests and linting
    
    ./script/test
    ./script/lint

If you want to see it in action
    
    ./script/start

### Improvements

- Add graceful shutdown
- Add another endpoint in the `box` service for file sync
- Add the possibility to sync different clients with diferent folders at the same time by namespacing the destination folder and having some kind of token to identify the client
