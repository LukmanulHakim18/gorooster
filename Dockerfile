FROM golang:1.19-alpine as build
ARG SSH_KEY

RUN apk add git && \
    apk add openssh

RUN mkdir -p /root/.ssh
RUN echo "$SSH_KEY" > /root/.ssh/id_rsa && chmod 0600 /root/.ssh/id_rsa

WORKDIR /app 
COPY . ./
RUN ls -all && pwd
RUN go mod tidy
# make sure your domain is accepted
RUN touch /root/.ssh/known_hosts
RUN echo -e "Host *\n\tStrictHostKeyChecking no\n\tForwardAgent yes\n\n" > ~/.ssh/config
RUN git config --global --add url."git@gitssh.bluebird.id:".insteadOf "https://git.bluebird.id"
RUN go clean && CGO_ENABLED=0 go build
RUN ls -all && pwd

FROM gcr.io/distroless/static
COPY --from=build /app/gorooster gorooster
ENTRYPOINT ["./gorooster"]


