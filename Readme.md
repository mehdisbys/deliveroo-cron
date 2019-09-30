
### Proposed solution

This repository aims to fulfill the requirements described below. 
In a nutshell this program allows to visualise a cron command (minus the year) in a table with a breakdown for each column.

The approach is as follows :

- Cron format supports multiple input types :

      - 5 : a fixed value
      - 5-10 : a range
      - */5 : a frequency
      - 5,7 : a list of values 
      
- For each type there is a different parser in [parsers.go](domain/parsers.go)
- Each type has a different delimiter defined in [cron.go](domain/cron.go)
- Each time component (minutes, hours, etc..) has a specific range defined in [cron.go](domain/cron.go)

#### How to run it 

Dependencies :

- Install `golang-ci` if not installed.
- go 1.12 or above

After cloning the repo do :

`export GO111MODULE=on` //enables go modules

`go mod download`

`make test`

`make binary`

Try it with a cron command. it can take any number of parameters as along as :
 - The program to execute is the last element
 - The 5 columns are right before the the program to execute.
 
 Example :
 
 `./deliveroo-cron foo bar xyzfile-name -arguments 3-45/15 0 1,15 2,3 1-5 /usr/bin/cat`
 
 yields :
 
 ```
 minutes       3 18 33
 hour          0
 day of month  1 15
 month         2 3
 day of week   1 2 3 4 5
 command       /usr/bin/cat
```

#### Errors


- `zsh: no matches found: */15` :Your shell could be trying to expand the `*` in the cron command try disabling __globbing__ with ` set -o noglob` (works for zsh. If that does not work research it for your particular shell.

- Entering invalid data such as impossible ranges is checked and will return an error like :

```
./deliveroo-cron file-name -arguments */15 0 1,15 * 1-9 /usr/bin/find
 input 1-9 : 9 is superior to allowed range [1 7]
```
-----------


Write a command line application or script which parses a cron string and expands each field to show the times at which it will run. You may use any of the following languages: Ruby, Scala, JavaScript, Python, Go, Java, or C# (using .NET Core as we need it to run on OS X/Linux).
You should only consider the standard cron format with five time fields (minute, hour, day of month, month, and day of week) plus a command; you do not need to handle the special time strings such as "@yearly".


The cron string will be passed to your application on the command line as a single argument on a single line. For example:


`~$ your-program */15 0 1,15 * 1-5 /usr/bin/find`


Or:


`~$ application-commands file-name -arguments */15 0 1,15 * 1-5 /usr/bin/find`


The output should be formatted as a table with the field name taking the first 14 columns and the times as a space-separated list following it.
For example, the following input argument:
should yield the following output:

```
minutes       0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find


```