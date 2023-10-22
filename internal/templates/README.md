
###### To gen node_modules, run
```sh
cd ./internal/templates
npm ci
```

###### To gen and update dist, run
```sh
cd ./internal/templates
npx tailwindcss -i ./src/input.css -o ./dist/output.css --watch
```
