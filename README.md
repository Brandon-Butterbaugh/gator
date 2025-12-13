Postgres and Go will need to be installed to run the program.

Make sure to install the gator CLI before running.
 - Do this by running 
        go install github.com/Brandon-Butterbaugh/gator@latest
 - or if the repo is already cloned locally, go to the project director and use
        go install

Enter the psql shell with
 - Mac: psql postgres
 - Linux: sudo -u postgres psql
and create a new gator database with
CREATE DATABASE gator;

get the connection string for this database. It will look something like
    "postgres://bobbill:@localhost:5432/gator"

Move into the sql/schema directory with
cd sql/schema
and run the up migration. It will look like
goose postgres "postgres://bobbill:@localhost:5432/gator" up

To run the program you will need to create a .gatorconfig.json in your home directory
After making the .gatorconfig.json file add 
```
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

the string for db_url is going to be your connection string and has format
    protocol://username:password@host:port/database
replace the string with your connection string and don't forget to add
    ?sslmode=disable
to the end of the connection string

current_user_name can remain as it is since the first command you should run is
    gator register <name>
this will create a new user and set the current user to the one registered.

A list of commands are

    login
        logs a user in
    example:    gator login bob

    register
        creates a new user
    example:    gator register bob

    reset
        clears all databases
    example:    gator reset

    users
        prints all registered users
    example:    gator users

    addfeed
        Add a feed to the database with it's name and url
    example:    gator addfeed "Hacker News RSS" "https://hnrss.org/newest"

    feeds
        prints all feeds in the database
    example:    gator feeds

    follow
        follow a feed added by a different user using the url
    example: gator follow "https://hnrss.org/newest"

    following
        prints the feeds the current user is following
    example: gator following

    unfollow
        removes a feed from your list of follows using the url
    example: gator unfollow "https://hnrss.org/newest"

    agg
        aggregates the feeds into posts every designated interval
            intervals can be like 10s  5m  1h
    example: gator agg 10s

    browse
        will print two aggregated posts unless specified
    example: gator browse 5