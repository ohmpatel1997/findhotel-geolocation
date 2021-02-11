FROM alpine:3.10.3

LABEL maintainer="ohm patel"
LABEL Description="find hotel assignment"

RUN apk update && \
    adduser -S -D -H -h /app findhotel_user

# add a bit of security here and run as user, instead of root
USER findhotel_user

# Copy over the app files
COPY ./main.go /app/

WORKDIR /app

CMD ["./main"]
