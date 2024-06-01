FROM python:3.10.13-slim

WORKDIR /docs

# mkdocs mkdocs-material mkdocs-minify-plugin mkdocs-material-extensions
COPY pyproject.toml poetry.lock ./

RUN pip install poetry==1.7.0 && poetry install --no-root --no-directory

COPY . .
RUN poetry install

EXPOSE 8123
ENTRYPOINT ["poetry", "run", "mkdocs", "serve", "--dev-addr=0.0.0.0:8123"]
