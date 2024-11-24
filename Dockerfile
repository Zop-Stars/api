FROM alpine:edge

RUN apk add --no-cache tzdata ca-certificates

COPY configs ./configs
COPY service_instruction.txt ./service_instruction.txt
COPY main ./main

RUN chmod +x /main

EXPOSE 8000

CMD ["/main"]
