<h3>Invoice Service</h3>

<h3>Description</h3>

<p>This repository will be used to manage invoice app</p>

<h3>Directory Structure</h3>

```
app
    L .github                        → Github CI/CD
        L workflows                   → Github workflow
    L cmd                            → Main Config
    L common                         → Common function
        L state                      → State FSM
    L config                         → Config app
    L constants                      → Constanta
    L controllers                    → Controller app
    L domain                         → Domain app
        L dto                        → DTO
        L models                     → Model object
    L middlewares                    → Middleware app
    L migrations                     → Migration app
    L mocks                          → Mocks
    L repositories                   → Repository app
    L routes                         → Route app
    L services                       → Service app
    L utils                          → Utils as helper
```

<h3>How to run</h3>

- Clone this repository
- go mod tidy
- copy .env.example to .env (if you want to run with consul)
- copy .config.example.json to .config.json
- make start

<h3>How to run with docker</h3>

- docker-compose up -d --build --force-recreate

<h3>How to run test</h3>

- make test

<h3>How to run linter</h3>

- make linter

<h3>How to build</h3>

- make build
