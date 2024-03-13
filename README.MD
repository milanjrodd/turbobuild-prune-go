# pipeline

turbo prune $project --docker
turbo-prune go $project --docker
  - go work use -r .
  - copy go.work out/json and out/full folders
  - copy go.work.sum file to out/json folder
  - go to apps/$project folder
    - copy go.mod file to out/json/apps/$project folder
    - copy $project folder to out/full/apps/$project
  - go to packages folder
    - copy all go.mod files from all only go packages to out/json/packages/* folders
    - copy all only go packages to out/full/packages/* folders