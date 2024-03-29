FROM ubuntu:22.04

# update container certificates
RUN apt-get -y update && apt-get install -y ca-certificates && update-ca-certificates

# set the Current Working Directory inside the container
WORKDIR ~/app

# copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY ./build/ .

# this container exposes port 8081 to the outside world
EXPOSE 8081

# run the server executable
CMD ["./client"]
