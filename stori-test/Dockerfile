FROM public.ecr.aws/docker/library/golang:1.22.2 as build-image
WORKDIR /src
COPY go.mod go.sum main.go ./
COPY infra ./infra
COPY entity ./entity
COPY repository ./repository
COPY usecase ./usecase
COPY adapters ./adapters
COPY utils ./utils

RUN go build -o lambda-handler

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build-image /src/lambda-handler .
COPY --from=build-image /src/utils/data/txns.csv .
COPY --from=build-image /src/utils/templates/txns_template.html .
ENTRYPOINT ./lambda-handler
