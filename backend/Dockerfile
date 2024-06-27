FROM ubuntu

WORKDIR /app

COPY go.mod go.sum ./

RUN apt-get update
RUN echo Y | apt-get install curl golang ca-certificates
RUN update-ca-certificates
RUN echo Y | (curl -s https://raw.githubusercontent.com/objectbox/objectbox-go/main/install.sh | bash)

COPY . .

RUN go build

EXPOSE 12345
EXPOSE 3000

CMD [ "./backend", "product_master_data.csv", "store.csv" ]
