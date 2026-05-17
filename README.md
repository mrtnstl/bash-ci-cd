# bash-ci-cd

![demo](https://github.com/mrtnstl/bash-ci-cd/blob/main/docs/bash_demo.gif "demo")

Continuous Integration and Continuous Delivery runner in bash.

The pipeline steps are defined in the `tasks` directory. These could be easily extended if needed.

## What it does? (NOT IMPLEMENTED)

the pipeline does
- recieve a webhook from your github
- checks the code out from a git repository
- installs dependencies
- runs linter
- runs tests
- builds docker image
- pushes docker image to artifact repository
- spins up the new instance of the application on the server
- sends notification when it's done or an error occurs

## Run

1. Set your variables in `config.sh`

    WARNING: these values are sitting in your config file in plain text!

    ```bash
    ENV="production" # currently used for testing

    ```

2. Make sure you have execution privileges for `start.sh`. The script will attempt to set it to the rest of the files where needed.

    ```bash
    chmod u+x ./start.sh
    ```

3. Run script

    ```bash
    ./start.sh
    ```
