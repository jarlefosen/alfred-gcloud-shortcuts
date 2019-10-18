# g: GCP shortcuts for Alfred

## Download

An exported workflow should be available under [releases](https://github.com/jarlefosen/alfred-gcloud-shortcuts/releases) make sure you meet the [requirements](#requirements).

## Usage

Commands:
- `g <query>` for search
- `g-refresh` for updating list of projects

You are free to change the hotkey by editing the variable `hotkey` when importing the workflow.

Changing the variable `hotkey` from `g` to `gcp` results in commands like `gcp <query>` and `gcp-refresh`.

### Refresh projects list

Initially run `g-refresh` in Alfred to update the list of authenticated projects.

### Open product page

`g <project filter>` ↩️️ `BigQuery` ➡️ Opens BigQuery for the selected project.

`g My Project` ↩️ `Kube` ➡️ Opens Kubernetes Engine in GCP for project My Project.

## Requirements

```
If you initialized gcloud recently, make sure to save the authentication locally, see below.
```

- installed and authenticated `gcloud` https://cloud.google.com/sdk/
- coreutils `brew install coreutils`
- save [auth locally](https://github.com/jarlefosen/alfred-gcloud-shortcuts/issues/5#issuecomment-537852834): `gcloud auth application-default login`