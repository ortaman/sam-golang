FROM public.ecr.aws/docker/library/golang:1.22.2 as build-image
WORKDIR /src
COPY go.mod go.sum main.go ./
RUN go build -o lambda-handler
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build-image /src/lambda-handler .
ENTRYPOINT ./lambda-handler