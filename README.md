# <p style="text-align: center">democrator</p>

A fullstack application for hosting live polls, (to be) built with Go, GraphQL, GORM and React.
Before running the project, make sure your PostgreSQL database instance is running and the connection string in `./api/.env` is correct.

## Re-generating Proto files

After introducing changes to the GQL schema in `./api/graph/schema.graphqls` make sure to re-generate the output with this command:

```
go run github.com/99designs/gqlgen
```

## Goals

This application will allow the user to:

- Create a live poll and generate a random link to it
- Enter the poll by with an appropriate link and vote
- See live poll results in the poll page
- Set a time in which the poll expires
