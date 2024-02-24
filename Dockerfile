FROM golang:1.21

WORKDIR /app

COPY . .

WORKDIR /app/cmd

CMD ["go", "run", "."]

