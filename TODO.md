## Backend

- [x] Replace hardcoded logins with DB calls
  - [x] Add hashing and salting for passwords
  - [x] Use JWTs for auth
    - [ ] Implement roles (basic admin + user for now)
- [x] Modularise code instead of main.go monolith
  - [ ] Move db and token config into config package (after implementing JWT)
- [ ] Add ability for verified user to update their password
  - [ ] Add checks for registry/update to ensure password has special char and length

## Frontend

- [ ] Add a frontend beyond the basic CRA (after tokens implemented)
  - [ ] Add register/login calls from frontend
  - [ ] Means of displaying the results of GetUserList() and welcome funcs - maybe a modal?
