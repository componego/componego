#!/usr/bin/env sh

set -o errexit
set -o nounset

cd "$(dirname "$(realpath -- "$0")")";

docker-compose run --entrypoint "poetry run mkdocs build --clean --site-dir=/docs/build" docs

cp ../LICENSE ./build/LICENSE
if [ -e "../NOTICE" ]; then
    cp ../NOTICE ./build/NOTICE
fi

cat > build/README.md << EOF
Website available [here](https://componego.github.io/).

---

These files are auto-generated files from the [main repository](https://github.com/componego/componego)
so you don't have to commit changes directly to this repository.

The [license](./LICENSE) of this repository matches the license of the parent repository.
EOF

cat > build/.gitignore << EOF
.idea
**/.idea
.vs-code
**/.vs-code
EOF
