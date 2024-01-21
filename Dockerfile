FROM golang

WORKDIR /app

COPY . .

WORKDIR /app/cmd

CMD ["go", "run", "."]